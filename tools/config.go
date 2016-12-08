package tools

// Monolithic Message-Oriented Application (MMOA)
// Config
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "errors"
import "time"

// Types
type (
	TypeCID     uint64
	TypeTHEME   string
	TypeSERVICE string
	TypeTIME    int64
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
