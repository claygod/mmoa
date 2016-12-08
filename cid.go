package mmoa

// Monolithic Message-Oriented Application (MMOA)
// Correlation identifier (cID)
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"runtime"
	"sync/atomic"

	"github.com/claygod/mmoa/tools"
)

// NewCid - create a new Cid
func NewCid() *Cid {
	c := &Cid{}
	return c
}

// Cid structure
type Cid struct {
	hasp    int32
	counter tools.TypeCID
}

// Get - get the current count
func (c *Cid) Get() tools.TypeCID {
	c.lock()
	cnt := c.counter
	c.counter++
	c.unlock()
	return cnt
}

// lock - block Cid
func (c *Cid) lock() bool {
	for {
		if c.hasp == 0 && atomic.CompareAndSwapInt32(&c.hasp, 0, 1) {
			break
		}
		runtime.Gosched()
	}
	return true
}

// unlock - block Cid
func (c *Cid) unlock() bool {
	for {
		if c.hasp == 1 && atomic.CompareAndSwapInt32(&c.hasp, 1, 0) {
			break
		}
		runtime.Gosched()
	}
	return true
}
