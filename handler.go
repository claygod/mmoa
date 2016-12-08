package mmoa

// Monolithic Message-Oriented Application (MMOA)
// Handler
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"html/template"
	"io"
	"net/http"

	"github.com/claygod/mmoa/service"
	"github.com/claygod/mmoa/support"
	"github.com/claygod/mmoa/tools"
)

// NewHandler - create a new Handler
func NewHandler(path string, chBus chan *tools.Message, aggregator *support.Aggregator, cid *Cid) *Handler {
	h := &Handler{
		the:        tools.NewThemes(),
		cid:        cid,
		parts:      make([]*tools.HandlerPart, 0),
		chBus:      chBus,
		aggregator: aggregator,
		templates:  make(map[tools.TypeSERVICE]map[tools.TypeTHEME]*template.Template),
		view:       NewView(),
	}
	h.view.TemplatePage(path)
	return h
}

// Handler structure
type Handler struct {
	the          *tools.Themes
	cid          *Cid
	tpl          *template.Template
	parts        []*tools.HandlerPart
	chBus        chan *tools.Message
	aggregator   *support.Aggregator
	templates    map[tools.TypeSERVICE]map[tools.TypeTHEME]*template.Template
	statusCodeOf tools.TypeTHEME
	contentType  string
	view         *View
}

// StatusCodeOf - the answer to what theme will form status
func (h *Handler) StatusCodeOf(theme tools.TypeTHEME) *Handler {
	h.view.StatusCodeOf(theme)
	return h
}

// Service - add the service to the handler
func (h *Handler) Service(hp *tools.HandlerPart) *Handler {
	h.view.TemplateService(hp.PartName, hp.PartTheme, hp.PartTemplate)
	h.parts = append(h.parts, hp)
	return h
}

// ContentType - for output
func (h *Handler) ContentType(ct string) *Handler {
	h.view.ContentType(ct)
	return h
}

// Do - launch handler
func (h *Handler) Do(w http.ResponseWriter, req *http.Request) {
	w.Header().Del("Content-Type")
	w.Header().Set("Content-Type", h.contentType)
	ch := h.handlePush(req)

	msgAgg := <-ch
	var agg *service.Aggregate
	if ag, ok := msgAgg.MsgCtx[h.the.Attach.Aggregate]; ok {
		agg = ag.(*service.Aggregate)
	} else {
		w.WriteHeader(404)
		io.WriteString(w, "404 Page not found")
		return
	}
	arr, statusCode := h.view.ProcessingAggregate(agg.Messages, msgAgg.MsgStatusCode)
	w.WriteHeader(statusCode)
	h.view.tpl.Execute(w, arr)
}

func (h *Handler) handlePush(req *http.Request) chan *tools.Message {
	a := &service.Aggregate{}
	var msg *tools.Message
	cid := h.cid.Get()
	ch := make(chan *tools.Message, len(h.parts))
	messages := make(map[string]*tools.Message)
	for _, p := range h.parts {
		key := a.GenKey(p.PartName, p.PartTheme)
		messages[key] = nil
	}
	h.aggregator.Aggregate(cid, tools.DurationHandle, messages, ch)
	for _, p := range h.parts {
		msg = tools.NewMessage().Cid(cid).
			From(h.the.Service.Controller).To(p.PartName).Re(h.the.Service.Aggregator).
			Field(h.the.Attach.Request, req).
			Theme(p.PartTheme)
		h.chBus <- msg
	}
	return ch
}
