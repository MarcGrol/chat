package webapi

import (
	"bytes"
	"testing"

	"github.com/MarcGrol/chat/model"
	"github.com/stretchr/testify/assert"
)

// Need to wrap data in "Closer" buffer to simulate reading from http response bodies
type ClosingBuffer struct {
	*bytes.Buffer
}

func (cb ClosingBuffer) Close() error {
	return nil
}

func TestFailEncodeOnMissingMandatoryField(t *testing.T) {
	peer := model.Peer{Name: "nyname"}
	_, err := encodeAnything(peer)
	assert.NoError(t, err) // TODO should fail
}

func TestFailDecodeOnMissingMandatoryField(t *testing.T) {
	json := `{"name":"my-name"}`
	data := bytes.NewBufferString(json)
	_, err := decodePeer(ClosingBuffer{data})
	assert.NoError(t, err) // TODO should fail
}

func TestEncodeDecodeResponse(t *testing.T) {
	response := model.Response{
		Status:   false,
		ErrorMsg: "Whatever error",
		Peers: []model.Peer{
			model.Peer{Name: "aname", Url: "http://anything"},
			model.Peer{Name: "anothername", Url: "http://otherthing"}}}

	// encode to json
	data, err := encodeAnything(response)
	assert.NoError(t, err)

	// decode from json
	responseAgain, err := decodeResponse(ClosingBuffer{data})

	// verify same as original
	assert.NotNil(t, responseAgain)
	assert.Equal(t, response.Status, responseAgain.Status)
	assert.Equal(t, response.ErrorMsg, responseAgain.ErrorMsg)
	for idx, expected := range responseAgain.Peers {
		assert.Equal(t, expected.Name, responseAgain.Peers[idx].Name)
		assert.Equal(t, expected.Url, responseAgain.Peers[idx].Url)
	}
}

func TestDecodeEncodePeer(t *testing.T) {
	json := `{"name":"aname","url":"http://anything"}
`
	data := bytes.NewBufferString(json)

	// decode from json
	peer, err := decodePeer(ClosingBuffer{data})
	assert.NoError(t, err)
	assert.NotNil(t, peer)

	// encode to json
	dataAgain, err := encodeAnything(peer)
	assert.NoError(t, err)

	// verify same as before
	assert.NotNil(t, dataAgain)
	assert.Equal(t, json, dataAgain.String())
}
