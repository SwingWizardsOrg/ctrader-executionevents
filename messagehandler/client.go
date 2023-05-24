package messagehandler

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"ctraderapi/appmessages"
	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"ctraderapi/mappers"
	"ctraderapi/middlewares"
	"ctraderapi/models"
	"ctraderapi/persistence"
	"ctraderapi/service"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

// // Client is a middleman between the websocket connection and the hub.
// type Client struct {
// 	hub *Hub

// 	// ctrader websocket connection.
// 	conn *websocket.Conn

// 	//will hold ctrader protomessages
// 	protomessages chan messages.ProtoMessage

// 	//app's websocket connection
// 	appconn *websocket.Conn

// 	//will hold json messages
// 	resources chan models.ResourceId
// }

const (
	// Time allowed to write a message to the peer.
	writeWait = 60 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024
)

var (
	CanAcessAccount *bool
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Reads messages from  a client websocket connection to the hub
// func (c *Client) readAppMessages() {
// 	fmt.Println("Reading messages from client  app...")
// 	defer func() {
// 		c.hub.unregister <- c
// 		c.appconn.Close()
// 	}()
// 	c.appconn.SetReadLimit(maxMessageSize)
// 	c.appconn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.appconn.SetPongHandler(func(string) error { c.appconn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

// 	for {
// 		var resourceId models.ResourceId
// 		err := c.appconn.ReadJSON(&resourceId)
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		fmt.Println(resourceId.ResourceName)
// 		c.hub.resourceid <- resourceId
// 	}

// }

// Reads messages from the ctrader websocket connection to the hub.
// func (c *Client) readCtraderMessages() {
// 	fmt.Println("Reading from ctrader....")
// 	defer func() {
// 		c.conn.Close()
// 	}()
// 	// c.conn.SetReadLimit(maxMessageSize)
// 	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	// c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
// 	for {
// 		msg := &messages.ProtoMessage{}
// 		_, readmessage, err := c.conn.ReadMessage()

// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		unmarsherr := proto.Unmarshal(readmessage, msg)

// 		if unmarsherr != nil {
// 			fmt.Println(unmarsherr)
// 		}
// 		fmt.Println(msg)
// 		c.hub.protos <- messages.ProtoMessage{
// 			PayloadType: msg.PayloadType,
// 			Payload:     msg.Payload,
// 			ClientMsgId: msg.ClientMsgId,
// 		}
// 	}
// }

// writePump pumps messages from the hub to the  app's websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// // executing all writes from this goroutine.
func writePump(c *middlewares.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Appconn.Close()
	}()
	for {
		select {

		case <-ticker.C:
			c.Appconn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Appconn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		case accountModel, ok := <-c.Hub.AccountModelChannel:

			accountModelUsecase := models.AccountModelUseCase{}
			c.Appconn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Appconn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			accountModelUsecase.Balance = accountModel.Balance
			accountModelUsecase.Equity = accountModel.Equity
			accountModelUsecase.Positions = accountModel.Positions

			err := c.Appconn.WriteJSON(accountModelUsecase)

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// 		case jsonData, ok := <-c.hub.protosback:
// 			//Write the ProtoMessage  as Json back to the app's connection
// 			fmt.Println("Protomessage")
// 			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				// The hub closed the channel.
// 				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			fmt.Println("WHATCURRENT:", jsonData)
// 			if *jsonData.PayloadType == 2103 {
// 				IsAuthenticated = true
// 			}

// 			NextMessage(*jsonData.PayloadType, c.conn)
// 			DeSerialize(*jsonData.PayloadType, jsonData.Payload, c.appconn)
// 			//result := GetMessage(jsonData)
// 			//fmt.Println("Thisresult:", result.getName())

// 		case resource, ok := <-c.resources:
// 			if !ok {
// 				// The hub closed the channel.
// 				c.appconn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}
// 			fmt.Println("passedIn:", resource)
// 			if IsAuthenticated == true {
// 				SendProtoMessage(resource.ResourceName, c.conn)
// 			}

// 		}
// 	}
// }

// func ConnectToOpenAPI(hub *Hub, host string, port int, w http.ResponseWriter, r *http.Request) {

// 	appconn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	// Set up a dialer with the desired options
// 	dialer := websocket.DefaultDialer
// 	dialer.EnableCompression = true

// 	// Connect to the  Ctrader WebSocket endpoint
// 	url := fmt.Sprintf("wss://%s:%d", host, port)
// 	fmt.Println(url)
// 	conn, _, err := dialer.Dial(url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	client := &Client{hub: hub, conn: conn, protomessages: make(chan messages.ProtoMessage), appconn: appconn, resources: make(chan models.ResourceId)}
// 	client.hub.register <- client

// 	//sendMessage(conn)
// 	//sendMessage(conn)
// 	// Allow collection of memory referenced by the caller by doing all work in
// 	// new goroutines.
// 	go client.readAppMessages()
// 	go client.readCtraderMessages()
// 	go client.writePump()

// }

func ConnectToOpen(host string, port int, hub *middlewares.Hub, w http.ResponseWriter, r *http.Request) {

	Appconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Set up a dialer with the desired options
	dialer := websocket.DefaultDialer
	dialer.EnableCompression = true

	// Connect to the  Ctrader WebSocket endpoint
	url := fmt.Sprintf("wss://%s:%d", host, port)
	fmt.Println(url)
	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	client := &middlewares.Client{Hub: hub, Conn: conn, Protomessages: make(chan messages.ProtoMessage), Appconn: Appconn}
	appAuth := &service.AppAuth{}
	accountAuth := &service.AccountAuth{}
	appAuth.SetNext(accountAuth)

	appAuth.Execute(conn, hub)

	go service.ReadCtraderMessages(conn, *client)
	go writePump(client)
	service.CollectAllMessages(hub, conn)

}

func sendMessage(conn *websocket.Conn) {
	clientId := "4082_lBzetXS3g9sTx3XROfYcZ6g3vzYk4c6k87WG1EpjTdeMFJQlsT"
	clientSecret := "hVYWOjeVAd9eIfuoqGUV0sBJ4l3jMsotFYNBY22R2ieYmNbKes"
	authReq := &messages.ProtoOAApplicationAuthReq{
		ClientId:     &clientId,
		ClientSecret: &clientSecret,
	}
	messageId := "myid"
	var x = uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_REQ)
	fmt.Println("x:", x)

	authReqBytes, err := proto.Marshal(authReq)
	if err != nil {
		fmt.Println("error marshaling message:", err)
	}
	message := &messages.ProtoMessage{
		PayloadType: &x,
		Payload:     authReqBytes,
		ClientMsgId: &messageId,
	}
	tobesent, _ := proto.Marshal(message)

	// // Serialize the message to a byte slice
	writeerror := conn.WriteMessage(2, tobesent)

	if writeerror != nil {
		fmt.Println("writeerror", writeerror)
	}

}

func SendTraderRequest(conn *websocket.Conn) {
	fmt.Println("Sending Trading request...")
	var tt = uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_REQ)
	id := int64(25675710)
	nmess := "443"
	acrequest := &messages.ProtoOATraderReq{
		CtidTraderAccountId: &id,
	}
	acBytes, peer := proto.Marshal(acrequest)
	if peer != nil {
		fmt.Println(peer)
	}

	nessage := &messages.ProtoMessage{
		PayloadType: &tt,
		Payload:     acBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	if verr != nil {
		fmt.Println("verr:", verr)
	}

}

// Will Send an Account Auth Request to Ctrader.
func SendAuthRequest(conn *websocket.Conn) {
	var tt = uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ)
	id := int64(25675710)
	token := "uHhSL-hDqtpe47UGm2RnexKDX-CDeChUircGfF7-dQY"
	nmess := "443"
	acrequest := &messages.ProtoOAAccountAuthReq{
		CtidTraderAccountId: &id,
		AccessToken:         &token,
	}
	acBytes, peer := proto.Marshal(acrequest)
	if peer != nil {
		fmt.Println(peer)
	}

	nessage := &messages.ProtoMessage{
		PayloadType: &tt,
		Payload:     acBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	if verr != nil {
		fmt.Println("verr:", verr)
	}

}

func NextMessage(value uint32, conn *websocket.Conn) {
	if value == 2101 {
		SendAuthRequest(conn)
	}
}

func SendProtoMessage(resourceId string, conn *websocket.Conn) {
	switch resourceId {
	case "TRADER_INFO":
		SendTraderRequest(conn)
	case "OPEN_POSITIONS":
		appmessages.SendPositionsRequest(conn)
	case "SYMBOL_INFO":
		appmessages.SendSymbolRequest(conn)
	case "SUBSCRIBE_REQ":
		appmessages.SendSubscribeSpotsRequest(conn)
	case "ASSET_LIST":
		appmessages.SendProtoAssetListReq(conn)
	case "CONVERSION_REQ":
		appmessages.SendProtoOAsymbolConversion(conn)
	case "SYMBOLS":
		appmessages.SendSymbolListRequest(conn)
	default:
	}
}

func DeSerialize(prototype uint32, payload []byte, appconn *websocket.Conn) {
	switch prototype {
	case 2122:
		{
			traderRes := messages.ProtoOATraderRes{}
			proto.Unmarshal(payload, &traderRes)
			traderInfo := models.Trader{
				Balance:     *traderRes.Trader.Balance,
				TraderLogin: *traderRes.Trader.TraderLogin,
			}
			appconn.WriteJSON(traderInfo)
			fmt.Println("trader:", *traderRes.Trader.Balance)
		}

	case 2125:
		{
			positions := messages.ProtoOAReconcileRes{}
			proto.Unmarshal(payload, &positions)
			traderOpenPositions := models.PositionsInfo{
				Positions: positions.Position,
			}
			appconn.WriteJSON(traderOpenPositions)
			fmt.Println("positions:", positions.Position)
		}
	case 2117:
		{
			symbolinfo := messages.ProtoOASymbolByIdRes{}
			proto.Unmarshal(payload, &symbolinfo)
			for _, symbol := range symbolinfo.Symbol {
				fmt.Println("symbol:", symbol)
				swingsymbol := mappers.SwingSymbol{
					Digits:   *symbol.Digits,
					SymbolId: *symbol.SymbolId,
				}
				fmt.Println("Mapped:", swingsymbol.Digits)
				persistence.InsertSwingSymbol(swingsymbol)
			}
			symboldata := models.SymbolInformation{
				Symbols: symbolinfo.Symbol,
			}
			appconn.WriteJSON(symboldata)
			fmt.Println("symbols:", symbolinfo.Symbol)
		}
	case 2128:
		{

		}
	case 2131:
		{
			var bid, ask float64
			// var (
			// 	conversionassets []helpers.ConversionAsset
			// )

			eventInfo := messages.ProtoOASpotEvent{}
			err := proto.Unmarshal(payload, &eventInfo)
			if err != nil {
				fmt.Println("eventerror:", err)
			}
			fmt.Println("REQUETSID:", *eventInfo.SymbolId)
			lightsymbol := persistence.GetSwingLightSymbol(*eventInfo.SymbolId)
			symbolQuoteAsset := persistence.GetSwingAsset(lightsymbol.QuoteAssetId)

			// lightsymbol, _ := persistence.ReadLightSymbolData(*eventInfo.SymbolId)
			fmt.Println("LightHere:", lightsymbol)
			fmt.Println("Assethere:", symbolQuoteAsset)
			if eventInfo.Ask != nil {
				// symbodata, _ := persistence.ReadSymbolData(1)
				// fmt.Println("symbolId:", symbodata.SymbolId)

				// if *symbodata.SymbolId == *eventInfo.SymbolId {
				// 	fmt.Println("databid:", *eventInfo.Ask)
				// 	ask = helpers.GetPriceRelative(symbodata, float64(*eventInfo.Ask))
				// }
				fmt.Println("ask:", ask)
			} else if eventInfo.Bid != nil {
				// symbodata, _ := persistence.ReadSymbolData(1)
				// if *symbodata.SymbolId == *eventInfo.SymbolId {
				// 	bid = helpers.GetPriceRelative(symbodata, float64(*eventInfo.Bid))
				// }
				fmt.Println("bid:", bid)
			}
		}
	case 2113:
		{
			assetInfo := messages.ProtoOAAssetListRes{}
			err := proto.Unmarshal(payload, &assetInfo)
			if err != nil {
				log.Fatal(err)
			}
			for _, asset := range assetInfo.Asset {
				if *asset.AssetId > 25 {
					break
				}
				swingAsset := mappers.SwingAsset{
					AssetId:     *asset.AssetId,
					Name:        *asset.Name,
					DisplayName: *asset.DisplayName,
				}
				persistence.InsertSwingAsset(swingAsset)
			}
			protoasset := models.AssetInfo{
				Assets: assetInfo.Asset,
			}
			appconn.WriteJSON(protoasset)

		}
	case 2119:
		{
			conversionInfo := messages.ProtoOASymbolsForConversionRes{}
			err := proto.Unmarshal(payload, &conversionInfo)
			if err != nil {
				fmt.Println("eventerror:", err)
			}
			protoasset := models.ConversionInfo{
				SymbolChain: conversionInfo.Symbol,
			}
			appconn.WriteJSON(protoasset)

		}
	case 51:
		{
			fmt.Println("HeartBeat..")
		}
	case 2115:
		{
			allSymbols := messages.ProtoOASymbolsListRes{}
			err := proto.Unmarshal(payload, &allSymbols)
			if err != nil {
				fmt.Println("eventerror:", err)
			}
			for _, symbol := range allSymbols.Symbol {
				if *symbol.SymbolId >= 25 {
					break
				}

				lightSymbol := mappers.SwingLightSymbol{
					SymbolId:     *symbol.SymbolId,
					SymbolName:   *symbol.SymbolName,
					BaseAssetId:  *symbol.BaseAssetId,
					QuoteAssetId: *symbol.QuoteAssetId,
				}
				persistence.InsertSwingLightSymbol(lightSymbol)
			}

		}
	default:
		{
			fmt.Println("Unprocessed protomessage:", prototype)
		}

	}

}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
