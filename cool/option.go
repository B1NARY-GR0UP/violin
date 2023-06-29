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

import "time"

type Option func(*options)

type options struct {
	connIdleTimeout time.Duration
}

func newOptions(opts ...Option) *options {
	options := &options{}
	options.apply(opts...)
	return options
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// WithConnIdleTimeout will set the connection idle timeout
func WithConnIdleTimeout(timeout time.Duration) Option {
	return func(o *options) {
		o.connIdleTimeout = timeout
	}
}
