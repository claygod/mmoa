package tools

// Monolithic Message-Oriented Application (MMOA)
// Reply structures
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

// Article - send in aggregate
type Article struct {
	Id    string
	Title string
	Text  interface{}
}

// List - key-article send in menu
type List map[string]string

// NewPart - create HandlerPart
func NewPart(name TypeSERVICE) *HandlerPart {
	return &HandlerPart{PartName: name}
}

// HandlerPart to simplify the configuration of the handler
type HandlerPart struct {
	PartName     TypeSERVICE
	PartTheme    TypeTHEME
	PartTemplate string
}

// Theme - set theme
func (hp *HandlerPart) Theme(theme TypeTHEME) *HandlerPart {
	hp.PartTheme = theme
	return hp
}

// Template - set template
func (hp *HandlerPart) Template(path string) *HandlerPart {
	hp.PartTemplate = path
	return hp
}
