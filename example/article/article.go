package article

// Monolithic Message-Oriented Application (MMOA)
// Service Article
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/claygod/mmoa/service"
	"github.com/claygod/mmoa/tools"
)

// import "fmt"

// NewServiceArticle - create a new ServiceArticle
func NewServiceArticle(chIn chan *tools.Message, chBus chan *tools.Message) *ServiceArticle {
	the := tools.NewThemes()
	s := &ServiceArticle{
		service.NewService(the.Service.Article, chIn, chBus),
		make(map[string]tools.Article),
	}
	s.MethodWork = s.Work
	s.initService(dbPath)
	s.setEvents()
	s.Start()
	return s
}

// ServiceArticle structure
type ServiceArticle struct {
	service.Service
	articles map[string]tools.Article
}

func (s *ServiceArticle) setEvents() {
	s.Methods[s.The.Article.Record] = s.recordEvent
	s.Methods[s.The.Article.List] = s.listEvent
}

func (s *ServiceArticle) recordEvent(msgIn *tools.Message) {
	if msgIn.AddsRe != tools.EmptyServiceAddress {
		id := s.getId(msgIn)
		if id != "" {
			if a, ok := s.articles[id]; ok {
				msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
					From(s.Name).To(msgIn.AddsRe).
					Theme(msgIn.MsgTheme).StatusCode(tools.StatusOK).
					Field("Title", a.Title).
					Field("Text", a.Text)
				s.ChBus <- msgOut
				return
			}
		}
	}
	msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
		From(s.Name).To(msgIn.AddsRe).
		Theme(msgIn.MsgTheme).StatusCode(tools.StatusNotFound).
		Field("Title", "404").
		Field("Text", "Not found")
	s.ChBus <- msgOut
}

func (s *ServiceArticle) getId(msgIn *tools.Message) string {
	var id string
	req := msgIn.MsgCtx[s.The.Attach.Request].(*http.Request)
	ctx := req.Context()

	idi := ctx.Value("id")
	if req.URL.String() == "/" {
		id = "/"
	} else if idi != nil {
		id = idi.(string)
	}
	return id
}

func (s *ServiceArticle) listEvent(msgIn *tools.Message) {
	listArticles := make(tools.List)
	for k, v := range s.articles {
		listArticles[k] = v.Title
	}
	msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
		From(s.Name).To(msgIn.AddsRe).
		Theme(msgIn.MsgTheme).StatusCode(tools.StatusOK).
		Field(string(s.The.Article.List), listArticles)
	s.ChBus <- msgOut
}

func (s *ServiceArticle) initService(path string) {
	dbBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	dbStringArr := strings.Split(string(dbBytes), "\r\n")
	for _, str := range dbStringArr {
		strArr := strings.Split(str, `|`)
		ra := tools.Article{
			Id:    strArr[0],
			Title: strArr[1],
			Text:  template.HTML(strArr[2]),
		}
		s.articles[strArr[0]] = ra
	}
}
