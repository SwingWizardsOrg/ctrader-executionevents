package messagebroker

import (
	"ctrader_events/credentials"
	"ctrader_events/messages/github.com/Carlosokumu/messages"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"google.golang.org/protobuf/proto"
)

const (
	MessageType = 2
	// Time allowed to read the next pong message from the peer.
	pongWait = 35 * time.Second
)

type Hub struct {
	CtraderMessages chan messages.ProtoMessage
	Conn            *websocket.Conn
	TimerChannel    chan time.Time
}

func NewHub(conn *websocket.Conn) *Hub {
	return &Hub{
		CtraderMessages: make(chan messages.ProtoMessage),
		Conn:            conn,
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(pongWait)
	for {
		select {

		case <-ticker.C:
			//Send a heartbeat message.
		case protoMessage := <-h.CtraderMessages:
			handleMessage(protoMessage, h)

		}
	}
}

func handleMessage(protomessage messages.ProtoMessage, h *Hub) {
	switch *protomessage.PayloadType {
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_RES):
		{
			fmt.Println("Application has been authorized 📿")
			AuthorizeAccount(h.Conn)
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_RES):
		{
			fmt.Println("Service is Live 🚀")
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES):
		{

		}
	case uint32(messages.ProtoPayloadType_HEARTBEAT_EVENT):
		{
			select {
			case <-time.After(10 * time.Second):
				// Send Back a heartBeat Message to the Server to keep it Alive.
				SendHeartBeatMessage(h.Conn)

			}

		}
	case uint32(messages.ProtoPayloadType_ERROR_RES):
		{

			panic("Unimplemented")

		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ERROR_RES):
		{

			panic("Unimplemented")

		}
	default:
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

func SendHeartBeatMessage(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoPayloadType_HEARTBEAT_EVENT)

	payloadType := messages.ProtoPayloadType_HEARTBEAT_EVENT

	heartbeatEvent := &messages.ProtoHeartbeatEvent{
		PayloadType: &payloadType,
	}

	heartbeatEventBytes, err := proto.Marshal(heartbeatEvent)

	if err != nil {
		log.Fatal(err)
	}

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     heartbeatEventBytes,
	}

	protoMessage, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.WriteMessage(2, protoMessage)
	if err != nil {
		log.Fatal(err)
	}

}
