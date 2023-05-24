package middlewares

import (
	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// ctrader websocket connection.
	Conn *websocket.Conn

	//will hold ctrader protomessages
	Protomessages chan messages.ProtoMessage

	// //app's websocket connection
	Appconn *websocket.Conn

	// //will hold json messages
	// resources chan models.ResourceId
}
