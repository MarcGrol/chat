package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MarcGrol/chat/model"
)

var (
	listenPort *int
)

type peerHandler struct {
	peers []model.Peer
}

func main() {
	readCommandLineFlags(os.Args[0])
	startServer(*listenPort)
}

func readCommandLineFlags(progname string) error {
	flags := flag.NewFlagSet(progname, flag.ExitOnError)
	listenPort = flags.Int("listenPort", 8081, "Hostname of the central chat registry")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing command-line: %v", err)
		return err
	}
	return nil
}

func (ph *peerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// read request from peer
	var peer model.Peer
	err := json.NewDecoder(r.Body).Decode(peer)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// append to list
	ph.peers = append(ph.peers, peer)

	// return response to caller
	resp := model.Response{Status: true}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func startServer(port int) {
	log.Printf("Start listening on port %d", port)

	mux := http.NewServeMux()

	peerHandler := peerHandler{peers: make([]model.Peer, 0, 10)}
	mux.Handle("/char/peer", &peerHandler)

	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
