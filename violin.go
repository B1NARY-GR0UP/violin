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
	"context"
	"sync"
	"sync/atomic"

	"github.com/B1NARY-GR0UP/violin/ring"
)

// Violin VIOLIN worker pool
type Violin struct {
	options *options

	mu   sync.RWMutex
	once sync.Once

	workerNum uint32
	taskNum   uint32
	status    uint32

	waitingQ *ring.Ring[func()]

	workerC   chan func()
	taskC     chan func()
	dismissC  chan struct{}
	pauseC    chan struct{}
	shutdownC chan struct{}
}

const (
	_ uint32 = iota
	statusInitialized
	statusPlaying
	statusCleaning
	statusShutdown
)

// New VIOLIN worker pool
// TODO: improve performance
func New(opts ...Option) *Violin {
	options := newOptions(opts...)
	v := &Violin{
		options:   options,
		waitingQ:  ring.New[func()](),
		workerC:   make(chan func()),
		taskC:     make(chan func()),
		dismissC:  make(chan struct{}),
		pauseC:    make(chan struct{}),
		shutdownC: make(chan struct{}),
	}
	_ = atomic.CompareAndSwapUint32(&v.status, 0, statusInitialized)
	go v.play()
	return v
}

// Submit a task to the worker pool
func (v *Violin) Submit(task func()) {
	v.submit(false, task)
}

// SubmitWait submit a task to the worker pool and wait for it to complete
func (v *Violin) SubmitWait(task func()) {
	v.submit(true, task)
}

// Consume the current tasks in the taskC
// Note: Consume will not execute the tasks which put into the taskC after calling Consume
func (v *Violin) Consume(taskC chan func()) {
	v.consume(false, taskC)
}

// ConsumeWait consume tasks in the taskC and wait for them to complete
// Note: ConsumeWait will not execute the tasks which put into the taskC after calling ConsumeWait
func (v *Violin) ConsumeWait(taskC chan func()) {
	v.consume(true, taskC)
}

// Pause the worker pool
func (v *Violin) Pause(ctx context.Context) {
	v.pause(ctx)
}

// Shutdown the worker pool
func (v *Violin) Shutdown() {
	v.shutdown(false)
}

// ShutdownWait graceful shutdown, wait for all tasks to complete
func (v *Violin) ShutdownWait() {
	v.shutdown(true)
}

// IsPlaying returns true if the worker pool is playing (running)
func (v *Violin) IsPlaying() bool {
	return atomic.LoadUint32(&v.status) == statusPlaying
}

// IsCleaning returns true if the worker pool is cleaning (graceful shutdown)
func (v *Violin) IsCleaning() bool {
	return atomic.LoadUint32(&v.status) == statusCleaning
}

// IsShutdown returns true if the worker pool is shutdown
func (v *Violin) IsShutdown() bool {
	return atomic.LoadUint32(&v.status) == statusShutdown
}

// MinWorkerNum returns the minimum number of workers
func (v *Violin) MinWorkerNum() int {
	return v.options.minWorkers
}

// MaxWorkerNum returns the maximum number of workers
func (v *Violin) MaxWorkerNum() int {
	return v.options.maxWorkers
}

// TaskNum returns the number of tasks
func (v *Violin) TaskNum() uint32 {
	return atomic.LoadUint32(&v.taskNum)
}

// WaitingTaskNum returns the number of waiting tasks
func (v *Violin) WaitingTaskNum() int {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.waitingQ.Size()
}

// WorkerNum returns the number of workers
func (v *Violin) WorkerNum() uint32 {
	return atomic.LoadUint32(&v.workerNum)
}
