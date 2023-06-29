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
	"net"
	"sync"
	"time"
)

var _ net.Conn = (*CConn)(nil)

type CConn struct {
	net.Conn
	sync.RWMutex
	unusable bool
	c        *cool
	t        time.Time
}

// Close overrides the net.Conn Close method
// put the connection back to the pool instead of closing it
func (cc *CConn) Close() error {
	cc.RLock()
	defer cc.RUnlock()
	if cc.unusable {
		if cc.Conn != nil {
			return cc.Conn.Close()
		}
		return nil
	}
	return cc.c.put(cc.Conn)
}

func (cc *CConn) MarkUnusable() {
	cc.Lock()
	defer cc.Unlock()
	cc.unusable = true
}

func (cc *CConn) IsUnusable() bool {
	cc.RLock()
	defer cc.RUnlock()
	return cc.unusable
}
