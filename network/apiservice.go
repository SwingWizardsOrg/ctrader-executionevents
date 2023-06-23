package network

import (
	"ctrader_events/credentials"
	"ctrader_events/messages/github.com/Carlosokumu/messages"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

var (
	MessageType = 2
)

func AuthorizeApp(conn *websocket.Conn) {
	clientId := credentials.ClientId
	clientSecret := credentials.ClientSecret
	authReq := &messages.ProtoOAApplicationAuthReq{
		ClientId:     &clientId,
		ClientSecret: &clientSecret,
	}
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_REQ)
	messageId := "APP_AUTH_REQ"

	authReqBytes, err := proto.Marshal(authReq)
	if err != nil {
		log.Fatal(err)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     authReqBytes,
		ClientMsgId: &messageId,
	}
	protoMessage, _ := proto.Marshal(message)

	// Serialize the message to a byte slice
	writeerror := conn.WriteMessage(MessageType, protoMessage)

	if writeerror != nil {
		log.Fatal(writeerror)
	}
}
