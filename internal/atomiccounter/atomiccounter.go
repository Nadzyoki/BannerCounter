package atomiccounter

import (
	"sync"
	"sync/atomic"
)

type AtomicCounter struct {
	mu       sync.RWMutex
	counters map[string]*uint64
}

func NewAtomicCounter() *AtomicCounter {
	return &AtomicCounter{
		counters: make(map[string]*uint64),
	}
}

func (c *AtomicCounter) Add(id string) {
	c.mu.RLock()
	if cntPtr, exists := c.counters[id]; exists {
		c.mu.RUnlock()
		atomic.AddUint64(cntPtr, 1)
		return
	}
	c.mu.RUnlock()

	c.mu.Lock()
	if cntPtr, exists := c.counters[id]; exists {
		atomic.AddUint64(cntPtr, 1)
		c.mu.Unlock()
		return
	}
	
	cntPtr := new(uint64)
	c.counters[id] = cntPtr
	atomic.AddUint64(cntPtr, 1)
	c.mu.Unlock()
}

func (c *AtomicCounter) Get(id string) uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if cntPtr, exists := c.counters[id]; exists {
		return atomic.LoadUint64(cntPtr)
	}
	return 0
}

func (c *AtomicCounter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k := range c.counters {
		delete(c.counters, k)
	}
}

func (c *AtomicCounter) GetAndReset() map[string]uint64 {
	c.mu.Lock()
	defer c.mu.Unlock()

	res := make(map[string]uint64, len(c.counters))
	for k, v := range c.counters {
		res[k] = atomic.LoadUint64(v)
	}
	c.counters = make(map[string]*uint64)
	return res
}
