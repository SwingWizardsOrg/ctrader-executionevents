package middlewares

import (
	"ctrader_events/messages/github.com/swingwizards/messages"
	"ctrader_events/models"

	"fmt"
)

type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	//  Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	//holds messages from the ctrader
	Protos chan messages.ProtoMessage

	//holds messages from the ctrader
	protosback chan messages.ProtoMessage

	//holds messages from  our frontend apps
	resourceid chan models.ResourceId

	AccountAuthResChannnel chan messages.ProtoMessage

	AppAuthResChannnel             chan messages.ProtoMessage
	TraderResChannnel              chan messages.ProtoMessage
	AssetListChannnel              chan messages.ProtoMessage
	MarketOrderListChannnel        chan messages.ProtoMessage
	LightSymbolChannel             chan messages.ProtoMessage
	Symbols                        chan messages.ProtoMessage
	SymbolModelChannel             chan []models.SymbolModel
	AccountOrdersChannel           chan messages.ProtoOAReconcileRes
	SpotEventChannel               chan messages.ProtoMessage
	ConversionLightSymbols         chan messages.ProtoMessage
	LightSymbolsChannel            chan []messages.ProtoOALightSymbol
	AccounConversionSymbolsChannel chan []models.SymbolModel
	AccountModelChannel            chan models.AccountModel
	SeparateSpotChannel            chan messages.ProtoMessage
	SubChannel                     chan models.AccountModel
}

func NewHub() *Hub {
	return &Hub{
		Register:                       make(chan *Client),
		Unregister:                     make(chan *Client),
		clients:                        make(map[*Client]bool),
		Protos:                         make(chan messages.ProtoMessage, 10000000),
		resourceid:                     make(chan models.ResourceId),
		protosback:                     make(chan messages.ProtoMessage, 1000000),
		AssetListChannnel:              make(chan messages.ProtoMessage, 100000),
		AccountAuthResChannnel:         make(chan messages.ProtoMessage, 100000),
		AppAuthResChannnel:             make(chan messages.ProtoMessage, 1000000),
		TraderResChannnel:              make(chan messages.ProtoMessage, 1000000),
		MarketOrderListChannnel:        make(chan messages.ProtoMessage, 100000),
		LightSymbolChannel:             make(chan messages.ProtoMessage, 1000000),
		Symbols:                        make(chan messages.ProtoMessage, 100000),
		SymbolModelChannel:             make(chan []models.SymbolModel, 1000000),
		AccountOrdersChannel:           make(chan messages.ProtoOAReconcileRes, 1000000),
		SeparateSpotChannel:            make(chan messages.ProtoMessage, 1000000),
		SpotEventChannel:               make(chan messages.ProtoMessage, 10000000),
		ConversionLightSymbols:         make(chan messages.ProtoMessage, 100000),
		LightSymbolsChannel:            make(chan []messages.ProtoOALightSymbol, 100000000),
		AccounConversionSymbolsChannel: make(chan []models.SymbolModel, 100000),
		AccountModelChannel:            make(chan models.AccountModel),
		SubChannel:                     make(chan models.AccountModel),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)

			}
		case protoMessage := <-h.Protos:
			ChannelMessage(protoMessage, h)
		}
	}
}

func ChannelMessage(protomessage messages.ProtoMessage, h *Hub) {
	switch *protomessage.PayloadType {
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_RES):
		{
			h.AppAuthResChannnel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_RES):
		{
			h.AccountAuthResChannnel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES):
		{
			h.TraderResChannnel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ASSET_LIST_RES):
		{
			h.AssetListChannnel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_RECONCILE_RES):
		{
			h.MarketOrderListChannnel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_RES):
		{
			h.LightSymbolChannel <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_RES):
		{
			h.Symbols <- protomessage
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_SPOT_EVENT):
		{

			fmt.Println("Spot...")
			//h.SpotEventChannel <- protomessage
			go func() {
				for {
					h.SpotEventChannel <- protomessage
					//fmt.Println("Passed spot", <-h.SpotEventChannel)
				}
			}()

		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_RES):
		{
			h.ConversionLightSymbols <- protomessage
		}

	default:
	}
}
