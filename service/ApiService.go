package service

import (
	"ctraderapi/helpers"
	"ctraderapi/mappers"
	"ctraderapi/messages/github.com/Carlosokumu/messages"
	"ctraderapi/middlewares"
	"ctraderapi/models"
	"ctraderapi/persistence"
	"fmt"
	"log"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

var (
	MessageType = 2
)

func AuthorizeApp(conn *websocket.Conn, h *middlewares.Hub) {
	clientId := helpers.ClientId
	clientSecret := helpers.ClientSecret
	authReq := &messages.ProtoOAApplicationAuthReq{
		ClientId:     &clientId,
		ClientSecret: &clientSecret,
	}
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_APPLICATION_AUTH_REQ)
	messageId := "APP_AUTH_REQ"

	authReqBytes, err := proto.Marshal(authReq)
	if err != nil {
		log.Fatal(err)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     authReqBytes,
		ClientMsgId: &messageId,
	}
	protoMessage, _ := proto.Marshal(message)

	// Serialize the message to a byte slice
	writeerror := conn.WriteMessage(MessageType, protoMessage)

	if writeerror != nil {
		log.Fatal(writeerror)
	}

	go func() {

		appauthres := <-h.AppAuthResChannnel
		//Means app is  authorized,we can now authorize the Trading account
		if *appauthres.PayloadType == 2101 {
			fmt.Println("Setting  App Account auth...")
			appAuth := AppAuth{}
			accountAuth := &AccountAuth{}
			appAuth.SetNext(accountAuth)
			accountAuth.Execute(conn, h)
		}

	}()

}

func AuthorizeAccount(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_REQ)
	id := helpers.AccountId
	token := helpers.AccessToken
	messageId := "A/C_AUTH_REQ"
	acReq := &messages.ProtoOAAccountAuthReq{
		CtidTraderAccountId: &id,
		AccessToken:         &token,
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
	go func() {
		accounthAuthRes := <-h.AccountAuthResChannnel
		authRes := messages.ProtoOAAccountAuthRes{}

		err := proto.Unmarshal(accounthAuthRes.Payload, &authRes)
		if err != nil {
			log.Fatal(err)
		}
		if *accounthAuthRes.PayloadType == uint32(messages.ProtoOAPayloadType_PROTO_OA_ACCOUNT_AUTH_RES) {
			accounthAuth := AccountAuth{}

			swingassests := persistence.GetAllSwingAssets()
			lightsymbols := persistence.GetAllSwingLightSymmbol()

			if len(swingassests) == 0 && len(lightsymbols) == 0 {
				assetlistinitializer := &AssetListInitializer{}
				//
				accounthAuth.SetNext(assetlistinitializer)
				assetlistinitializer.Execute(conn, h)
			} else {

				traderinfo := &TraderInfo{}
				accounthAuth.SetNext(traderinfo)
				traderinfo.Execute(conn, h)
			}

		}

	}()

}

func GetTrader(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_TRADER_REQ)
	id := helpers.AccountId
	messageId := "TRADER_REQ"
	acrequest := &messages.ProtoOATraderReq{
		CtidTraderAccountId: &id,
	}
	acBytes, peer := proto.Marshal(acrequest)
	if peer != nil {
		fmt.Println(peer)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     acBytes,
		ClientMsgId: &messageId,
	}
	protomessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protomessage)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		traderInfoRes := <-h.TraderResChannnel
		traderRes := messages.ProtoOATraderRes{}
		proto.Unmarshal(traderInfoRes.Payload, &traderRes)
		trader := TraderInfo{}
		marketorder := &MarketOrder{}
		trader.SetNext(marketorder)
		marketorder.Execute(conn, h)
		h.TraderResChannnel <- traderInfoRes

	}()
}

