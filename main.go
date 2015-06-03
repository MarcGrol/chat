package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/MarcGrol/chat/core"
	"github.com/MarcGrol/chat/ui"
	"github.com/MarcGrol/chat/webapi"
)

var (
	username             *string
	registryHostnamePort *string
	localListenPort      *int
)

func main() {
	readCommandLineFlags(os.Args[0])

	// Create state-machine that dispaches events back and  forth
	stateMachine := core.NewStateMachine(*username, createRegisterUrl(*localListenPort))

	// Start the http interface to communicate with others
	webApi := webapi.NewWebApi(createListenUrl(*localListenPort), *registryHostnamePort)
	webChannel := webApi.Start(stateMachine.Channel)

	// Start the UI
	ui := ui.NewUi(*username)
	uiChannel := ui.Start(stateMachine.Channel)

	// Start the state machine
	stateMachine.Start(uiChannel, webChannel)

}

func readCommandLineFlags(progname string) error {
	flags := flag.NewFlagSet(progname, flag.ExitOnError)
	username = flags.String("username", "marc", "name chat client")
	registryHostnamePort = flags.String("registry", "http://localhost:8081", "Hostname of the central chat registry")
	localListenPort = flags.Int("localPort", 8081, "Port to listen on for incoming messages")
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("Error parsing command-line: %v", err)
		return err
	}
	return nil
}

func createListenUrl(port int) string {
	return fmt.Sprintf(":%d", port)
}

func createRegisterUrl(port int) string {
	return fmt.Sprintf("http://%s:%d", getLocalIpv4(), port)
}

func getLocalIpv4() string {
	host, _ := os.Hostname()
	addrs, _ := net.LookupIP(host)
	for _, addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			return fmt.Sprintf("%s", ipv4)
		}
	}
	return "localhost"
}
