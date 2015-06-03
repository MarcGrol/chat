package webapi

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcGrol/chat/model"
	"github.com/gorilla/mux"
)

func startServer(listenUrl string, outputQueue chan model.Event) {
	log.Printf("Start listening on url %s", listenUrl)

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/chat/peer", handlePeerRegistration(outputQueue))
	router.HandleFunc("/chat/msg", handleChatMsg(outputQueue))

	http.ListenAndServe(listenUrl, router)
}

func handlePeerRegistration(outputQueue chan model.Event) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// read request from peer
		var peers []model.Peer
		err := json.NewDecoder(r.Body).Decode(peers)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// report event to core
		outputQueue <- model.Event{Type: model.EventTypeNewPeersReceived, Peers: peers}

		// return response to caller
		resp := model.Response{Status: true}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func handleChatMsg(outputQueue chan model.Event) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		// read request from peer
		var msg model.Msg
		err := json.NewDecoder(r.Body).Decode(msg)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// report event to core
		outputQueue <- model.Event{Type: model.EventTypeMsgReceived, Msg: msg}

		// return response to caller
		resp := model.Response{Status: true}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
