package ui

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/MarcGrol/chat/model"
)

type Ui struct {
	username   string
	inputQueue chan model.Event
	reader     *bufio.Reader
	writer     *bufio.Writer
}

func NewUi(username string) *Ui {
	ui := new(Ui)
	ui.username = username
	ui.reader = bufio.NewReader(os.Stdin)
	ui.writer = bufio.NewWriter(os.Stdout)
	return ui
}

func (ui *Ui) Start(outputQueue chan model.Event) chan model.Event {

	ui.inputQueue = make(chan model.Event)

	go listenForEvents(ui.inputQueue, ui.writer)

	go ListenForUserInput(ui.reader, outputQueue)

	return ui.inputQueue
}

func (ui *Ui) Stop() {
	close(ui.inputQueue)
}

func ListenForUserInput(reader *bufio.Reader, outputQueue chan model.Event) {

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		event, err := parseCommandIntoEvent(string(line))
		if err != nil {
			event = &model.Event{Type: model.EventTypeError, ErrorMsg: err.Error()}
		}
		outputQueue <- *event
	}
}

func parseCommandIntoEvent(line string) (*model.Event, error) {
	parts := strings.Split(line, " ")
	if len(parts) > 0 {
		if parts[0] == "msg" && len(parts) >= 2 {
			msg := strings.Join(parts[1:], " ")
			return &model.Event{Type: model.EventTypeSendMsg, Msg: model.Msg{MsgText: msg}}, nil
		} else if parts[0] == "reg" && len(parts) >= 2 {
			return &model.Event{Type: model.EventTypeRegister, Peer: model.Peer{Name: parts[1]}}, nil
		} else if parts[0] == "unreg" && len(parts) >= 2 {
			return &model.Event{Type: model.EventTypeUnRegister, Peer: model.Peer{Name: parts[1]}}, nil
		} else {
			return nil, fmt.Errorf("unrecognized sub-command %s", line)
		}
	}
	return nil, fmt.Errorf("missing sub-command %s", line)
}

func listenForEvents(inputQueue chan model.Event, writer *bufio.Writer) {

	fmt.Fprintf(writer, "chat> ")
	for event := range inputQueue {
		log.Printf("Got incoming event: %+v", event)
		writeEvent(writer, event)
		fmt.Fprintf(writer, "chat> ")
	}
}

func writeEvent(writer *bufio.Writer, event model.Event) {
	switch event.Type {

	case model.EventTypeError:
		fmt.Fprintf(writer, "Error: %s\n", event.ErrorMsg)
		fmt.Fprintf(writer, "\tHelp: msg, reg, unreg\n")

	case model.EventTypeCompleted:
		fmt.Fprintf(writer, "%s\n", event.CompletedMsg)

	case model.EventTypeNewPeersReceived:
		fmt.Fprintf(writer, "Peers:\n")
		for _, peer := range event.Peers {
			fmt.Fprintf(writer, "Peer: %s @ %s\n", peer.Name, peer.Url)
		}

	case model.EventTypeMsgReceived:
		fmt.Fprintf(writer, "%s: '%s'\n", event.Msg.Sender.Name, event.Msg.MsgText)

	default:
		fmt.Fprintf(writer, "Error: %d %s\n", event.Type, event.ErrorMsg)
	}
}
