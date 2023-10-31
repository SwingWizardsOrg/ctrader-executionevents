package messagebroker

import (
	"ctrader_events/credentials"
	"ctrader_events/database"
	"ctrader_events/helpers"
	"ctrader_events/messages/github.com/swingwizards/messages"
	"fmt"

	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"google.golang.org/protobuf/proto"
)

const (
	MessageType = 2
	// Time allowed to read the next pong message from the peer.
	pongWait = 35 * time.Second

	STANDARD_10000 = 100000
	STANDARD_20000 = 2000000
	STANDARD_30000 = 3000000
	STANDARD_50000 = 5000000
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
			GetTrader(h.Conn)
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_RES):
		{
			//update the master account details
			traderRes := messages.ProtoOATraderRes{}
			proto.Unmarshal(protomessage.Payload, &traderRes)
			// masterAccount := database.MasterAccount{
			// 	Balance:      traderRes.Trader.Balance,
			// 	AccountLogin: uint(*traderRes.Trader.TraderLogin),
			// }
			fmt.Println(traderRes)
			//database.InsertMasterAccountDetails(masterAccount)

		}
	case uint32(messages.ProtoPayloadType_HEARTBEAT_EVENT):
		{
			select {
			case <-time.After(10 * time.Second):
				// Send Back a heartBeat Message to the Server to keep the connection Alive.
				SendHeartBeatMessage(h.Conn)

			}
		}
	case uint32(messages.ProtoPayloadType_ERROR_RES):
		{

		}

	case uint32(messages.ProtoOAPayloadType_PROTO_OA_ERROR_RES):
		{
			panic("Ctrader Error response...")
		}
	case uint32(messages.ProtoOAPayloadType_PROTO_OA_EXECUTION_EVENT):
		{
			fmt.Println("Handle execution events..")
			eventProtoMessage := messages.ProtoOAExecutionEvent{}
			err := proto.Unmarshal(protomessage.Payload, &eventProtoMessage)

			if err != nil {
				log.Fatal(err)
			}
			if *eventProtoMessage.ExecutionType == messages.ProtoOAExecutionType_ORDER_FILLED {
				symbolId64 := *eventProtoMessage.Position.TradeData.SymbolId
				symbolId := strconv.Itoa(int(symbolId64))
				symbolEntity, err := database.GetSymbolEntity(symbolId)
				if err != nil {
					log.Fatal(err)
				}

				tradeInfo := *eventProtoMessage.Position.TradeData
				orderInfo := eventProtoMessage.Order
				switch *tradeInfo.Volume {
				case STANDARD_10000:
					{

						tradeSide := helpers.DetermineTradeSide(tradeInfo.TradeSide)
						positionReward := helpers.GetTakeProfitPips(orderInfo.RelativeTakeProfit, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroOne)
						positionRisk := helpers.GetTakeProfitPips(orderInfo.RelativeStopLoss, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroOne)

						fmt.Println("PIPS:", helpers.GetTakeProfitPips(orderInfo.RelativeTakeProfit, tradeInfo.Volume))
						fmt.Println("PIPSNN:", *symbolEntity.Lot_ZeroOne)
						fmt.Println("POSITIONREWARD:", positionReward)
						fmt.Println("POSITIONREWARD:", positionReward)
						fmt.Println("POSITIONRISK:", positionRisk)

						runningPosition := &database.RunningPosition{
							Volume:          tradeInfo.Volume,
							Price:           eventProtoMessage.Position.Price,
							TradeSide:       &tradeSide,
							Commission:      eventProtoMessage.Position.Commission,
							MoneyDigits:     eventProtoMessage.Position.MoneyDigits,
							Swap:            eventProtoMessage.Position.Swap,
							OpenTime:        eventProtoMessage.Position.GetTradeData().OpenTimestamp,
							SymbolId:        tradeInfo.SymbolId,
							PositionsReward: &positionReward,
							PositionRisk:    &positionRisk,
						}
						users := database.GetAllUsers()
						for _, user := range users {
							userRisk := float64(positionRisk)
							userRiskProportion := *user.PercentageContribution * userRisk
							userBalance := *user.Balance
							user.Positions = append(user.Positions, *runningPosition)
							newBalance := userBalance - userRiskProportion
							user.Balance = &newBalance
							database.SaveUser(user)
						}

					}
				case STANDARD_20000:
					{

						tradeSide := helpers.DetermineTradeSide(tradeInfo.TradeSide)
						positionReward := helpers.GetTakeProfitPips(orderInfo.RelativeTakeProfit, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroTwenty)
						positionRisk := helpers.GetTakeProfitPips(orderInfo.RelativeStopLoss, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroTwenty)
						fmt.Println("POSITIONREWARD:", positionReward)
						fmt.Println("POSITIONREWARD:", positionReward)
						runningPosition := &database.RunningPosition{
							Volume:          tradeInfo.Volume,
							Price:           eventProtoMessage.Position.Price,
							TradeSide:       &tradeSide,
							Commission:      eventProtoMessage.Position.Commission,
							MoneyDigits:     eventProtoMessage.Position.MoneyDigits,
							Swap:            eventProtoMessage.Position.Swap,
							OpenTime:        eventProtoMessage.Position.GetTradeData().OpenTimestamp,
							SymbolId:        tradeInfo.SymbolId,
							PositionsReward: &positionReward,
							PositionRisk:    &positionRisk,
						}

						users := database.GetAllUsers()
						for _, user := range users {

							user.Positions = append(user.Positions, *runningPosition)
							database.SaveUser(user)
						}

					}
				case STANDARD_30000:
					{

						tradeSide := helpers.DetermineTradeSide(tradeInfo.TradeSide)
						positionReward := helpers.GetTakeProfitPips(orderInfo.RelativeTakeProfit, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroThirty)
						positionRisk := helpers.GetTakeProfitPips(orderInfo.RelativeStopLoss, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroThirty)
						fmt.Println("POSITIONREWARD:", positionReward)
						fmt.Println("POSITIONREWARD:", positionReward)
						runningPosition := &database.RunningPosition{
							Volume:          tradeInfo.Volume,
							Price:           eventProtoMessage.Position.Price,
							TradeSide:       &tradeSide,
							Commission:      eventProtoMessage.Position.Commission,
							MoneyDigits:     eventProtoMessage.Position.MoneyDigits,
							Swap:            eventProtoMessage.Position.Swap,
							OpenTime:        eventProtoMessage.Position.GetTradeData().OpenTimestamp,
							SymbolId:        tradeInfo.SymbolId,
							PositionsReward: &positionReward,
							PositionRisk:    &positionRisk,
						}

						users := database.GetAllUsers()
						for _, user := range users {
							user.Positions = append(user.Positions, *runningPosition)
							database.SaveUser(user)
						}

					}
				case STANDARD_50000:
					{

						tradeSide := helpers.DetermineTradeSide(tradeInfo.TradeSide)
						positionReward := helpers.GetTakeProfitPips(orderInfo.RelativeTakeProfit, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroFive)
						positionRisk := helpers.GetTakeProfitPips(orderInfo.RelativeStopLoss, tradeInfo.Volume) * (*symbolEntity.Lot_ZeroFive)
						fmt.Println("POSITIONREWARD:", positionReward)
						fmt.Println("POSITIONREWARD:", positionReward)
						runningPosition := &database.RunningPosition{
							Volume:          tradeInfo.Volume,
							Price:           eventProtoMessage.Position.Price,
							TradeSide:       &tradeSide,
							Commission:      eventProtoMessage.Position.Commission,
							MoneyDigits:     eventProtoMessage.Position.MoneyDigits,
							Swap:            eventProtoMessage.Position.Swap,
							OpenTime:        eventProtoMessage.Position.GetTradeData().OpenTimestamp,
							SymbolId:        tradeInfo.SymbolId,
							PositionsReward: &positionReward,
							PositionRisk:    &positionRisk,
						}

						users := database.GetAllUsers()
						for _, user := range users {
							user.Positions = append(user.Positions, *runningPosition)
							database.SaveUser(user)
						}

					}
				}
				fmt.Println("RateConversion:", *symbolEntity.Lot_ZeroFive)
			}
			fmt.Println(eventProtoMessage)

		}
	default:
		{
			fmt.Println("Running default")
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

func GetTrader(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_REQ)
	accountId := credentials.AccountId
	messageId := "TRADER_REQ"
	traderrequest := &messages.ProtoOATraderReq{
		CtidTraderAccountId: &accountId,
	}
	traderRequestBytes, err := proto.Marshal(traderrequest)
	if err != nil {
		fmt.Println(err)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     traderRequestBytes,
		ClientMsgId: &messageId,
	}
	protomessage, err := proto.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.WriteMessage(MessageType, protomessage)
	if err != nil {
		log.Fatal(err)
	}

}

func HandleDatabaseResponse(users []database.User) {
	for _, user := range users {
		//newPosition := database.RunningPosition{UserID: 24}
		///user.Positions = append(user.Positions, newPosition)
		database.SaveUser(user)
		fmt.Println(user.Username)

		fmt.Println("User ID:", database.Instance.Model(&database.User{}).Preload("Positions").First(&user))

		fmt.Println("USERPOST:", user.Positions)
	}
}