func GetAccountAssets(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_ASSET_LIST_REQ)
	fmt.Println(payloadtype)
	id := helpers.AccountId
	nmess := "ASSETS_REQ"
	assetReq := &messages.ProtoOAAssetListReq{
		CtidTraderAccountId: &id,
	}
	assetBytes, peer := proto.Marshal(assetReq)
	if peer != nil {
		fmt.Println(peer)
	}

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     assetBytes,
		ClientMsgId: &nmess,
	}

	protoMessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protoMessage)

	if err != nil {
		log.Fatal(err)
	}

	go func() {
		assetListRes := <-h.AssetListChannnel
		accountAssets := messages.ProtoOAAssetListRes{}
		err := proto.Unmarshal(assetListRes.Payload, &accountAssets)
		if err != nil {
			log.Fatal(err)
		}
		swingwizardassets := persistence.GetAllSwingAssets()
		if len(swingwizardassets) == 0 {
			for _, asset := range accountAssets.Asset {
				swingAsset := mappers.SwingAsset{
					AssetId:     *asset.AssetId,
					Name:        *asset.Name,
					DisplayName: *asset.DisplayName,
				}
				persistence.InsertSwingAsset(swingAsset)

			}
		}
		assetListInitializer := AssetListInitializer{}
		lightSymbolInitializer := &LightSymbolInitializer{}
		assetListInitializer.SetNext(lightSymbolInitializer)
		lightSymbolInitializer.Execute(conn, h)

	}()

}

func GetAccountOrders(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_RECONCILE_REQ)
	id := helpers.AccountId
	messageId := "ORDERS_REQ"
	assetReq := &messages.ProtoOAReconcileReq{
		CtidTraderAccountId: &id,
	}

	acBytes, peer := proto.Marshal(assetReq)
	if peer != nil {
		fmt.Println(peer)
	}

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     acBytes,
		ClientMsgId: &messageId,
	}
	protomessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protomessage)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		marketorders := <-h.MarketOrderListChannnel
		reconcileRes := messages.ProtoOAReconcileRes{}
		proto.Unmarshal(marketorders.Payload, &reconcileRes)
		marketorder := MarketOrder{}
		symbol := &Symbol{}
		marketorder.SetNext(symbol)
		symbol.Execute(conn, h)
		h.AccountOrdersChannel <- reconcileRes
	}()
}

func GetLightSymbolList(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_REQ)
	messageId := "LIGHT_SYMBOLS"
	accountId := helpers.AccountId
	IncludeArchivedSymbols := false

	allSymbolsReq := &messages.ProtoOASymbolsListReq{
		CtidTraderAccountId:    &accountId,
		IncludeArchivedSymbols: &IncludeArchivedSymbols,
	}
	symbolBytes, peer := proto.Marshal(allSymbolsReq)
	if peer != nil {
		fmt.Println(peer)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &messageId,
	}
	protoMessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protoMessage)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		lightsymbols := <-h.LightSymbolChannel
		ctraderSymbols := messages.ProtoOASymbolsListRes{}
		SwingLightSymbol := persistence.GetAllSwingLightSymmbol()
		err := proto.Unmarshal(lightsymbols.Payload, &ctraderSymbols)
		if err != nil {
			log.Fatal(err)
		}
		if len(SwingLightSymbol) == 0 {
			for _, symbol := range ctraderSymbols.Symbol {
				lightSymbol := mappers.SwingLightSymbol{
					SymbolId:     *symbol.SymbolId,
					SymbolName:   *symbol.SymbolName,
					BaseAssetId:  *symbol.BaseAssetId,
					QuoteAssetId: *symbol.QuoteAssetId,
				}
				persistence.InsertSwingLightSymbol(lightSymbol)

			}
		}
		lightSymbolInitializer := LightSymbolInitializer{}

		traderinfo := &TraderInfo{}
		lightSymbolInitializer.SetNext(traderinfo)
		traderinfo.Execute(conn, h)

	}()
}

