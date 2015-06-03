package webapi

import (
	"fmt"
	"log"

	"github.com/MarcGrol/chat/model"
)

type WebApi struct {
	listenUrl            string
	registryHostnamePort string
	inputQueue           chan model.Event
}

func NewWebApi(listenUrl string, registryHostnamePort string) *WebApi {
	api := new(WebApi)
	api.listenUrl = listenUrl
	api.registryHostnamePort = registryHostnamePort
	return api
}

func (api *WebApi) Start(outputQueue chan model.Event) chan model.Event {
	log.Printf("Start web-api")

	api.inputQueue = make(chan model.Event)

	go startServer(api.listenUrl, outputQueue)

	go startClient(api.inputQueue, api.registryHostnamePort, outputQueue)

	return api.inputQueue
}

func startClient(inputQeue chan model.Event, registryHostnamePort string, outputQueue chan model.Event) {
	log.Printf("Start listening for event to forward over the web")
	for event := range inputQeue {
		switch event.Type {
		case model.EventTypeRegister:
			_, err := registerSelf(registryHostnamePort, event.Peer)
			if err != nil {
				outputQueue <- model.Event{Type: model.EventTypeError, ErrorMsg: fmt.Sprintf("Registration error: %s", err.Error())}
			} else {
				outputQueue <- model.Event{Type: model.EventTypeCompleted, CompletedMsg: "Registration success"}
			}
		case model.EventTypeSendMsg:
			log.Printf("got %+v", event)
			err := sendMsg(event.Msg.Recipient.Url, event.Msg.Sender, event.Msg.Recipient, event.Msg.MsgText)
			if err != nil {
				outputQueue <- model.Event{Type: model.EventTypeError, ErrorMsg: fmt.Sprintf("Send msg error: %s", err.Error())}
			} else {
				outputQueue <- model.Event{Type: model.EventTypeCompleted, CompletedMsg: "Send msg success"}
			}
		}
	}
}
func (api *WebApi) Stop() {
	close(api.inputQueue)
}
