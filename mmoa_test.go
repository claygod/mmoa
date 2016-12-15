package mmoa

// Monolithic Message-Oriented Application (MMOA)
// Test
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

import (
	"testing"
)

// Application tests combined with an example and there are: example/app_test.go

const (
	pathTemplate string = "./data/template.html"
	pathSitemap  string = "./data/sitemap.html"
	pathRecord   string = "./data/record.html"
	contentType  string = "./text/html; charset=UTF-8"
)

func TestCid(t *testing.T) {
	cid := NewCid()
	if cid.Get() != 0 {
		t.Error("Incorrect number from CID")
	}
	if cid.Get() != 1 {
		t.Error("Incorrect number from CID")
	}
}

func TestView(t *testing.T) {
	view := NewView()
	view.ContentType("utf-8")
	if view.contentType != "utf-8" {
		t.Error("Error in method ContentType")
	}
	view.StatusCodeOf("article")
	if view.statusCodeOf != "article" {
		t.Error("Error in method StatusCodeOf")
	}
}
