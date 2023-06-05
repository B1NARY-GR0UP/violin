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

package violin

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestQuickStart(t *testing.T) {
	defer goleak.VerifyNone(t)
	v := New()
	defer v.Shutdown()
	v.Submit(func() {
		fmt.Println("Hello, VIOLIN!")
	})
}

func TestVIOLIN(t *testing.T) {
	defer goleak.VerifyNone(t)
	v := New(WithMaxWorkers(2))
	requests := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	respChan := make(chan string, len(requests))
	for _, r := range requests {
		r := r
		v.Submit(func() {
			respChan <- r
		})
	}
	v.ShutdownWait()
	close(respChan)
	respSet := map[string]struct{}{}
	for rsp := range respChan {
		respSet[rsp] = struct{}{}
	}
	assert.GreaterOrEqual(t, len(respSet), len(requests))
	for _, req := range requests {
		_, ok := respSet[req]
		assert.True(t, ok)
	}
}

func TestMaxWorkers(t *testing.T) {
	defer goleak.VerifyNone(t)

	v := New(WithMaxWorkers(0))
	v.Shutdown()
	if v.MaxWorkerNum() != 1 {
		t.Fatal("should have created one worker")
	}

	max := 13
	v = New(WithMaxWorkers(max))
	defer v.Shutdown()

	assert.Equal(t, max, v.MaxWorkerNum())

	started := make(chan struct{}, max)
	release := make(chan struct{})

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < max; i++ {
		v.Submit(func() {
			started <- struct{}{}
			<-release
		})
	}

	// Wait for all queued tasks to be dispatched to workers.
	assert.Equal(t, len(v.waitingChan), int(v.WaitingTaskNum()))

	timeout := time.After(5 * time.Second)
	for startCount := 0; startCount < max; {
		select {
		case <-started:
			startCount++
		case <-timeout:
			t.Fatal("timed out waiting for workers to start")
		}
	}

	// Release workers.
	close(release)
}

func TestSubmitWait(t *testing.T) {
	defer goleak.VerifyNone(t)

	v := New(WithMaxWorkers(1))
	defer v.Shutdown()

	// Check that these are noop.
	v.Submit(nil)
	v.SubmitWait(nil)

	done1 := make(chan struct{})
	v.Submit(func() {
		time.Sleep(100 * time.Millisecond)
		close(done1)
	})
	select {
	case <-done1:
		t.Fatal("Submit did not return immediately")
	default:
	}

	done2 := make(chan struct{})
	v.SubmitWait(func() {
		time.Sleep(100 * time.Millisecond)
		close(done2)
	})
	select {
	case <-done2:
	default:
		t.Fatal("SubmitWait did not wait for function to execute")
	}
}

func TestStopRace(t *testing.T) {
	defer goleak.VerifyNone(t)

	max := 13

	v := New(WithMaxWorkers(max))
	defer v.Shutdown()

	workRelChan := make(chan struct{})

	var started sync.WaitGroup
	started.Add(max)

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < max; i++ {
		v.Submit(func() {
			started.Done()
			<-workRelChan
		})
	}

	started.Wait()

	const doneCallers = 5
	stopDone := make(chan struct{}, doneCallers)
	for i := 0; i < doneCallers; i++ {
		go func() {
			v.Shutdown()
			stopDone <- struct{}{}
		}()
	}

	select {
	case <-stopDone:
		t.Fatal("Stop should not return in any goroutine")
	default:
	}

	close(workRelChan)

	timeout := time.After(time.Second)
	for i := 0; i < doneCallers; i++ {
		select {
		case <-stopDone:
		case <-timeout:
			v.Shutdown()
			t.Fatal("Timeout waiting for Stop to return")
		}
	}
}

func TestWorkerLeak(t *testing.T) {
	defer goleak.VerifyNone(t)

	const workerCount = 100

	v := New(WithMaxWorkers(workerCount))

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < workerCount; i++ {
		v.Submit(func() {
			time.Sleep(time.Millisecond)
		})
	}

	// If v.Shutdown() is not waiting for all workers to complete, then goleak
	// should catch that
	v.Shutdown()
}

func BenchmarkExecute1Worker(b *testing.B) {
	benchmarkExecWorkers(1, b)
}

func BenchmarkExecute2Worker(b *testing.B) {
	benchmarkExecWorkers(2, b)
}

func BenchmarkExecute4Workers(b *testing.B) {
	benchmarkExecWorkers(4, b)
}

func BenchmarkExecute16Workers(b *testing.B) {
	benchmarkExecWorkers(16, b)
}

func BenchmarkExecute64Workers(b *testing.B) {
	benchmarkExecWorkers(64, b)
}

func BenchmarkExecute1024Workers(b *testing.B) {
	benchmarkExecWorkers(1024, b)
}

func benchmarkExecWorkers(n int, b *testing.B) {
	v := New(WithMaxWorkers(n))
	defer v.Shutdown()
	var allDone sync.WaitGroup
	allDone.Add(b.N * n)

	b.ResetTimer()

	// Start workers, and have them all wait on a channel before completing.
	for i := 0; i < b.N; i++ {
		for j := 0; j < n; j++ {
			v.Submit(func() {
				allDone.Done()
			})
		}
	}
	allDone.Wait()
}
