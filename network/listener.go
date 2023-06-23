package network

import (
	"fmt"
	"log"
	"time"

	"ctrader_events/messages/github.com/Carlosokumu/messages"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
)

func ReadCtraderMessages(conn *websocket.Conn) {
	fmt.Println("Reading Messages from Ctrader....🧔🏽‍♂️")

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msg := &messages.ProtoMessage{}
		_, readmessage, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		unmarsherr := proto.Unmarshal(readmessage, msg)

		if unmarsherr != nil {
			fmt.Println(unmarsherr)
		}

		fmt.Println("Message..")
		fmt.Println(msg)

	}
}
