package messagebroker

import (
	"ctrader_events/credentials"
	"ctrader_events/messages/github.com/Carlosokumu/messages"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"google.golang.org/protobuf/proto"
)

const (
	MessageType = 2
)

type Hub struct {
	CtraderMessages chan messages.ProtoMessage
	Conn            *websocket.Conn
}

func NewHub(conn *websocket.Conn) *Hub {
	return &Hub{
		CtraderMessages: make(chan messages.ProtoMessage),
		Conn:            conn,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case protoMessage := <-h.CtraderMessages:
			handleMessage(protoMessage, h)
		}
	}
}

func handleMessage(protomessage messages.ProtoMessage, h *Hub) {
	switch *protomessage.PayloadType {
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_RES):
		{
			fmt.Println("Application has been authorized")
			AuthorizeAccount(h.Conn)
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_RES):
		{
			fmt.Println("Account has been authorized")
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES):
		{

		}
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
