package webapi

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MarcGrol/chat/model"
	"github.com/stretchr/testify/assert"
)

func TestIntegrationRegister(t *testing.T) {
	peer := &model.Peer{Name: "me", Url: "TO_BE_RELACED"}

	// define server mock behavior
	server := mockServer(t,
		"POST", "/chat/peer", peer,
		200, model.Response{Status: true,
			Peers: []model.Peer{
				model.Peer{Name: "Eva", Url: "http://b"},
				model.Peer{Name: "Pien", Url: "http://c"},
			}})
	peer.Url = server.URL
	defer server.Close()

	// start sub-system under test
	webApi, outputQueue := startWebApi(server.URL)
	defer webApi.Stop()

	// trigger registration
	webApi.inputQueue <- model.Event{Type: model.EventTypeRegister, Peer: model.Peer{Name: "me", Url: server.URL}}

	// waitfor completion event and verify completion event
	event := <-outputQueue
	assert.Equal(t, model.EventTypeCompleted, event.Type)
	assert.NotEmpty(t, event.CompletedMsg)
	assert.Empty(t, event.ErrorMsg)
}

func TestIntegrationMsgOut(t *testing.T) {
	msg := &model.Msg{
		Sender:    model.Peer{Name: "me", Url: "localhost:54321"},
		Recipient: model.Peer{Name: "you", Url: "TO_BE_REPACED"},
		MsgText:   "msg from me to you"}

	// define server mock behavior
	server := mockServer(t,
		"POST", "/chat/msg", msg,
		200, model.Response{Status: true})
	msg.Recipient.Url = server.URL
	defer server.Close()

	// start sub-system under test
	webApi, outputQueue := startWebApi(server.URL)
	defer webApi.Stop()

	// trigger msg sent out
	webApi.inputQueue <- model.Event{Type: model.EventTypeSendMsg, Msg: *msg}

	// waitfor completion event and verify completion event
	event := <-outputQueue
	assert.Equal(t, model.EventTypeCompleted, event.Type)
	assert.NotEmpty(t, event.CompletedMsg)
	assert.Empty(t, event.ErrorMsg)
}

func TestIntegrationUnRegister(t *testing.T) {
}

func TestIntegrationPeers(t *testing.T) {
}

func TestIntegrationMsgIn(t *testing.T) {
}

func TestIntegrationQuit(t *testing.T) {
}

func mockServer(t *testing.T, expectedMethod string, expectedUrl string, expectedRequestBody interface{},
	httpCodeToReturn int, responseBodyToReturn interface{}) *httptest.Server {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedUrl, r.RequestURI)

		if expectedRequestBody != nil {
			// verify if client send the right request
			expectedBodyAsBytes, _ := encodeAnything(expectedRequestBody)
			actualBodyAsBytes, _ := ioutil.ReadAll(r.Body)
			assert.Equal(t, expectedBodyAsBytes.Bytes(), actualBodyAsBytes)
		}

		//log.Printf("Returning %d %+v", responseBodyToReturn, actualRequest)
		w.WriteHeader(httpCodeToReturn)
		w.Header().Set("Content-Type", "application/json")

		if responseBodyToReturn != nil {
			// return success response
			responseData, _ := encodeAnything(responseBodyToReturn)
			log.Printf("Returning response: %+v", responseData)
			fmt.Fprint(w, responseData)
		}
	}))

	return server
}

func startWebApi(url string) (*WebApi, chan model.Event) {
	stateMachineChannel := make(chan model.Event)
	webApi := NewWebApi(":54321", url)
	webApi.Start(stateMachineChannel)
	return webApi, stateMachineChannel
}
