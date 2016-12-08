package support

// Monolithic Message-Oriented Application (MMOA)
// Aggregator
// Copyright © 2016 Eduard Sesigin. All rights reserved. Contacts: <claygod@yandex.ru>

//import "reflect"
//import "sync"
//import "time"
import "errors"
import "github.com/claygod/mmoa/tools"
import "github.com/claygod/mmoa/service"

//import "fmt"

// NewAggregator - create a new Aggregator
func NewAggregator(chIn chan *tools.Message, chBus chan *tools.Message) *Aggregator {
	The := tools.NewThemes()
	s := &Aggregator{
		The,
		service.NewService(The.Service.Aggregator, chIn, chBus),
		service.NewLogger(),
	}
	s.MethodWork = s.Work
	s.Start()
	return s
}

// Aggregator structure
type Aggregator struct {
	The *tools.Themes
	service.Service
	logger *service.Logger
}

// Aggregate - add a new aggregate in The aggregator
func (a *Aggregator) Aggregate(cid tools.TypeCID, duration tools.TypeTIME, messages map[string]*tools.Message, ch chan *tools.Message) error {
	//currentTime := tools.TypeTIME(time.Now().UnixNano())
	//fmt.Print("\n", a.The.Service.Controller, a.The.Service)
	msg := tools.NewMessage().Cid(cid).
		From(a.Name).To(a.The.Service.Controller).
		Theme(a.The.Aggregator.Filled)
	//ae := service.NewAggregate(duration, messages, ch)
	/*
		ae := &Aggregate{
			0,
			//cid,
			currentTime,
			currentTime + duration,
			messages,
			len(messages),
			ch,
		}
	*/
	//a.WaitingFor.Lock()
	//if _, ok := a.WaitingFor.Arr[cid]; ok {
	//	a.WaitingFor.Unlock()
	//	return errors.New("Aggregate with the cID already exists")
	//}
	//wg := &service.Waiting{
	//	ae,
	//	msg,
	//}

	if err := a.WaitingFor.NewWaiting(msg, duration, messages, ch); err != nil {
		return errors.New("Aggregate with the cID already exists")
	}

	//a.WaitingFor.Arr[msg.MsgCid] = wg
	//a.WaitingFor.Unlock()
	return nil
}
