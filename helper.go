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
	"time"
)

func (v *Violin) submit(wait bool, task func()) {
	if task == nil {
		return
	}
	if wait {
		ctx, done := context.WithCancel(context.Background())
		v.taskChan <- func() {
			task()
			done()
		}
		<-ctx.Done()
	} else {
		v.taskChan <- task
	}
	_ = atomic.AddUint32(&v.taskNum, 1)
}

func (v *Violin) pause(ctx context.Context) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.IsShutdown() {
		return
	}
	wg := new(sync.WaitGroup)
	wg.Add(v.MaxWorkerNum())
	for i := 0; i < v.MaxWorkerNum(); i++ {
		v.Submit(func() {
			wg.Done()
			select {
			case <-ctx.Done():
			case <-v.pauseChan:
			}
		})
	}
	wg.Wait()
}

func (v *Violin) shutdown(wait bool) {
	v.once.Do(func() {
		close(v.taskChan)
		close(v.pauseChan)
		if wait {
			if !atomic.CompareAndSwapUint32(&v.status, statusPlaying, statusCleaning) {
				panic(cleaningFailed)
			}
			v.clean()
			v.waitClose()
		} else {
			v.waitClose()
			if !atomic.CompareAndSwapUint32(&v.status, statusPlaying, statusShutdown) {
				panic(shutdownFailed)
			}
		}
	})
}

func (v *Violin) play() {
	defer func() {
		close(v.shutdownChan)
		if v.IsCleaning() {
			_ = atomic.CompareAndSwapUint32(&v.status, statusCleaning, statusShutdown)
		}
	}()
	if !atomic.CompareAndSwapUint32(&v.status, statusInitialized, statusPlaying) {
		panic(playingFailed)
	}
	wg := new(sync.WaitGroup)
	timer := time.NewTimer(v.options.workerIdleTimeout)
	defer timer.Stop()
LOOP:
	for {
		if int(v.WorkerNum()) < v.MaxWorkerNum() {
			_ = atomic.AddUint32(&v.workerNum, 1)
			go v.recruit(wg)
		}
		select {
		case task, ok := <-v.waitingChan:
			if !ok {
				break LOOP
			}
			select {
			case v.workerChan <- task:
				_ = atomic.AddUint32(&v.waitingTaskNum, ^uint32(0))
			default:
				v.waitingChan <- task
			}
		default:
			select {
			case task, ok := <-v.taskChan:
				if !ok {
					break LOOP
				}
				select {
				case v.workerChan <- task:
				default:
					v.waitingChan <- task
					_ = atomic.AddUint32(&v.waitingTaskNum, 1)
				}
			case <-timer.C:
				if v.WorkerNum() > 0 {
					v.dismiss()
				}
				timer.Reset(v.options.workerIdleTimeout)
			}
		}
	}
	v.dismissAll()
	wg.Wait()
}

func (v *Violin) clean() {
	for v.WaitingTaskNum() > 0 {
		task, ok := <-v.waitingChan
		if !ok {
			break
		}
		_ = atomic.AddUint32(&v.waitingTaskNum, ^uint32(0))
		v.workerChan <- task
	}
}

func (v *Violin) recruit(wg *sync.WaitGroup) {
	wg.Add(1)
	defer func() {
		_ = atomic.AddUint32(&v.workerNum, ^uint32(0))
		wg.Done()
	}()
LOOP:
	for {
		if atomic.LoadUint32(&v.status) == statusShutdown {
			break
		}
		select {
		case <-v.dismissChan:
			break LOOP
		case task, ok := <-v.workerChan:
			if !ok {
				break LOOP
			}
			task()
			_ = atomic.AddUint32(&v.taskNum, ^uint32(0))
		}
	}
}

func (v *Violin) dismiss() {
	select {
	case v.dismissChan <- struct{}{}:
	default:
	}
}

func (v *Violin) dismissAll() {
	for v.WorkerNum() > 0 {
		v.dismiss()
	}
	close(v.dismissChan)
}

func (v *Violin) waitClose() {
	<-v.shutdownChan
	close(v.waitingChan)
	close(v.workerChan)
}
