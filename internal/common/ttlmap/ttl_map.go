// Copyright The OpenTelemetry Authors
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

package ttlmap // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/common/ttlmap"

import (
	"sync"
	"time"
)

// TTLMap is a map that evicts entries after the configured ttl has elapsed.
type TTLMap struct {
	md            *ttlMapData
	sweepInterval int64
}

// New creates a TTLMap. The sweepIntervalSeconds arg indicates how often
// entries are checked for expiration. The maxAgeSeconds arg indicates how long
// entries can persist before getting evicted. Call Start() on the returned
// TTLMap to begin periodic sweeps which check for expiration and evict entries
// as needed.
func New(sweepIntervalSeconds int64, maxAgeSeconds int64) *TTLMap {
	return &TTLMap{
		sweepInterval: sweepIntervalSeconds,
		md:            newTTLMapData(maxAgeSeconds),
	}
}

// Start starts periodic sweeps for expired entries in the underlying map.
func (m *TTLMap) Start() {
	go func() {
		d := time.Duration(m.sweepInterval) * time.Second
		for now := range time.Tick(d) {
			m.md.sweep(now.Unix())
		}
	}()
}

// Put adds the passed-in key and value to the underlying map. The current time
// is attached to the entry for periodic expiration checking and eviction when
// necessary.
func (m *TTLMap) Put(k string, v interface{}) {
	m.md.put(k, v, time.Now().Unix())
}

// Get returns the object in the underlying map at the given key. If there is no
// value at that key, Get returns nil.
func (m *TTLMap) Get(k string) interface{} {
	return m.md.get(k)
}

type entry struct {
	createTime int64
	v          interface{}
}

type ttlMapData struct {
	m      map[string]entry
	maxAge int64
	mux    sync.Mutex
}

func newTTLMapData(maxAgeSeconds int64) *ttlMapData {
	return &ttlMapData{
		maxAge: maxAgeSeconds,
		m:      map[string]entry{},
		mux:    sync.Mutex{},
	}
}

func (d *ttlMapData) put(k string, v interface{}, currTime int64) {
	d.mux.Lock()
	d.m[k] = entry{v: v, createTime: currTime}
	d.mux.Unlock()
}

func (d *ttlMapData) get(k string) interface{} {
	d.mux.Lock()
	defer d.mux.Unlock()
	entry, ok := d.m[k]
	if !ok {
		return nil
	}
	return entry.v
}

func (d *ttlMapData) sweep(currTime int64) {
	d.mux.Lock()
	for k, v := range d.m {
		if currTime-v.createTime > d.maxAge {
			delete(d.m, k)
		}
	}
	d.mux.Unlock()
}
