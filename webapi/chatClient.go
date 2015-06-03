package webapi

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/MarcGrol/chat/model"
)

func registerSelf(registerHostnamePort string, peer model.Peer) ([]model.Peer, error) {
	registerUrl := fmt.Sprintf("%s/chat/peer", registerHostnamePort)

	log.Printf("Perform registration to registry to %s", registerHostnamePort)
	response, err := doHttp("POST", registerUrl, peer)
	if err != nil {
		return nil, fmt.Errorf("Error posting register-request to %v: %s",
			registerUrl, err.Error())
	} else if response.Status == false {
		return nil, fmt.Errorf("Error posting register-request to %v: %s",
			registerUrl, response.ErrorMsg)
	}

	return response.Peers, nil
}

func sendMsg(recipientHostnamePort string, sender model.Peer, recipient model.Peer, msgText string) error {
	recipientUrl := fmt.Sprintf("%s/chat/msg", recipientHostnamePort)

	log.Printf("Send message to peer on %s", recipientUrl )
	response, err := doHttp("POST", recipientUrl, model.Msg{Sender: sender, Recipient: recipient, MsgText: msgText})
	if err != nil {
		return fmt.Errorf("Error posting chat-request to %v: %s",
			recipientUrl, err.Error())
	} else if response.Status == false {
		return fmt.Errorf("Error posting chat-msg to %v: %s",
			recipientUrl, response.ErrorMsg)
	}

	return nil
}

func doHttp(method string, url string, command interface{}) (*model.Response, error) {
	var requestBody *bytes.Buffer = nil
	var err error
	if method == "POST" || method == "PUT" {
		// serialize request to json
		requestBody, err = encodeAnything(command)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling request: %s", err.Error())
		}
	}

	log.Printf("Perform HTTP %s on url '%s' : '%v'", method, url, requestBody)

	// perform http call
	req, err := http.NewRequest(method, url, requestBody)
	httpResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		//log.Printf("Error performing http POST to url %s: %s", url, err.Error())
		return nil, fmt.Errorf("Error performing http POST to url %s: %s", url, err.Error())
	}

	// decode response
	applicationResponse, err := decodeResponse(httpResponse.Body)
	if err != nil {
		log.Printf("Error unmarshalling response: %s", err.Error())
		return nil, fmt.Errorf("Error unmarshalling response: %s", err.Error())
	}

	log.Printf("Received '%v' : '%+v'", httpResponse.StatusCode, applicationResponse)

	// evaluate application status
	if httpResponse.StatusCode != http.StatusOK {
		return applicationResponse, fmt.Errorf(http.StatusText(httpResponse.StatusCode))
	}

	return applicationResponse, nil
}
