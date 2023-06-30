// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except In compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to In writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package cool

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	InitialCap = 5
	MaximumCap = 30
	network    = "tcp"
	address    = "127.0.0.1:7246"
	producer   = func() (net.Conn, error) {
		return net.Dial(network, address)
	}
)

func init() {
	// used for producer function
	go simpleTCPServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled
	rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
}

func TestNew(t *testing.T) {
	_, err := newCool()
	assert.Nil(t, err)
}

func TestPool_Get_Impl(t *testing.T) {
	p, _ := newCool()
	defer p.Close()

	conn, err := p.Get()
	assert.Nil(t, err)

	_, ok := conn.(*Conn)
	assert.True(t, ok)
}

func TestPool_Get(t *testing.T) {
	p, _ := newCool()
	defer p.Close()

	_, err := p.Get()
	assert.Nil(t, err)

	// after one get, current capacity should be lowered by one.
	assert.Equal(t, InitialCap-1, p.Len())

	// get them all
	var wg sync.WaitGroup
	for i := 0; i < (InitialCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Get()
			assert.Nil(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, 0, p.Len())

	_, err = p.Get()
	assert.Nil(t, err)
}

func TestPool_Put(t *testing.T) {
	p, err := New(0, MaximumCap, producer)
	assert.Nil(t, err)
	defer p.Close()

	// get/create from the pool
	connC := make([]net.Conn, MaximumCap)
	for i := 0; i < MaximumCap; i++ {
		conn, _ := p.Get()
		connC[i] = conn
	}

	// now put them all back
	for _, conn := range connC {
		_ = conn.Close()
	}

	if p.Len() != MaximumCap {
		t.Errorf("Put error len. Expecting %d, got %d", 1, p.Len())
	}

	conn, _ := p.Get()
	p.Close() // close pool

	_ = conn.Close() // try to put into a full pool
	assert.Equal(t, 0, p.Len())
}

func TestPool_PutUnusableConn(t *testing.T) {
	p, _ := newCool()
	defer p.Close()

	// ensure pool is not empty
	conn, _ := p.Get()
	_ = conn.Close()

	poolSize := p.Len()
	conn, _ = p.Get()
	_ = conn.Close()
	assert.Equal(t, poolSize, p.Len())

	conn, _ = p.Get()
	if pc, ok := conn.(*Conn); !ok {
		t.Errorf("impossible")
	} else {
		pc.MarkUnusable()
	}
	_ = conn.Close()
	assert.Equal(t, poolSize-1, p.Len())
}

func TestPool_UsedCapacity(t *testing.T) {
	p, _ := newCool()
	defer p.Close()
	assert.Equal(t, InitialCap, p.Len())
}

func TestPool_Close(t *testing.T) {
	p, _ := newCool()

	// now close it and test all cases we are expecting.
	p.Close()

	c := p.(*cool)

	assert.Nil(t, c.connC)
	assert.Nil(t, c.producer)

	_, err := p.Get()
	assert.NotNil(t, err)
	assert.Equal(t, 0, p.Len())
}

func TestPoolConcurrent(t *testing.T) {
	p, _ := newCool()
	pipe := make(chan net.Conn, 0)

	go func() {
		p.Close()
	}()

	for i := 0; i < MaximumCap; i++ {
		go func() {
			conn, _ := p.Get()

			pipe <- conn
		}()

		go func() {
			conn := <-pipe
			if conn == nil {
				return
			}
			_ = conn.Close()
		}()
	}
}

func TestPoolWriteRead(t *testing.T) {
	p, _ := New(0, 30, producer)

	conn, _ := p.Get()

	msg := "hello"
	_, err := conn.Write([]byte(msg))
	assert.Nil(t, err)
}

func TestPoolConcurrent2(t *testing.T) {
	p, _ := New(0, 30, producer)

	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				conn, _ := p.Get()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				_ = conn.Close()
				wg.Done()
			}(i)
		}
	}()

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			conn, _ := p.Get()
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			_ = conn.Close()
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestPoolConcurrent3(t *testing.T) {
	p, _ := New(0, 1, producer)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		p.Close()
		wg.Done()
	}()

	if conn, err := p.Get(); err == nil {
		_ = conn.Close()
	}

	wg.Wait()
}

func newCool() (Cool, error) {
	return New(InitialCap, MaximumCap, producer)
}

func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			buffer := make([]byte, 256)
			_, _ = conn.Read(buffer)
		}()
	}
}
