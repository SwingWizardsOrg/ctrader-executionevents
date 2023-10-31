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
func AuthorizeAccount(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ)
	accountId := credentials.AccountId
	accessToken := credentials.AccessToken
	messageId := "A/C_AUTH_REQ"
	acReq := &messages.ProtoOAAccountAuthReq{
		CtidTraderAccountId: &accountId,
		AccessToken:         &accessToken,
	}
	acBytes, err := proto.Marshal(acReq)
	if err != nil {
		log.Fatal(err)
	}

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     acBytes,
		ClientMsgId: &messageId,
	}
	protoMessage, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.WriteMessage(MessageType, protoMessage)
	if err != nil {
		log.Fatal(err)
	}

}

func SendHeartBeatMessage(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ)

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
	}
	protoMessage, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.WriteMessage(MessageType, protoMessage)
	if err != nil {
		log.Fatal(err)
	}

}
