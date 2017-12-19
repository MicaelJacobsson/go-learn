package main

import (
	"flag"
	"fmt"
	zmq "github.com/pebbe/zmq4"
)

type client struct {
	id          int
	query       string
	last_result string
}

var getClients = client{0, "", ""}

func main() {
	var subPortPtr = flag.Int("server_port", 9303, "Port to publish subscription server")
	var pubPortPtr = flag.Int("publisher_port", 9304, "Port to publish data publisher")
	var consulPortPtr = flag.Int("consul_port", 8500, "Consul REST API port")
	var oamIpPtr = flag.String("oam_ip", "", "If excluded will query consul at localhost for om service IP")

	flag.Parse()

	fmt.Printf("Subscription server port: %d, publisher port: %d, consul port: %d, oamIp: %v\n", *subPortPtr, *pubPortPtr, *consulPortPtr, *oamIpPtr)

	//TODO: Query consul for oam ip, fail if not found

	toPub := make(chan client)
	clientsCh := make(chan client)

	go clientHandler(clientsCh)
	go subscriptionServer(clientsCh)
	go publisher(toPub)
	go oamPoller(toPub, clientsCh)

}

// publisher publish a PUB ZeroMQ socket and broadcast client data
// received from oamPoller
func publisher(data <-chan client) {

}

// clientHandler will keep a list of client to poll for and notify
// with new data.
func clientHandler(clientCh chan client) {

}

// oamPoller will poll oam with all client queries.
// TODO: Should react to events from event bus instead?
func oamPoller(toPub chan<- client, getClients chan client) {

}

// subscriptionServer will read subscription requests,
// and add to list of clients
func subscriptionServer(addClient chan<- client) {
	context, _ := zmq.NewContext()
	//defer context.Close()

	server, err := context.NewSocket(zmq.REP)
	if err != nil {
		fmt.Println("%v", err)
		return
	}
	defer server.Close()
	server.Bind("tcp://127.0.0.1:5555")

	for {
		_, err := server.RecvMessage(0)
		if err != nil {
			fmt.Println("%v", err)
		}
		server.SendMessage("")
	}
}
