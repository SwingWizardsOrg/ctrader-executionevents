package client

import (
	"ctrader_events/messagebroker"
	"ctrader_events/network"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

func ConnectToCtrader(host string, port int) {

	// Set up a dialer with the desired options
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true

	// Connect to the  Ctrader WebSocket endpoint
	url := fmt.Sprintf("wss://%s:%d", host, port)
	fmt.Println(url)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	network.AuthorizeApp(conn)

	hub := messagebroker.NewHub(conn)
	go hub.Run()
	go network.ReadCtraderMessages(conn, *hub)
}
