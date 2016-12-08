package tools

// Monolithic Message-Oriented Application (MMOA)
// Reply structures
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

type Article struct {
	Id    string
	Title string
	Text  interface{}
}

type List map[string]string

func NewPart(name TypeSERVICE) *HandlerPart {
	return &HandlerPart{PartName: name}
}

type HandlerPart struct {
	PartName     TypeSERVICE
	PartTheme    TypeTHEME
	PartTemplate string
}

func (hp *HandlerPart) Theme(theme TypeTHEME) *HandlerPart {
	hp.PartTheme = theme
	return hp
}

func (hp *HandlerPart) Template(path string) *HandlerPart {
	hp.PartTemplate = path
	return hp
}
