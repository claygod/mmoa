package tools

// Monolithic Message-Oriented Application (MMOA)
// Message
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "time"

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

// Cid - correlation identifier
func (m *Message) Cid(cid TypeCID) *Message {
	m.MsgCid = cid
	return m
}

// From - sende
func (m *Message) From(from TypeSERVICE) *Message {
	m.AddsFrom = from
	return m
}

// To - destination
func (m *Message) To(to TypeSERVICE) *Message {
	m.AddsTo = to
	return m
}

// Re - where to send the response
func (m *Message) Re(re TypeSERVICE) *Message {
	m.AddsRe = re
	return m
}

// Theme - the meggage theme
func (m *Message) Theme(theme TypeTHEME) *Message {
	m.MsgTheme = theme
	return m
}

// Field - the contents of the letter
func (m *Message) Field(key string, value interface{}) *Message {
	m.MsgCtx[key] = value
	return m
}

// StatusCode - response status
func (m *Message) StatusCode(h int) *Message {
	m.MsgStatusCode = h
	return m
}
