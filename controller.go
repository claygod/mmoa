package mmoa

// Monolithic Message-Oriented Application (MMOA)
// Controller
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"github.com/claygod/mmoa/support"
	"github.com/claygod/mmoa/tools"
)

// NewController - create a new Controller
func NewController(chBus chan *tools.Message) *Controller {
	c := &Controller{
		cid:   NewCid(),
		the:   tools.NewThemes(),
		chBus: chBus,
	}
	c.addBus()
	c.addAggregator()
	c.addTrash()
	return c
}

// Controller - main struct.
type Controller struct {
	the        *tools.Themes
	cid        *Cid
	bus        *support.Bus
	aggregator *support.Aggregator
	chBus      chan *tools.Message
}

func (c *Controller) Handler(template string) *Handler {
	h := NewHandler(template, c.chBus, c.aggregator, c.cid)
	return h
}

func (c *Controller) AddService(nameService tools.TypeSERVICE, chIn chan *tools.Message) {
	c.bus.Set(nameService, chIn)
}

func (c *Controller) addBus() {
	c.bus = support.NewBus(c.chBus)
}

func (c *Controller) addAggregator() {
	chAggregator := make(chan *tools.Message, 1000)
	c.aggregator = support.NewAggregator(chAggregator, c.chBus)
	c.bus.Set(c.the.Service.Aggregator, chAggregator)
}

func (c *Controller) addTrash() {
	chTrash := make(chan *tools.Message, 1000)
	support.NewTrash(chTrash, c.chBus)
	c.bus.Set(c.the.Service.Trash, chTrash)
}