func GetSymbols(conn *websocket.Conn, h *middlewares.Hub) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_REQ)
	messageId := "SYMBOLS"
	accountId := helpers.AccountId
	var symbolIds []int64

	lightsymbols := persistence.GetAllSwingLightSymmbol()
	swingassets := persistence.GetAllSwingAssets()

	for _, lightsymbol := range lightsymbols {
		symbolIds = append(symbolIds, lightsymbol.SymbolId)
	}

	allSymbolsReq := &messages.ProtoOASymbolByIdReq{
		CtidTraderAccountId: &accountId,
		SymbolId:            symbolIds,
	}
	symbolBytes, peer := proto.Marshal(allSymbolsReq)
	if peer != nil {
		fmt.Println(peer)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &messageId,
	}
	protoMessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(2, protoMessage)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		symbols := <-h.Symbols

		symbolRes := messages.ProtoOASymbolByIdRes{}
		proto.Unmarshal(symbols.Payload, &symbolRes)

		symbolmodels := ProcessLightSymbols(lightsymbols, symbolRes.Symbol, swingassets)

		h.SymbolModelChannel <- symbolmodels
		for _, symbolmodel := range symbolmodels {

			//Not Equal to account Deposit Asset
			if symbolmodel.QuoteAsset.AssetId != 11 {
				fmt.Println("Not equal so re-tech...")
				GetConversionSymbols(symbolmodel.QuoteAsset.AssetId, 11, conn)
				lightsymbol := <-h.ConversionLightSymbols
				lightSymbolResponse := messages.ProtoOASymbolsForConversionRes{}
				err := proto.Unmarshal(lightsymbol.Payload, &lightSymbolResponse)
				if err != nil {
					log.Fatal(err)
				}
				accountConversionSymbols := HandleLightSymbols(lightSymbolResponse.Symbol, symbolmodels)

				go func() {
					h.AccounConversionSymbolsChannel <- accountConversionSymbols
				}()
				break

			} else {
				var symbols []models.SymbolModel
				go func() {
					symbols = append(symbols, symbolmodel)
					h.AccounConversionSymbolsChannel <- symbols
				}()

			}

		}
		symbol := Symbol{}
		spotsubscriber := &SpotSubscriber{}
		symbol.SetNext(spotsubscriber)
		spotsubscriber.Execute(conn, h)
	}()
}

func SendSubscribeSpotsRequest(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_SPOTS_REQ)
	id := helpers.AccountId
	ids := []int64{41, 3}
	nmess := "SUB_REQ"

	symbolsRequest := &messages.ProtoOASubscribeSpotsReq{
		CtidTraderAccountId: &id,
		SymbolId:            ids,
	}
	symbolBytes, eer := proto.Marshal(symbolsRequest)
	if eer != nil {
		log.Fatal(eer)
	}

	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &nmess,
	}
	protomessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protomessage)
	if err != nil {
		log.Fatal(err)
	}
}

func GetConversionSymbols(baseAssetId int64, quoteAssetId int64, conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ)
	nmess := "CONVERSION_REQ"
	id := helpers.AccountId

	conversionreq := &messages.ProtoOASymbolsForConversionReq{
		CtidTraderAccountId: &id,
		FirstAssetId:        &baseAssetId,
		LastAssetId:         &quoteAssetId,
	}
	convBytes, peer := proto.Marshal(conversionreq)
	if peer != nil {
		fmt.Println(peer)
	}
	message := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     convBytes,
		ClientMsgId: &nmess,
	}
	protomessage, _ := proto.Marshal(message)
	err := conn.WriteMessage(MessageType, protomessage)
	if err != nil {
		log.Fatal(err)
	}
}

func ProcessLightSymbols(lightSymbols []mappers.SwingLightSymbol, symbols []*messages.ProtoOASymbol, assets []mappers.SwingAsset) []models.SymbolModel {
	var result []models.SymbolModel

	for _, lightSymbol := range lightSymbols {
		for _, symbol := range symbols {
			if lightSymbol.SymbolId == *symbol.SymbolId {
				data := &messages.ProtoOASymbol{}
				baseAsset := assets[0]
				quoteAsset := assets[0]

				for _, s := range symbols {
					if *s.SymbolId == lightSymbol.SymbolId {
						data = s
						break
					}
				}

				for _, asset := range assets {
					if asset.AssetId == lightSymbol.BaseAssetId {
						baseAsset = asset
					}
					if asset.AssetId == lightSymbol.QuoteAssetId {
						quoteAsset = asset
					}
				}
				symbolModel := models.SymbolModel{}
				symbolModel.LightSymbol = lightSymbol
				symbolModel.BaseAsset = baseAsset
				symbolModel.QuoteAsset = quoteAsset
				symbolModel.Id = lightSymbol.SymbolId
				symbolModel.Data = data

				result = append(result, symbolModel)
			}
		}
	}

	return result
}

func HandleLightSymbols(conversionLightSymbols []*messages.ProtoOALightSymbol, accountsymbols []models.SymbolModel) []models.SymbolModel {
	var conversionSymbolModels []models.SymbolModel
	for _, iLightSymbol := range conversionLightSymbols {
		for _, iSymbol := range accountsymbols {
			if iSymbol.Id == *iLightSymbol.SymbolId {
				conversionSymbolModels = append(conversionSymbolModels, iSymbol)
				break
			}
		}
	}
	return conversionSymbolModels
}
