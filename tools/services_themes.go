package tools

// Monolithic Message-Oriented Application (MMOA)
// Themes
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import "reflect"

//import "github.com/claygod/mmoa/tools"

// import	"runtime"
//import "fmt"

// import "errors"

// NewThemes - create a new Cid
func NewThemes() *Themes {
	t := &Themes{}
	t.genThemeName()
	return t
}

// Theme structure
type Themes struct {
	Service    ListServices
	Attach     ListAttach
	Article    ThemeArticle
	Menu       ThemeMenu
	Trash      ThemeTrash
	Aggregator ThemeAggregator
}

// ListServices structure
type ListServices struct {
	Controller TypeSERVICE
	Bus        TypeSERVICE
	Service    TypeSERVICE
	Article    TypeSERVICE
	Menu       TypeSERVICE
	Trash      TypeSERVICE
	Aggregator TypeSERVICE
}

// ListAttach structure
type ListAttach struct {
	Message   string
	Aggregate string
	Title     string
	Request   string
}

// ThemeArticle structure
type ThemeArticle struct {
	Record TypeTHEME
	List   TypeTHEME
}

// ThemeMenu structure
type ThemeMenu struct {
	Sitemap TypeTHEME
}

// ThemeTrash structure
type ThemeTrash struct {
	Uncorrect   TypeTHEME
	Undelivered TypeTHEME
	Timeout     TypeTHEME
	Double      TypeTHEME
}

// ThemeAggregator structure
type ThemeAggregator struct {
	Filled TypeTHEME
}

func (t *Themes) genThemeName() {
	t1 := reflect.TypeOf(*t)
	v1 := reflect.ValueOf(t)
	v1 = reflect.Indirect(v1)
	for i := t1.NumField() - 1; i >= 0; i-- {
		t2 := t1.Field(i).Type
		v2 := v1.Field(i).Addr()
		v2 = reflect.Indirect(v2)
		for i2 := t2.NumField() - 1; i2 >= 0; i2-- {
			v3 := v2.Field(i2).Addr()
			v3 = reflect.Indirect(v3)
			v3.SetString(t2.Field(i2).Name)
		}
	}
	//fmt.Print("\n - t.Article.Article=", t.Article.Article)
	//fmt.Print("\n - t.Service.Article=", t.Service.Article)
	//fmt.Print("\n - t.Service.Aggregator=", t.Service.Aggregator)
	//return t
}
