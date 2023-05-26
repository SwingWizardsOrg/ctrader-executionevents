package service

import (
	"fmt"
	"log"
	"math"
	"time"

	"ctraderapi/helpers"
	"ctraderapi/mappers"
	"ctraderapi/messages/github.com/Carlosokumu/messages"
	"ctraderapi/middlewares"
	"ctraderapi/models"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
)

func ReadCtraderMessages(conn *websocket.Conn, handler middlewares.Client) {
	fmt.Println("Reading Messages from Ctrader....🧔🏽‍♂️")
	defer func() {
		conn.Close()
	}()

	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msg := &messages.ProtoMessage{}
		_, readmessage, err := conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		unmarsherr := proto.Unmarshal(readmessage, msg)

		if unmarsherr != nil {
			fmt.Println(unmarsherr)
		}

		handler.Hub.Protos <- messages.ProtoMessage{
			PayloadType: msg.PayloadType,
			Payload:     msg.Payload,
			ClientMsgId: msg.ClientMsgId,
		}
		fmt.Println(msg)

	}
}

func CollectAllMessages(h *middlewares.Hub, conn *websocket.Conn) {

	// Order Matters else channels will block hence no execution
	symbolmodels := <-h.SymbolModelChannel
	swingtrader := <-h.TraderResChannnel
	accountorders := <-h.AccountOrdersChannel

	accountModel := models.AccountModel{}

	accountModel.Symbols = symbolmodels

	fmt.Println("Traderin:", swingtrader)
	fmt.Println("symbolmodels:", symbolmodels[0])
	fmt.Println("accountorders:", accountorders)

	Trader := messages.ProtoOATraderRes{}
	err := proto.Unmarshal(swingtrader.Payload, &Trader)
	if err != nil {
		log.Fatal(err)
	}
	accountModel.DepositAsset = mappers.SwingAsset{
		AssetId: *Trader.Trader.DepositAssetId,
	}
	accountModel.SwingTrader = models.Trader{
		Balance:     *Trader.Trader.Balance,
		TraderLogin: *Trader.Trader.TraderLogin,
	}

	var positions []models.MarketOrderModel
	for _, position := range accountorders.Position {
		var positionSymbol models.SymbolModel
		for _, symbol := range accountModel.Symbols {
			if symbol.Id == *position.TradeData.SymbolId {
				positionSymbol = symbol
				break
			}

		}
		mareketOrderModel := models.MarketOrderModel{}
		mareketOrderModel.Ordermodel.Symbol = positionSymbol
		mareketOrderModel.Ordermodel.TradeSide = *position.TradeData.TradeSide
		mareketOrderModel.Ordermodel.Price = *position.Price
		mareketOrderModel.Ordermodel.Volume = *position.TradeData.Volume
		mareketOrderModel.Commision = *position.Commission
		mareketOrderModel.Swap = *position.Swap
		mareketOrderModel.MoneyDigits = int32(*position.MoneyDigits)
		mareketOrderModel.CommissionMonetary = float64(mareketOrderModel.Commision / int64(math.Pow(10, float64(mareketOrderModel.MoneyDigits))))
		mareketOrderModel.DoubleCommissionMonetary = mareketOrderModel.CommissionMonetary * 2
		power := math.Pow(10, float64(mareketOrderModel.MoneyDigits))
		SwapMoney := float64(*position.Swap) / float64(power)
		mareketOrderModel.SwapMonetary = SwapMoney
		positions = append(positions, mareketOrderModel)
	}
	accountModel.Positions = positions

	go func() {
		event := <-h.SpotEventChannel
		spotEvent := &messages.ProtoOASpotEvent{}
		unmarsherr := proto.Unmarshal(event.Payload, spotEvent)
		if unmarsherr != nil {
			log.Fatal(unmarsherr)
		}
		var (
			bid, ask, TickValue float64
			symbol              models.SymbolModel
		)
		for _, symbomodel := range accountModel.Symbols {
			if symbomodel.Id == *spotEvent.SymbolId {
				symbol = symbomodel
				break
			}
		}
		conversionSymbols := <-h.AccounConversionSymbolsChannel
		symbol.ConversionSymbols = conversionSymbols

		if spotEvent.Bid != nil {
			bid = helpers.GetPriceRelative(symbol.Data, float64(*spotEvent.Bid))

			symbol.OnTick(bid)
			fmt.Println("bid:", symbol.Bid)
		}

		if spotEvent.Ask != nil {
			ask = helpers.GetPriceRelative(symbol.Data, float64(*spotEvent.Ask))

			fmt.Println("ask:", ask)
		}
		symbol.Bid = bid

		if symbol.QuoteAsset.AssetId == accountModel.DepositAsset.AssetId {
			TickValue = helpers.GetTickValue(symbol.Data, symbol.QuoteAsset, *Trader.Trader.DepositAssetId, nil)
		} else if (len(symbol.ConversionSymbols)) > 0 {
			var updatedSymbols []models.SymbolModel
			var tickValues []models.Tuple
			for _, iSymbol := range symbol.ConversionSymbols {
				iSymbol.Bid = bid
				updatedSymbols = append(updatedSymbols, iSymbol)

			}
			allNonZero := true
			for _, iSymbol := range updatedSymbols {
				if iSymbol.Bid == 0 {
					allNonZero = false
					break
				}
				if allNonZero {
					for _, iSymbol := range updatedSymbols {
						tickValues = append(tickValues, models.Tuple{BaseAsset: iSymbol.BaseAsset, QuoteAsset: iSymbol.QuoteAsset, Bid: iSymbol.Bid})
					}
					TickValue = helpers.GetTickValue(symbol.Data, symbol.QuoteAsset, *Trader.Trader.DepositAssetId, tickValues)

				}

			}

		}
		var positionsPLStatus []float64
		sum := 0.0

		for _, position := range accountModel.Positions {
			if position.Ordermodel.TradeSide == messages.ProtoOATradeSide_BUY {
				Symbol := position.Ordermodel.Symbol.Data
				pipValue := helpers.GetPipValue(*Symbol, TickValue)
				positionReturn := symbol.Bid - position.Ordermodel.Price
				pips := helpers.GetPipsFromPrice(Symbol, positionReturn)
				grossProfit := pips * pipValue * helpers.FromMonetary(float64(position.Ordermodel.Volume))
				roundedGross := math.Round(grossProfit*100) / 100
				netProfit := roundedGross + position.DoubleCommissionMonetary + position.SwapMonetary
				roundedNet := math.Round(netProfit*100) / 100
				positionsPLStatus = append(positionsPLStatus, roundedNet)

			} else {

			}
		}

		// Iterate over the array and add each element to the sum
		for _, netProfit := range positionsPLStatus {
			sum += netProfit
		}
		accountModel.Balance = helpers.FromMonetary(float64(accountModel.SwingTrader.Balance))
		equity := accountModel.Balance + sum
		accountModel.Equity = equity

		h.AccountModelChannel <- accountModel

		fmt.Println("BALANCE:", helpers.FromMonetary(float64(accountModel.SwingTrader.Balance)))
		fmt.Println("Sum:", sum)
		fmt.Println("EQUITY:", equity)

	}()

}
