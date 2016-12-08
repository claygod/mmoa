package service

// Monolithic tools.Message-Oriented Application (MMOA)
// Waiting
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/claygod/mmoa/tools"
)

// Waitings structure
type Waitings struct {
	sync.Mutex
	The   *tools.Themes
	Name  tools.TypeSERVICE
	ChBus chan *tools.Message
	Arr   map[tools.TypeCID]*Waiting
}

func (w *Waitings) Cleaner() {
	c := time.Tick(tools.CleanerTimerSec * time.Nanosecond)
	for _ = range c {
		w.Lock()
		for cid, wg2 := range w.Arr {
			wg := wg2
			if wg != nil && wg.Aggregate.TimeExpire < tools.TypeTIME(time.Now().UnixNano()) {
				w.delWaiting(cid)
			}
		}
		w.Unlock()
		runtime.Gosched()
	}
}

// newWaiting - add a new aggregate
func (w *Waitings) NewWaiting(msg *tools.Message, duration tools.TypeTIME, messages map[string]*tools.Message, ch chan *tools.Message) error {
	w.Lock()
	if _, ok := w.Arr[msg.MsgCid]; ok {
		w.Unlock()
		return errors.New("Aggregate with the cID already exists")
	}
	ae := &Aggregate{
		0,
		tools.TypeTIME(time.Now().UnixNano()),
		tools.TypeTIME(time.Now().UnixNano()) + duration,
		messages,
		len(messages),
		ch,
	}
	wg := &Waiting{
		ae,
		msg,
	}
	w.Arr[msg.MsgCid] = wg
	w.Unlock()
	return nil
}

func (w *Waitings) MsgToWaiting(msg *tools.Message) (chan *tools.Message, *tools.Message) {
	w.Lock()
	wg, ok := w.Arr[msg.MsgCid]
	w.Unlock()
	if wg != nil && ok {
		x := wg.Aggregate.TimeExpire
		if x < tools.TypeTIME(time.Now().UnixNano()) { // the duration of the run-time error
			msgErr := tools.NewMessage().Cid(msg.MsgCid).
				From(w.Name).To(w.The.Service.Trash).
				Theme(w.The.Trash.Timeout).Field(w.The.Attach.Message, msg)
			w.ChBus <- msgErr
		} else if c, err := wg.Aggregate.Add(msg); err != nil { // re-adding the error
			msgErr := tools.NewMessage().Cid(msg.MsgCid).
				From(w.Name).To(w.The.Service.Trash).
				Theme(w.The.Trash.Double).Field(w.The.Attach.Message, msg)
			w.ChBus <- msgErr
		} else if c == 0 { // sent filled aggregate
			wg.Msg.Field(w.The.Attach.Aggregate, wg.Aggregate)
			w.Lock()
			delete(w.Arr, msg.MsgCid)
			w.Unlock()
			return wg.Aggregate.Ch, wg.Msg

		}
	} else { // not aggregate
		msgErr := tools.NewMessage().Cid(msg.MsgCid).
			From(w.Name).To(w.The.Service.Trash).
			Theme(w.The.Trash.Undelivered).Field(w.The.Attach.Message, msg)
		w.ChBus <- msgErr
	}
	return nil, nil
}

func (w *Waitings) delWaiting(cid tools.TypeCID) {
	w.Lock()
	if wg, ok := w.Arr[cid]; ok && wg != nil {
		if wg.Aggregate.Ch != nil {
			wg.Msg.Field(w.The.Attach.Aggregate, wg.Aggregate).
				StatusCode(tools.StatusTimeout)
			wg.Aggregate.Ch <- wg.Msg
		}
		msgErr := tools.NewMessage().Cid(cid).
			From(w.Name).To(w.The.Service.Trash).
			Theme(w.The.Trash.Timeout).Field(w.The.Attach.Aggregate, w.Arr[cid].Aggregate)
		w.ChBus <- msgErr
		delete(w.Arr, cid)
	}
	w.Unlock()
}

// Waiting structure
type Waiting struct {
	Aggregate *Aggregate
	Msg       *tools.Message
}
