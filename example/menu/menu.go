package menu

// Monolithic Message-Oriented Application (MMOA)
// Service Menu
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"github.com/claygod/mmoa/service"
	"github.com/claygod/mmoa/tools"
)

// NewServiceMenu - create a new ServiceMenu
func NewServiceMenu(chIn chan *tools.Message, chBus chan *tools.Message) *ServiceMenu {
	The := tools.NewThemes()
	s := &ServiceMenu{
		service.NewService(The.Service.Menu, chIn, chBus),
	}
	s.MethodWork = s.Work
	s.setEvents()
	s.Start()
	return s
}

// ServiceMenu structure
type ServiceMenu struct {
	service.Service
}

func (s *ServiceMenu) setEvents() {
	s.Methods[s.The.Menu.Sitemap] = s.sitemapEvent
}

func (s *ServiceMenu) sitemapEvent(msgIn *tools.Message) {
	if _, ok := msgIn.MsgCtx[s.The.Attach.Aggregate]; ok {
		s.sitemapStep2(msgIn)
	} else {
		s.sitemapStep1(msgIn)
	}
}

func (s *ServiceMenu) sitemapStep1(msgIn *tools.Message) {
	a := &service.Aggregate{}
	key := a.GenKey(s.The.Service.Article, s.The.Article.List)
	messages := map[string]*tools.Message{key: nil}
	s.WaitingFor.NewWaiting(msgIn, tools.DurationHandle, messages, nil)
	// inquiry in article for obtaining `list`
	msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
		From(s.Name).To(s.The.Service.Article).Re(s.Name).
		Theme(s.The.Article.List)
	s.ChBus <- msgOut
}

func (s *ServiceMenu) sitemapStep2(msgIn *tools.Message) {
	ae := msgIn.MsgCtx[s.The.Attach.Aggregate].(*service.Aggregate) // take aggregate
	keyStr := ae.GenKey(s.The.Service.Article, s.The.Article.List)
	if m, ok := ae.Messages[keyStr]; ok { // take message
		if list, ok := m.MsgCtx[string(s.The.Article.List)]; ok { // take list
			list = map[string]string(list.(tools.List))
			msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
				From(s.Name).To(msgIn.AddsRe).
				Theme(msgIn.MsgTheme).StatusCode(tools.StatusOK).
				Field("Menu", list)
			s.ChBus <- msgOut
		}
	}
}
