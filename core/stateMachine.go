package core

import (
	"log"
	"time"

	"github.com/MarcGrol/chat/model"
)

type StateMachine struct {
	username     string
	localAddress string
	Channel      chan model.Event
}

func NewStateMachine(username string, localAddress string) *StateMachine {
	disp := new(StateMachine)
	disp.username = username
	disp.localAddress = localAddress
	disp.Channel = make(chan model.Event)
	return disp

}

func (d *StateMachine) Start(uiChannel chan model.Event, webChannel chan model.Event) {
	log.Printf("Start state-machibe")

	// kick of registration
	regEvent := model.Event{Type: model.EventTypeRegister, Peer: model.Peer{Name: d.username, Url: d.localAddress}}
	log.Printf("Register: %+v", regEvent)
	webChannel <- regEvent

	tickerChannel := time.Tick(time.Millisecond * 5000)
	for {
		select {
		case event := <-d.Channel:
			log.Printf("Got: %+v", event)
			webChannel <- event

		case <-tickerChannel:
			log.Printf("Re-register: %+v", regEvent)
			webChannel <- regEvent

		}
	}

}
