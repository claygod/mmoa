package tools

// Monolithic Message-Oriented Application (MMOA)
// Config
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "errors"
import "time"

type (
	// TypeCID - correlation identifier
	TypeCID uint64
	// TypeTHEME - the theme messages in service
	TypeTHEME string
	// TypeSERVICE - the name of the service
	TypeSERVICE string
	// TypeTIME - Unix time
	TypeTIME int64
)

// Conf
const (
	CleanerTimerSec     time.Duration = 1e9
	DurationHandle      TypeTIME      = 5e9
	EmptyServiceAddress TypeSERVICE   = ""
)

// Status codes
const (
	StatusOK       int = 200
	StatusNotFound int = 404
	StatusTimeout  int = 504
)
