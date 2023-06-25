package network

import (
	"fmt"
	"log"
	"time"

	"ctrader_events/messagebroker"
	"ctrader_events/messages/github.com/Carlosokumu/messages"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second
)

func ReadCtraderMessages(conn *websocket.Conn, messagehandler messagebroker.Hub) {
	fmt.Println("Reading Messages from Ctrader....🧔🏽‍♂️")
	defer func() {
		conn.Close()
	}()

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

		messagehandler.CtraderMessages <- messages.ProtoMessage{
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
			ClientMsgId: msg.ClientMsgId,
		}
		fmt.Println("Message..")
		fmt.Println(msg)

	}
}
