package webapi

import (
	"testing"

	"github.com/MarcGrol/chat/model"
	"github.com/stretchr/testify/assert"
)

func TestRegisterSuccess(t *testing.T) {
	expectedResponse := model.Response{
		Status:   true,
		ErrorMsg: "",
		Peers: []model.Peer{
			model.Peer{Name: "Eva", Url: "http://b"},
			model.Peer{Name: "Pien", Url: "http://c"},
		}}

	// define server mock behavior
	server := mockServer(t,
		"POST", "/chat/peer", model.Peer{Name: "Marc", Url: "http://a"},
		200, expectedResponse)
	defer server.Close()

	// verify if client can handle the server response
	actualPeers, err := registerSelf(server.URL, model.Peer{Name: "Marc", Url: "http://a"})
	assert.NoError(t, err)
	assert.Len(t, actualPeers, len(expectedResponse.Peers))
	assert.Equal(t, expectedResponse.Peers[0].Name, actualPeers[0].Name)
	assert.Equal(t, expectedResponse.Peers[0].Url, actualPeers[0].Url)
	assert.Equal(t, expectedResponse.Peers[1].Name, actualPeers[1].Name)
	assert.Equal(t, expectedResponse.Peers[1].Url, actualPeers[1].Url)
}

func TestRegisterAuthError(t *testing.T) {

	// define server mock behavior
	server := mockServer(t,
		"POST", "/chat/peer", model.Peer{Name: "Marc", Url: "http://a"},
		403, model.Response{Status: false, ErrorMsg: "Invalid request"})
	defer server.Close()

	// verify if client can handle the server response
	peers, err := registerSelf(server.URL, model.Peer{Name: "Marc", Url: "http://a"})
	assert.Error(t, err)
	assert.Nil(t, peers)
	assert.Contains(t, err.Error(), " register")
}

func TestRegisterInternalError(t *testing.T) {
	server := mockServer(t,
		"POST", "/chat/peer", model.Peer{Name: "Marc", Url: "http://a"},
		403, "garbage")
	defer server.Close()

	// verify if client can handle the server response
	peers, err := registerSelf(server.URL, model.Peer{Name: "Marc", Url: "http://a"})
	assert.Error(t, err)
	assert.Nil(t, peers)
}
