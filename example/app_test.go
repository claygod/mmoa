package main

// Monolithic Message-Oriented Application (MMOA)
// Application
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/claygod/Bxog"
	"github.com/claygod/mmoa"
	"github.com/claygod/mmoa/example/article"
	"github.com/claygod/mmoa/example/menu"
	"github.com/claygod/mmoa/tools"
)

const (
	pathTemplate string = "./data/template.html"
	pathSitemap  string = "./data/sitemap.html"
	pathRecord   string = "./data/record.html"
	contentType  string = "./text/html; charset=UTF-8"
)

func TestApp200(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler(pathTemplate).
		ContentType(contentType).
		Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template(pathSitemap)).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/", hPage.Do)
	m.Test()
	m.ServeHTTP(res, req)

	if res.Code != 200 {
		t.Error("Wrong code, I expect 200 and received: ", res.Code, res.Body)
	}
}

func TestApp404(t *testing.T) {
	req, _ := http.NewRequest("GET", "/xx", nil)
	res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler(pathTemplate).
		ContentType(contentType).
		Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template(pathSitemap)).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Test()
	m.ServeHTTP(res, req)

	if res.Code != 404 {
		t.Error("Wrong code, I expect 404 and received: ", res.Code)
	}
}

func BenchmarkOnlyArticle(b *testing.B) { // BenchmarkOnlyArticle-4      	   50000	     28822 ns/op
	b.StopTimer()
	req, _ := http.NewRequest("GET", "/", nil)
	//res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler("template.html").
		ContentType(contentType).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Test()
	// m.ServeHTTP(res, req)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
	}
}

func BenchmarkArticleAndMenu(b *testing.B) { // BenchmarkApp-4   	   20000	     96359 ns/op
	b.StopTimer()
	req, _ := http.NewRequest("GET", "/", nil)
	//res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler("template.html").
		ContentType(contentType).
		Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template("sitemap.html")).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Test()
	// m.ServeHTTP(res, req)

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
	}
}

func BenchmarkOnlyArticleParallel(b *testing.B) {
	b.StopTimer()
	req, _ := http.NewRequest("GET", "/", nil)
	//res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler("template.html").
		ContentType(contentType).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Test()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
		}
	})
}

func BenchmarkArticleAndMenuParallel(b *testing.B) {
	b.StopTimer()
	req, _ := http.NewRequest("GET", "/", nil)
	//res := httptest.NewRecorder()

	chBus := make(chan *tools.Message, 1000)
	chMenu := make(chan *tools.Message, 1000)
	chArticle := make(chan *tools.Message, 1000)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler("template.html").
		ContentType(contentType).
		Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template("sitemap.html")).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template(pathRecord)).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Test()

	b.StartTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
		}
	})
}
