package service

// Monolithic tools.Message-Oriented Application (MMOA)
// Service
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"sync"

	"github.com/claygod/mmoa/tools"
)

//import "fmt"

// NewService - create a new Service
func NewService(nameServiceMenu tools.TypeSERVICE, chIn chan *tools.Message, chBus chan *tools.Message) Service {
	wts := &Waitings{
		Arr:   make(map[tools.TypeCID]*Waiting),
		Name:  nameServiceMenu,
		ChBus: chBus,
		The:   tools.NewThemes(),
	}
	return Service{
		sync.Mutex{},
		tools.NewThemes(),
		nameServiceMenu,
		chIn,
		chBus,
		make(map[tools.TypeTHEME]func(*tools.Message)),
		nil,
		wts,
		NewLogger(),
	}
}

// Service structure
type Service struct {
	sync.Mutex
	The        *tools.Themes
	Name       tools.TypeSERVICE
	ChIn       chan *tools.Message
	ChBus      chan *tools.Message
	Methods    map[tools.TypeTHEME]func(*tools.Message)
	MethodWork func(*tools.Message)
	WaitingFor *Waitings
	Logger     *Logger
}

func (s *Service) Start() {
	go s.worker()
	go s.WaitingFor.Cleaner()
}

func (s *Service) worker() {
	for {
		msg := <-s.ChIn
		//s.Lock()
		s.WaitingFor.Lock()
		//defer s.WaitingFor.Unlock()
		wg, ok := s.WaitingFor.Arr[msg.MsgCid]
		s.WaitingFor.Unlock()
		//s.Unlock()
		if ok && wg != nil {
			s.addMessage(msg)
		} else {
			s.MethodWork(msg)
		}
	}
}

func (s *Service) Work(msg *tools.Message) {
	if m, ok := s.Methods[msg.MsgTheme]; ok {
		go m(msg)
	} else {
		msgErr := tools.NewMessage().Cid(msg.MsgCid).
			From(s.Name).To(s.The.Service.Trash).
			Theme(s.The.Trash.Uncorrect).Field(s.The.Attach.Message, msg)
		s.ChBus <- msgErr
	}
}

func (s *Service) addMessage(msg *tools.Message) {
	if ch, mes := s.WaitingFor.MsgToWaiting(msg); mes != nil {
		if ch != nil {
			ch <- mes
		} else if mtd, ok := s.Methods[mes.MsgTheme]; ok {
			mtd(mes)
		}
	}
}
