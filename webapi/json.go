package webapi

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/MarcGrol/chat/model"
)

func encodeAnything(anything interface{}) (*bytes.Buffer, error) {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	err := enc.Encode(anything)
	if err != nil {
		return nil, err
	}
	return &buffer, nil
}

func decodePeer(requestBody io.ReadCloser) (*model.Peer, error) {
	dec := json.NewDecoder(requestBody)
	var peer model.Peer
	err := dec.Decode(&peer)
	if err != nil {
		return nil, err
	}
	return &peer, nil
}

func decodeMsg(requestBody io.ReadCloser) (*model.Msg, error) {
	dec := json.NewDecoder(requestBody)
	var msg model.Msg
	err := dec.Decode(&msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func decodeResponse(responseBody io.ReadCloser) (*model.Response, error) {
	dec := json.NewDecoder(responseBody)
	var resp model.Response
	err := dec.Decode(&resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
