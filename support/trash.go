package support

// Monolithic Message-Oriented Application (MMOA)
// Trash
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "github.com/claygod/mmoa/tools"
import "github.com/claygod/mmoa/service"

// NewTrash - create a new Trash
func NewTrash(chIn chan *tools.Message, chBus chan *tools.Message) *Trash {
	The := tools.NewThemes()
	s := &Trash{
		service.NewService(The.Service.Trash, chIn, chBus),
		service.NewLogger(),
	}
	s.MethodWork = s.toLog
	s.Start()
	return s
}

// Trash structure
type Trash struct {
	service.Service
	logger *service.Logger
}

func (s *Trash) toLog(msg7 *tools.Message) {
	if msg7.AddsFrom == s.The.Service.Aggregator {
		errAgg := msg7.MsgCtx[s.The.Attach.Aggregate].(*service.Aggregate)
		s.logger.Message().
			Field("service", s.Name).
			Field("resurce", s.The.Attach.Aggregate).
			Field("event", msg7.MsgTheme).
			Field("cid", msg7.MsgCid).
			Field("reporter", msg7.AddsFrom).
			Field("attach", errAgg.Messages).
			Send()
	} else {
		errMsg := msg7.MsgCtx[s.The.Attach.Message].(*tools.Message)
		s.logger.Message().
			Field("service", s.Name).
			Field("resurce", s.The.Attach.Message).
			Field("event", msg7.MsgTheme).
			Field("cid", msg7.MsgCid).
			Field("reporter", msg7.AddsFrom).
			Field("attach", errMsg).
			Send()

	}
}
