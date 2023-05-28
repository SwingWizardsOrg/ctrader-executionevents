package models

import (
	"ctraderapi/mappers"
	"ctraderapi/messages/github.com/Carlosokumu/messages"
	"fmt"
)

type Trader struct {
	Balance     int64 `json:"balance"`
	TraderLogin int64 `json:"traderlogin"`
}

type PositionsInfo struct {
	Positions []*messages.ProtoOAPosition `json:"positions"`
}
type SymbolInformation struct {
	Symbols []*messages.ProtoOASymbol `json:"symbolinfo"`
}

type AssetInfo struct {
	Assets []*messages.ProtoOAAsset `json:"assets"`
}

type ConversionInfo struct {
	SymbolChain []*messages.ProtoOALightSymbol
}

type AccountModel struct {
	Symbols      []SymbolModel
	SwingTrader  Trader
	DepositAsset mappers.SwingAsset
	Equity       float64
	Balance      float64
	Positions    []MarketOrderModel
}

func (accountModel *AccountModel) UpdateStatus() {

	var sum = 0.0
	for _, position := range accountModel.Positions {
		sum = sum + position.NetProfit
	}
	accountModel.Equity = float64(sum)
}

type SymbolModel struct {
	Bid               float64
	Ask               float64
	LightSymbol       mappers.SwingLightSymbol
	BaseAsset         mappers.SwingAsset
	QuoteAsset        mappers.SwingAsset
	ConversionSymbols []SymbolModel
	Id                int64
	Data              *messages.ProtoOASymbol
}

func (symbolmodel *SymbolModel) OnTick(bid float64) {
	fmt.Println("Called:", bid)
	symbolmodel.Bid = bid

}

type MarketOrderModel struct {
	GrossProfit              float64
	NetProfit                float64
	Pips                     float64
	Ordermodel               OrderModel
	Swap                     int64
	SwapMonetary             float64
	MoneyDigits              int32
	Commision                int64
	CommissionMonetary       float64
	DoubleCommissionMonetary float64
}

type OrderModel struct {
	Symbol    SymbolModel
	TradeSide messages.ProtoOATradeSide
	TradeData messages.ProtoOATradeData
	Volume    int64
	Price     float64
	OpenTime  int64
}

type Tuple struct {
	BaseAsset  mappers.SwingAsset
	QuoteAsset mappers.SwingAsset
	Bid        float64
}

func (marketOrder *MarketOrderModel) Update(symbol SymbolModel, position messages.ProtoOAPosition) {
	marketOrder.Ordermodel.Symbol = symbol
	marketOrder.Ordermodel.TradeSide = *position.TradeData.TradeSide
	marketOrder.Ordermodel.Volume = *position.TradeData.Volume
	marketOrder.Ordermodel.TradeData = *position.TradeData
}

func (marketOrder *MarketOrderModel) OnSymbolTick() {
	if marketOrder.Ordermodel.TradeSide == messages.ProtoOATradeSide_SELL {
		fmt.Println("Sell")
	} else {
		fmt.Println("Buy")
	}
}
