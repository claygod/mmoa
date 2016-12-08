# MMOA

[![API documentation](https://godoc.org/github.com/claygod/mmoa?status.svg)](https://godoc.org/github.com/claygod/mmoa)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/mmoa)](https://goreportcard.com/report/github.com/claygod/mmoa)

MMOA Library is microframework written in Go, which helps to create a monolithic message-oriented applications. The architecture of these applications makes it easy to allocate a separate micro-service, any service from the MMOA.

# General principles

![alt tag](https://github.com/claygod/mmoa/blob/master/mmoa1.gif?raw=true)

MMOA concept as follows: developed sufficiently independent services (for this framework I have divided into several libraries), these services are combined into a single application and in the process of exchanging data with each other by sending messages. Address delivery of the message is from the service's name and the message subject. (*Themes* - is in fact the events in services, but in this case I prefer to call them precisely themes for messages).

Initially, having read clever books, I have created a very complex and intricate architecture, but in the course of development time, the husk has fallen off and been left a small amount of ingredients. We had to cut to the quick. From remote perhaps mention balancing channels that support priority queues. Well, let him rest in peace))

# Quick start

As part of the framework has a ready example that clearly demonstrates the MMOA the example of a simple website. Navigate to the folder example, compile and run the application. The result of playing in the browser at *localhost*. As a router, I used my design *Bxog*, but you can use any router of your choice. Standardization is carried out regular *html/template* library.

![alt tag](https://github.com/claygod/mmoa/blob/master/mmoa2.gif?raw=true)

### Details

In the sample application created two services:

- *article* - service articles. In order not to complicate the example of the work with the database, the article is stored in a plain text file that parses when the application loads. On request *Record* service gives the title and text of the article, nothing complicated. He also has supported the theme *List*, which gives a list of available items in the database.
- *menu* - this service is to give an array of id - name, ie list of articles. But since the articles have their own service, the *menu* requests list in article and getting it sends a response to the aggregator on their behalf. This decision (not the most optimal performance), is intended to show the interaction between a service. Originally I just wanted to put the hardcore in this service an array of key-value and give it to the request, but it would not be interesting.

# Composition MMOA

For convenience and ease of service creation as part of the application, some of the MMOA divided into separate libraries.

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
