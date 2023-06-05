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
)

type Violin struct {
	options *options

	mu   sync.Mutex
	once sync.Once

	workerNum      uint32
	taskNum        uint32
	waitingTaskNum uint32
	status         uint32

	workerChan   chan func()
	taskChan     chan func()
	waitingChan  chan func()
	dismissChan  chan struct{}
	pauseChan    chan struct{}
	shutdownChan chan struct{}
}

const (
	_ uint32 = iota
	statusInitialized
	statusPlaying
	statusCleaning
	statusShutdown
)

const (
	initializedFailed = "violin initialize failed"
	playingFailed     = "violin playing failed"
	cleaningFailed    = "violin cleaning failed"
	shutdownFailed    = "violin shutdown failed"
)

// New VIOLIN worker pool
func New(opts ...Option) *Violin {
	options := newOptions(opts...)
	v := &Violin{
		options:      options,
		workerChan:   make(chan func()),
		taskChan:     make(chan func()),
		waitingChan:  make(chan func(), options.waitingQueueSize),
		dismissChan:  make(chan struct{}),
		pauseChan:    make(chan struct{}),
		shutdownChan: make(chan struct{}),
	}
	if !atomic.CompareAndSwapUint32(&v.status, 0, statusInitialized) {
		panic(initializedFailed)
	}
	go v.play()
	return v
}

// Submit a task to the worker pool
func (v *Violin) Submit(task func()) {
	v.submit(false, task)
}

// SubmitWait a task to the worker pool and wait for it to complete
func (v *Violin) SubmitWait(task func()) {
	v.submit(true, task)
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

// MaxWorkerNum returns the maximum number of workers
func (v *Violin) MaxWorkerNum() int {
	return v.options.maxWorkers
}

// WaitingQueueSize returns the size of the waiting queue
func (v *Violin) WaitingQueueSize() int {
	return v.options.waitingQueueSize
}

// TaskNum returns the number of tasks
func (v *Violin) TaskNum() uint32 {
	return atomic.LoadUint32(&v.taskNum)
}

// WaitingTaskNum returns the number of waiting tasks
func (v *Violin) WaitingTaskNum() uint32 {
	return atomic.LoadUint32(&v.waitingTaskNum)
}

// WorkerNum returns the number of workers
func (v *Violin) WorkerNum() uint32 {
	return atomic.LoadUint32(&v.workerNum)
}
