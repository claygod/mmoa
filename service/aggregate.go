package service

// Monolithic tools.Message-Oriented Application (MMOA)
// Aggregate
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/claygod/mmoa/tools"
)

// NewAggregate - create a new Aggregator
func NewAggregate(duration tools.TypeTIME, messages map[string]*tools.Message, ch chan *tools.Message) *Aggregate {
	currentTime := tools.TypeTIME(time.Now().UnixNano())
	return &Aggregate{
		0,
		currentTime,
		currentTime + duration,
		messages,
		len(messages),
		ch,
	}
}

// Aggregate - main struct.
type Aggregate struct {
	Hasp       int32
	TimeCreate tools.TypeTIME
	TimeExpire tools.TypeTIME
	Messages   map[string]*tools.Message
	Counter    int
	Ch         chan *tools.Message
}

// Add - add a message in the aggregate
func (a *Aggregate) Add(msg *tools.Message) (int, error) {
	keyStr := a.GenKey(msg.AddsFrom, msg.MsgTheme)
	if _, ok := a.Messages[keyStr]; !ok {
		return a.Counter, errors.New("This type of aggregate is not provided for in the")
	} else if a.Messages[keyStr] != nil {
		return a.Counter, errors.New("This message is already available")
	}
	a.Messages[keyStr] = msg
	a.Counter--
	return a.Counter, nil
}

// GenKey - generate the key from the service name and the theme
func (a *Aggregate) GenKey(service tools.TypeSERVICE, theme tools.TypeTHEME) string {
	keyStr := fmt.Sprintf(`%s/%s`, string(service), string(theme))
	return keyStr
}

// Lock - mutex analogue lock
func (a *Aggregate) Lock() bool {
	for {
		if a.Hasp == 0 && atomic.CompareAndSwapInt32(&a.Hasp, 0, 1) {
			break
		}
		runtime.Gosched()
	}
	return true
}

// Unlock - mutex analogue unlock
func (a *Aggregate) Unlock() bool {
	for {
		if a.Hasp == 1 && atomic.CompareAndSwapInt32(&a.Hasp, 1, 0) {
			break
		}
		runtime.Gosched()
	}
	return true
}
