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

import "time"

type Option func(*options)

type options struct {
	minWorkers        int
	maxWorkers        int
	workerIdleTimeout time.Duration
}

var defaultOptions = options{
	minWorkers:        0,
	maxWorkers:        DefaultViolinMaxWorkerNum,
	workerIdleTimeout: DefaultViolinWorkerIdleTimeout,
}

func newOptions(opts ...Option) *options {
	options := &options{
		minWorkers:        defaultOptions.minWorkers,
		maxWorkers:        defaultOptions.maxWorkers,
		workerIdleTimeout: defaultOptions.workerIdleTimeout,
	}
	options.apply(opts...)
	return options
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithMinWorkers set the minimum number of workers
func WithMinWorkers(min int) Option {
	if min < 0 {
		min = 0
	}
	return func(o *options) {
		o.minWorkers = min
	}
}

// WithMaxWorkers set the maximum number of workers
func WithMaxWorkers(max int) Option {
	if max < 1 {
		max = 1
	}
	return func(o *options) {
		o.maxWorkers = max
	}
}

// WithWorkerIdleTimeout set the destroyed timeout of idle workers
func WithWorkerIdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.workerIdleTimeout = timeout
	}
}
