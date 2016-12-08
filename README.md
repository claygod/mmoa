# MMOA


# Usage

An example of using the MMOA:

```go

package main

// Monolithic Message-Oriented Application (MMOA)
// Application
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"github.com/claygod/Bxog"
	"github.com/claygod/mmoa"
	"github.com/claygod/mmoa/example/article"
	"github.com/claygod/mmoa/example/menu"
	"github.com/claygod/mmoa/tools"
)

const chLen int = 1000

func main() {
	chBus := make(chan *tools.Message, chLen)
	chMenu := make(chan *tools.Message, chLen)
	chArticle := make(chan *tools.Message, chLen)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)

	hPage := app.Handler("./data/template.html").
		ContentType("text/html; charset=UTF-8").
		Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template("./data/sitemap.html")).
		Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template("./data/record.html")).
		StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Start(":80")
}

```
