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
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

var _ Cool = (*cool)(nil)

var (
	ErrInitialization = errors.New("cool initialization failed")
	ErrNilConn        = errors.New("connection is nil")
	ErrClosed         = errors.New("cool is closed")
)

// Cool net.Conn pool
type Cool interface {
	Get() (net.Conn, error)
	Close()
	Size() int
}

type Producer func() (net.Conn, error)

type cool struct {
	sync.Once
	sync.RWMutex
	options  *options
	connC    chan net.Conn
	producer Producer
}

func New(init, max int, producer Producer, opts ...Option) (Cool, error) {
	if init < 0 || max <= 0 || init > max || producer == nil {
		return nil, ErrInitialization
	}
	c := &cool{
		options:  newOptions(opts...),
		connC:    make(chan net.Conn, max),
		producer: producer,
	}
	for i := 0; i < init; i++ {
		conn, err := producer()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("error produce connection: %v", err)
		}
		c.connC <- c.wrap(conn)
	}
	return c, nil
}

func (c *cool) Get() (net.Conn, error) {
	c.RLock()
	connC := c.connC
	producer := c.producer
	c.RUnlock()
	if connC == nil {
		return nil, ErrClosed
	}
	select {
	case conn, ok := <-connC:
		if !ok {
			return nil, ErrClosed
		}
		if timeout := c.options.connIdleTimeout; timeout > 0 && time.Now().After(conn.(*Conn).t.Add(timeout)) {
			return c.produceWrap(producer)
		}
		return conn, nil
	default:
		return c.produceWrap(producer)
	}
}

func (c *cool) Close() {
	c.Do(func() {
		connC := c.connC
		c.options = nil
		c.connC = nil
		c.producer = nil
		close(connC)
		for conn := range connC {
			_ = conn.Close()
		}
	})
}

func (c *cool) Size() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.connC)
}

func (c *cool) put(conn net.Conn) error {
	if conn == nil {
		return ErrNilConn
	}
	c.RLock()
	defer c.RUnlock()
	if c.connC == nil {
		// cool is closed
		return conn.Close()
	}
	select {
	case c.connC <- c.wrap(conn):
		return nil
	default:
		return conn.Close()
	}
}

func (c *cool) wrap(conn net.Conn) net.Conn {
	return &Conn{
		Conn: conn,
		c:    c,
		t:    time.Now(),
	}
}

func (c *cool) produceWrap(producer Producer) (net.Conn, error) {
	conn, err := producer()
	if err != nil {
		return nil, err
	}
	return c.wrap(conn), nil
}
