package tools

// Monolithic Message-Oriented Application (MMOA)
// Message
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "time"

//import "github.com/claygod/mmoa/tools"

// Message - main struct.
type Message struct {
	MsgCid        TypeCID // correlation identifier
	AddsFrom      TypeSERVICE
	AddsTo        TypeSERVICE
	AddsRe        TypeSERVICE
	MsgTheme      TypeTHEME
	MsgCreated    TypeTIME
	MsgCtx        map[string]interface{}
	MsgStatusCode int
}

// NewMessage - create a new Message
func NewMessage() *Message {
	m := &Message{AddsRe: EmptyServiceAddress, MsgCreated: TypeTIME(time.Now().UnixNano()), MsgCtx: make(map[string]interface{})}
	return m
}

func (m *Message) Cid(cid TypeCID) *Message {
	m.MsgCid = cid
	return m
}

func (m *Message) From(from TypeSERVICE) *Message {
	m.AddsFrom = from
	return m
}

func (m *Message) To(to TypeSERVICE) *Message {
	m.AddsTo = to
	return m
}
func (m *Message) Re(re TypeSERVICE) *Message {
	m.AddsRe = re
	return m
}
func (m *Message) Theme(theme TypeTHEME) *Message {
	m.MsgTheme = theme
	return m
}
func (m *Message) Field(key string, value interface{}) *Message {
	m.MsgCtx[key] = value
	return m
}
func (m *Message) StatusCode(h int) *Message {
	m.MsgStatusCode = h
	return m
}
