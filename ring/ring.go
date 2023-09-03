// Copyright 2023 BINARY Members
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package ring

import (
	"sync"
)

type (
	// Ring deque
	Ring[T any] struct {
		mu       sync.RWMutex
		sentinel *node[T]
		size     int
	}
	node[T any] struct {
		prev, next *node[T]
		item       T
	}
)

func New[T any]() *Ring[T] {
	r := &Ring[T]{
		sentinel: &node[T]{},
	}
	r.sentinel.prev = r.sentinel
	r.sentinel.next = r.sentinel
	return r
}

func NewWithCopy[T any](r *Ring[T]) *Ring[T] {
	nr := New[T]()
	for i := 0; i < r.size; i++ {
		nr.AddLast(r.IterGet(i))
	}
	return nr
}

func (r *Ring[T]) Size() int {
	return r.size
}

func (r *Ring[T]) IsEmpty() bool {
	return r.size == 0
}

func (r *Ring[T]) AddLast(item T) {
	n := &node[T]{
		prev: r.sentinel.prev,
		next: r.sentinel,
		item: item,
	}
	r.sentinel.prev.next = n
	r.sentinel.prev = n
	r.size++
}

func (r *Ring[T]) AddFirst(item T) {
	n := &node[T]{
		prev: r.sentinel,
		next: r.sentinel.next,
		item: item,
	}
	r.sentinel.next.prev = n
	r.sentinel.next = n
	r.size++
}

// DelFirst remove the first node and return its value
func (r *Ring[T]) DelFirst() T {
	if r.size == 0 {
		return *new(T)
	}
	n := r.sentinel.next
	n.prev.next = n.next
	n.next.prev = n.prev
	r.size--
	return n.item
}

// DelLast remove the last node and return its value
func (r *Ring[T]) DelLast() T {
	if r.size == 0 {
		return *new(T)
	}
	n := r.sentinel.prev
	n.prev.next = n.next
	n.next.prev = n.prev
	r.size--
	return n.item
}

func (r *Ring[T]) IterGet(index int) T {
	if index < 0 || index >= r.size {
		return *new(T)
	}
	head := r.sentinel.next
	for i := 0; i < index; i++ {
		head = head.next
	}
	return head.item
}

func (r *Ring[T]) RecGet(index int) T {
	if index < 0 || index >= r.size {
		return *new(T)
	}
	return r.recGetHelper(r.sentinel.next, index)
}

func (r *Ring[T]) recGetHelper(head *node[T], index int) T {
	if index == 0 {
		return head.item
	}
	return r.recGetHelper(head.next, index-1)
}
