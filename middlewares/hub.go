package middlewares

import (
	"ctrader_events/messages/github.com/swingwizards/messages"
	"ctrader_events/models"

	"fmt"
)

//Hub Contains different channels that messages can be sent to:
//As example :- The AppAuthResChannnel channel will receive messages related to App Authetication

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
		Protos:                         make(chan messages.ProtoMessage),
		resourceid:                     make(chan models.ResourceId),
		protosback:                     make(chan messages.ProtoMessage),
		AssetListChannnel:              make(chan messages.ProtoMessage),
		AccountAuthResChannnel:         make(chan messages.ProtoMessage),
		AppAuthResChannnel:             make(chan messages.ProtoMessage),
		TraderResChannnel:              make(chan messages.ProtoMessage),
		MarketOrderListChannnel:        make(chan messages.ProtoMessage),
		LightSymbolChannel:             make(chan messages.ProtoMessage),
		Symbols:                        make(chan messages.ProtoMessage),
		SymbolModelChannel:             make(chan []models.SymbolModel),
		AccountOrdersChannel:           make(chan messages.ProtoOAReconcileRes),
		SeparateSpotChannel:            make(chan messages.ProtoMessage),
		SpotEventChannel:               make(chan messages.ProtoMessage),
		ConversionLightSymbols:         make(chan messages.ProtoMessage),
		LightSymbolsChannel:            make(chan []messages.ProtoOALightSymbol),
		AccounConversionSymbolsChannel: make(chan []models.SymbolModel),
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

//This will "channel" messages to the appropriate channels declared by the Hub

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
			go func() {
				for {
					h.SpotEventChannel <- protomessage

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
