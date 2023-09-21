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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	options := newOptions(
		WithMinWorkers(5),
		WithMaxWorkers(13),
		WithWorkerIdleTimeout(time.Second*10),
	)
	assert.Equal(t, 5, options.minWorkers)
	assert.Equal(t, 13, options.maxWorkers)
	assert.Equal(t, time.Second*10, options.workerIdleTimeout)
}

func TestDefaultOptions(t *testing.T) {
	options := newOptions()
	assert.Equal(t, 0, options.minWorkers)
	assert.Equal(t, 5, options.maxWorkers)
	assert.Equal(t, time.Second*3, options.workerIdleTimeout)
}
