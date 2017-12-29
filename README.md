# MMOA

[![API documentation](https://godoc.org/github.com/claygod/mmoa?status.svg)](https://godoc.org/github.com/claygod/mmoa)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/mmoa)](https://goreportcard.com/report/github.com/claygod/mmoa)

MMOA Library is microframework written in Go, which helps to create a monolithic message-oriented applications. The architecture of these applications makes it easy to allocate a separate micro-service, any service from the MMOA.

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

# General principles

![alt tag](https://github.com/claygod/mmoa/blob/master/mmoa1.gif?raw=true)

MMOA concept as follows: developed sufficiently independent services (for this framework I have divided into several libraries), these services are combined into a single application and in the process of exchanging data with each other by sending messages. Address delivery of the message is from the service's name and the message subject. (*Themes* - is in fact the events in services, but in this case I prefer to call them precisely themes for messages).

Initially, having read clever books, I have created a very complex and intricate architecture, but in the course of development time, the husk has fallen off and been left a small amount of ingredients. We had to cut to the quick. From remote perhaps mention balancing channels that support priority queues. Well, let him rest in peace))

# Quick start

As part of the framework has a ready *example* that clearly demonstrates the MMOA the example of a simple website. Navigate to the folder example, compile and run the application. The result of playing in the browser at *localhost*. As a router, I used my design *Bxog*, but you can use any router of your choice. Standardization is carried out regular *html/template* library.

![alt tag](https://github.com/claygod/mmoa/blob/master/mmoa2.gif?raw=true)

### Details

In the sample application created two services:

- *article* - service articles. In order not to complicate the example of the work with the database, the article is stored in a plain text file that parses when the application loads. On request *Record* service gives the title and text of the article, nothing complicated. He also has supported the theme *List*, which gives a list of available items in the database.
- *menu* - this service is to give an array of id - name, ie list of articles. But since the articles have their own service, the *menu* requests list in article and getting it sends a response to the aggregator on their behalf. This decision (not the most optimal performance), is intended to show the interaction between a service. Originally I just wanted to put the hardcore in this service an array of key-value and give it to the request, but it would not be interesting.

# Composition MMOA

For convenience and ease of service creation as part of the application, some of the MMOA divided into separate libraries.

### tools

This library is needed to create both applications and services.

- *config.go* - file contains the types used in the application, setting timers and status
- *message.go* - basic and single unit of information exchange between services, the message content is stored in *MsgCtx*, everything else is an envelope.
- *themes.go* - structure with a list of all the services in the application, and they have taken the messages. The structure is designed specifically for convenience (*IDE* will not write the name of a non-existent service or not it supports the message subject).
- *exchanger.go* - stores all data structures in a format which services can communicate with each other in this file

### service

This library is used to create auxiliary services and is required when writing custom services.

- *logger.go* - a simple library to output logs to the console
- *service.go* - service basis, listening to the input port, if the messages have a suitable waiting, defines it there, if not, tries to invoke a method - reserved for the subject of the message, if not, sends a message to the trash.
- *waiting.go* - keeps the unit waiting for the arrival of messages, if aggregate of filled, it is sent to the appropriate method, if the unit is out of date, it is sent to the trash, and waiting is removed.
- *aggregate.go* - collects messages from specified *CID* (correlation identifier), after adding the number of posts gives more expected messages.

### support

This library contains the services aggregator and trash, as well as the bus for messaging.

- *aggregator.go* - is an analogue of *waiting*, with the difference that here the units accumulate messages to be sent to the handler, and in contrast to the waiting units here come from the handlers.
- *trash.go* - here are sent messages to the wrong address, incorrect, and the overdue units.
- *bus.go* - Bus receives messages and forwards them to the destination channels. If the recipient is not available, a message is sent to the trash.

### root MMOA

These files are the core MMOA, and outside are not used.

- *cid.go* - correlation identifier, the correlation *ID* to help identify the services posts.
- *handler.go* - request handler creates unit for the request and sends the request messages to the right services.
- *controller.go* - conduct an initial application provisioning, and creates a service bus services - aggregator and trash.
- *view.go* - is responsible for standardization. It contains templates for answers and services for the page.

# How to add the service to MMOA

For the set of words can be quite unclear how easy/difficult to use the concept described MMOA, so a simple example will describe the process of adding a new service. For example, we wanted to on the page are always displayed date.

### Update the list of services and themes

The *tools/services_themes.go* add structure
```go
// ThemeCalendar structure
type ThemeCalendar struct {
	Date TypeTHEME
}
```

*Themes* to add a line *Calendar ThemeCalendar*, and in the structure of *ListServices* add string *Calendar TypeSERVICE*

### Create a service file

Create a directory *example/calendar/* and it *calendar.go* file with the following content:
```go
package calendar

// Monolithic Message-Oriented Application (MMOA)
// Calendar
// Copyright  2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"time"

	"github.com/claygod/mmoa/service"
	"github.com/claygod/mmoa/tools"
)

func NewServiceCalendar(chIn chan *tools.Message, chBus chan *tools.Message) *ServiceCalendar {
	the := tools.NewThemes()
	s := &ServiceCalendar{
		service.NewService(the.Service.Calendar, chIn, chBus),
	}
	s.MethodWork = s.Work
	s.setEvents()
	s.Start()
	return s
}

type ServiceCalendar struct {
	service.Service
}

func (s *ServiceCalendar) setEvents() {
	s.Methods[s.The.Calendar.Date] = s.dateEvent
}

func (s *ServiceCalendar) dateEvent(msgIn *tools.Message) {
	t := time.Now()
	msgOut := tools.NewMessage().Cid(msgIn.MsgCid).
		From(s.Name).To(msgIn.AddsRe).
		Theme(msgIn.MsgTheme).
		Field("Day", t.Day()).
		Field("Month", int(t.Month())).
		Field("Year", t.Year())
	s.ChBus <- msgOut
}
```
### Templating

In the catalog *example/data*  create a file *date.html* with the string *{{.Day}}.{{Month}.}.{{Year}.}*, And the total content of the page template to change:
```html
<!DOCTYPE html>
<html>
<head>
<title>{{.Title}}</title>
<meta http-equiv="content-type" content="text/html; charset=utf-8"/>
<link rel="stylesheet" href="file/twocolumn.css">
</head>
<body>
	<div id="header"><h1>Rumba</h1></div>
	<div id="sidebar">
	{{.Sitemap}}
	
	{{.Date}}	
	</div>
	<div id="content">
	{{.Record}}	
	</div>
</body>
</html>
```

### Change application

Now we have to add a new library of files to import, create a channel for it, create a new structure and add a line in the initialization handler:
```go
package main

// Monolithic Message-Oriented Application (MMOA)
// Application
// Copyright  2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"github.com/claygod/Bxog"
	"github.com/claygod/mmoa"
	"github.com/claygod/mmoa/example/article"
	"github.com/claygod/mmoa/example/calendar"
	"github.com/claygod/mmoa/example/menu"
	"github.com/claygod/mmoa/tools"
)

const chLen int = 1000

func main() {
	chBus := make(chan *tools.Message, chLen)
	chMenu := make(chan *tools.Message, chLen)
	chArticle := make(chan *tools.Message, chLen)
	chCalendar := make(chan *tools.Message, chLen)

	the := tools.NewThemes()
	app := mmoa.NewController(chBus)

	sm := menu.NewServiceMenu(chMenu, chBus)
	app.AddService(sm.Name, chMenu)
	sa := article.NewServiceArticle(chArticle, chBus)
	app.AddService(sa.Name, chArticle)
	sc := calendar.NewServiceCalendar(chCalendar, chBus)
	app.AddService(sc.Name, chCalendar)

	hPage := app.Handler("./data/template.html").
	ContentType("text/html; charset=UTF-8").
	Service(tools.NewPart(sm.Name).Theme(the.Menu.Sitemap).Template("./data/sitemap.html")).
	Service(tools.NewPart(sa.Name).Theme(the.Article.Record).Template("./data/record.html")).
	Service(tools.NewPart(sc.Name).Theme(the.Calendar.Date).Template("./data/date.html")).
	StatusCodeOf(the.Article.Record)

	m := bxog.New()
	m.Add("/:id", hPage.Do)
	m.Start(":80")
}
```
# Performance

As we have made clear, despite the fact that the application is compiled, MMOA, it is not the best solution for the problems that have the main and decisive factor is speed, because in the process of service applications do not just send messages to each other through the channels, which naturally It is a brake. To at least roughly understand how productive MMOA purely reference held a simple ab tests running example of the *example* folder

- ab -n 10000 -c 1 --> 3127 r/s
- ab -n 30000 -c 100 --> 6373 r/s

Following good benchmark indicates that an application running only with the service *article* is much faster than with the *menu*, which sends the request to the *article*, waiting for him, receives a response, and only then send its response to the aggregator. (Note: In parallel mode the difference is somewhat reduced.)

- BenchmarkOnlyArticle-4 50000 24722 ns/op
- BenchmarkArticleAndMenu-4 30000 43404 ns/op
- BenchmarkOnlyArticleParallel-4 100000 13831 ns/op
- BenchmarkArticleAndMenuParallel-4 100000 20752 ns/op

# About

Most likely in many MMOA something familiar: the patterns mikroservice, SOA, MQ, etc. This is good and I can assure you, MMOA does not purport to overthrow or assigning themselves another's laurels. It is only a tool, the idea that you might be interested in. I would add only one thing IMHO - largely MMOA written under the influence of *Golang*, which I think is quite interesting and very suitable for the development of a variety of applications, and language authors thank for their work.
