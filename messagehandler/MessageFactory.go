package messagehandler

import (
	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"google.golang.org/protobuf/proto"
)

func GetMessage(protoMessage messages.ProtoMessage) Imessage {
	switch *protoMessage.PayloadType {
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES):
		traderRes := messages.ProtoOATraderRes{}
		proto.Unmarshal(protoMessage.Payload, &traderRes)
		return &Trader{
			name: "Trader",
		}
	default:
		return nil

	}
}
