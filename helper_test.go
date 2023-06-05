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

//func TestCloseChan(t *testing.T) {
//	ch := make(chan int)
//	go func() {
//		select {
//		case num, ok := <-ch:
//			fmt.Println(num, ok)
//		}
//	}()
//	time.Sleep(time.Second)
//	close(ch)
//}
//
//func TestCancel(t *testing.T) {
//	ctx, done := context.WithCancel(context.Background())
//	go func() {
//		time.Sleep(time.Second * 3)
//		done()
//	}()
//	select {
//	case <-ctx.Done():
//		fmt.Println("done")
//	}
//}
//
//func TestLoop(t *testing.T) {
//	ctx, done := context.WithCancel(context.Background())
//	go func() {
//		time.Sleep(time.Second * 5)
//		done()
//	}()
//Loop:
//	for {
//		select {
//		case <-ctx.Done():
//			fmt.Println("done")
//			// In for-select, if you break in select, you can only exit select, not for
//			break Loop
//		default:
//			time.Sleep(time.Second * 1)
//			fmt.Println("loop")
//		}
//	}
//	fmt.Println("end")
//}
//
//func TestMaxLimit(t *testing.T) {
//	v := New(WithMaxWorkers(13))
//	defer v.Shutdown()
//	go func() {
//		for {
//			fmt.Println("Worker Number: ", v.WorkerNum())
//		}
//	}()
//	for {
//		v.Submit(func() {
//			fmt.Println("Hello, VIOLIN!")
//		})
//	}
//}
