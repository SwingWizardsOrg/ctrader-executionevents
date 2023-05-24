package service

import (
	"fmt"

	"ctraderapi/middlewares"

	"github.com/gorilla/websocket"
)

// messagechan := make(chan messages.ProtoMessage)
type MessageHandler interface {
	Execute(conn *websocket.Conn, h *middlewares.Hub)
	SetNext(MessageHandler)
}

type AppAuth struct {
	next         MessageHandler
	isAuthorized bool
}

type AccountAuth struct {
	next MessageHandler
}

type TraderInfo struct {
	next MessageHandler
}

type AssetList struct {
	next MessageHandler
}

type MarketOrder struct {
	next MessageHandler
}
type LightSymbol struct {
	next MessageHandler
}

type Symbol struct {
	next MessageHandler
}

type SpotSubscriber struct {
	next MessageHandler
}

func (appauth *AppAuth) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Authorizing the application...")
	AuthorizeApp(conn, h)
}

func (appauth *AppAuth) SetNext(next MessageHandler) {
	appauth.next = next
}

func (accountauth *AccountAuth) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Authorizing trading  account")
	AuthorizeAccount(conn, h)

}

func (accountauth *AccountAuth) SetNext(next MessageHandler) {
	accountauth.next = next
}

func (traderinfo *TraderInfo) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Fetching Trader Information...")
	GetTrader(conn, h)

}

func (traderinfo *TraderInfo) SetNext(next MessageHandler) {
	traderinfo.next = next
}

func (assetlist *AssetList) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Fetching trader's asset list...")
	GetAssets(conn, h)

}

func (assetlist *AssetList) SetNext(next MessageHandler) {
	assetlist.next = next
}

func (marketOrder *MarketOrder) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Fetching trader's open positions...")
	GetAccountOrders(conn, h)

}

func (marketOrder *MarketOrder) SetNext(next MessageHandler) {
	marketOrder.next = next
}

func (lightSymbol *LightSymbol) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Fetching  lightsymbols..")
	GetLightSymbolList(conn, h)

}

func (lightSymbol *LightSymbol) SetNext(next MessageHandler) {
	lightSymbol.next = next
}

func (symbol *Symbol) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Fetching  full symbol entity..")
	GetSymbols(conn, h)

}

func (symbol *Symbol) SetNext(next MessageHandler) {
	symbol.next = next
}

func (subscriber *SpotSubscriber) Execute(conn *websocket.Conn, h *middlewares.Hub) {
	fmt.Println("Subscribing to spots..")
	SendSubscribeSpotsRequest(conn)

}

func (subscriber *SpotSubscriber) SetNext(next MessageHandler) {
	subscriber.next = next
}
