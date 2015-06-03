package ui

import (
	"bufio"
	"strings"
	"testing"

	"github.com/MarcGrol/chat/model"
	"github.com/stretchr/testify/assert"
)

func TestParseRegister(t *testing.T) {
	outputQueue := make(chan model.Event)
	ui := NewUi("me")

	ui.reader = bufio.NewReader(strings.NewReader("reg me\n"))
	ui.Start(outputQueue)

	event := <-outputQueue
	assert.Equal(t, model.EventTypeRegister, event.Type)
	assert.Equal(t, "me", event.Peer.Name)
	assert.Equal(t, "", event.Peer.Url)
	assert.Empty(t, event.ErrorMsg)

	ui.Stop()
}

func TestParseSendMsg(t *testing.T) {
	outputQueue := make(chan model.Event)
	ui := NewUi("me")

	ui.reader = bufio.NewReader(strings.NewReader("msg hi there\n"))
	_ = ui.Start(outputQueue)

	event := <-outputQueue
	assert.Equal(t, model.EventTypeSendMsg, event.Type)
	assert.Equal(t, "hi there", event.Msg.MsgText)

	ui.Stop()
}

func TestParseUnregister(t *testing.T) {
	outputQueue := make(chan model.Event)
	ui := NewUi("me")

	ui.reader = bufio.NewReader(strings.NewReader("unreg me\n"))
	_ = ui.Start(outputQueue)

	event := <-outputQueue
	assert.Equal(t, model.EventTypeUnRegister, event.Type)
	assert.Equal(t, "me", event.Peer.Name)
	assert.Equal(t, "", event.Peer.Url)

	ui.Stop()
}

func TestParseGarbage(t *testing.T) {
	outputQueue := make(chan model.Event)
	ui := NewUi("me")

	ui.reader = bufio.NewReader(strings.NewReader("garbage\n"))
	_ = ui.Start(outputQueue)

	event := <-outputQueue
	assert.Equal(t, model.EventTypeError, event.Type)

	ui.Stop()
}

func TestQuit(t *testing.T) {
}

func TestMsgIn(t *testing.T) {
}

func TestPeers(t *testing.T) {
}

func TestRegistered(t *testing.T) {
}

func TestMsgSent(t *testing.T) {
}

func TestUnRegistered(t *testing.T) {
}
