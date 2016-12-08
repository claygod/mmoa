package mmoa

// Monolithic Message-Oriented Application (MMOA)
// View
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/claygod/mmoa/service"
	"github.com/claygod/mmoa/tools"
)

// NewView - create a new View
func NewView() *View {
	v := &View{
		the:       tools.NewThemes(),
		templates: make(map[tools.TypeSERVICE]map[tools.TypeTHEME]*template.Template),
	}
	return v
}

// View structure
type View struct {
	the          *tools.Themes
	tpl          *template.Template
	templates    map[tools.TypeSERVICE]map[tools.TypeTHEME]*template.Template
	statusCodeOf tools.TypeTHEME
	contentType  string
}

// ContentType - for output
func (v *View) ContentType(ct string) {
	v.contentType = ct
}

// StatusCodeOf - the answer to what theme will form status
func (v *View) StatusCodeOf(theme tools.TypeTHEME) {
	v.statusCodeOf = theme
}

// TemplatePage - Loading of the page template
func (v *View) TemplatePage(path string) {
	t, _ := ioutil.ReadFile(path)
	tpl, _ := template.New(path).Parse(string(t))
	v.tpl = tpl
}

// TemplateService - Loading of the service template
func (v *View) TemplateService(service tools.TypeSERVICE, theme tools.TypeTHEME, path string) {
	t, _ := ioutil.ReadFile(path)
	keyStr := string(theme)
	tpl, _ := template.New(keyStr).Parse(string(t))
	if _, ok := v.templates[service]; !ok {
		v.templates[service] = make(map[tools.TypeTHEME]*template.Template)
	}
	v.templates[service][theme] = tpl
}

// ProcessingAggregate - processing the resulting aggregate, messages are displayed on templates
func (v *View) ProcessingAggregate(messages map[string]*tools.Message, statusCode int) (map[string]template.HTML, int) {
	var sCode int = tools.StatusNotFound
	var title bytes.Buffer
	a := &service.Aggregate{}
	arr := make(map[string]template.HTML)
	for service, vw := range v.templates {
		for theme, tpl := range vw {
			keyStr := a.GenKey(service, theme)
			if msg, ok := messages[keyStr]; ok && msg != nil {
				if theme == v.statusCodeOf {
					sCode = msg.MsgStatusCode
					if ttl, ok := msg.MsgCtx[v.the.Attach.Title]; ok {
						title.WriteString(ttl.(string))
						title.WriteString(" ")
					}
				}
				var doc bytes.Buffer
				tpl.Execute(&doc, msg.MsgCtx)
				arr[tpl.Name()] = template.HTML(doc.String())
			} else {
				if theme == v.statusCodeOf {
					sCode = statusCode
					title.WriteString(string(v.the.Trash.Timeout))
				}
			}
		}
	}
	arr[v.the.Attach.Title] = template.HTML(strings.TrimSpace(title.String()))
	return arr, sCode
}
