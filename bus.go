package support

// Monolithic Message-Oriented Application (MMOA)
// Bus
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "sync"
import "github.com/claygod/mmoa/tools"

//import "fmt"

//import "time"

//import "sync/atomic"

// NewBus - create a new Bus
func NewBus(ch chan *tools.Message) *Bus {
	b := &Bus{
		sync.Mutex{},
		tools.NewThemes(),
		ch,
		make(map[tools.TypeSERVICE]chan *tools.Message),
		1,
	}
	go b.Worker()
	return b
}

// Bus structure
type Bus struct {
	sync.Mutex
	the        *tools.Themes
	chIn       chan *tools.Message
	services   map[tools.TypeSERVICE]chan *tools.Message
	numWorkers int32
}

func (b *Bus) Set(name tools.TypeSERVICE, ChIn chan *tools.Message) {
	b.services[name] = ChIn
}

func (b *Bus) Del(name tools.TypeSERVICE) {
	delete(b.services, name)
}

func (b *Bus) Worker() {
	for {
		msgIn := <-b.chIn

		if s, ok := b.services[msgIn.AddsTo]; ok {
			s <- msgIn
		} else {
			// response to messages from the wrong location (only if there is an address RE)
			if msgIn.AddsRe != tools.EmptyServiceAddress {
				msgErr := tools.NewMessage().Cid(msgIn.MsgCid).
					From(b.the.Service.Bus).To(msgIn.AddsRe).
					Theme(b.the.Trash.Undelivered).Field(b.the.Attach.Message, msgIn)
				b.services[msgIn.AddsRe] <- msgErr
			}
			// in trash
			if trash, ok := b.services[b.the.Service.Trash]; ok {
				msgErr := tools.NewMessage().Cid(msgIn.MsgCid).
					From(b.the.Service.Bus).To(msgIn.AddsFrom).
					Theme(b.the.Trash.Undelivered).Field(b.the.Attach.Message, msgIn)
				trash <- msgErr
			}
		}
	}
}
