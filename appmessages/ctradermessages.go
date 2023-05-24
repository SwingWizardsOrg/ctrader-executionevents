package appmessages

import (
	"fmt"

	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"google.golang.org/protobuf/proto"

	"github.com/gorilla/websocket"
)

// Will Send an Account Auth Request to Ctrader.
func SendPositionsRequest(conn *websocket.Conn) {
	var tt = uint32(messages.ProtoOAPayloadType_PROTO_OA_RECONCILE_REQ)
	id := int64(25675710)
	nmess := "443"
	positionsrequest := &messages.ProtoOAReconcileReq{
		CtidTraderAccountId: &id,
	}
	acBytes, peer := proto.Marshal(positionsrequest)
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

// Request Symbol Information By Id on the Ctrader Platform
func SendSymbolRequest(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOL_BY_ID_REQ)
	id := int64(25675710)
	ids := []int64{1, 2}
	nmess := "Symbol_request"

	symbolsRequest := &messages.ProtoOASymbolByIdReq{
		CtidTraderAccountId: &id,
		SymbolId:            ids,
	}
	symbolBytes, peer := proto.Marshal(symbolsRequest)
	if peer != nil {
		fmt.Println(peer)
	}

	nessage := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	if verr != nil {
		fmt.Println("verr:", verr)
	}

}

func SendSubscribeSpotsRequest(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SUBSCRIBE_SPOTS_REQ)
	id := int64(25675710)
	ids := []int64{1, 2}
	nmess := "Subscribe_request"

	symbolsRequest := &messages.ProtoOASubscribeSpotsReq{
		CtidTraderAccountId: &id,
		SymbolId:            ids,
	}
	symbolBytes, peer := proto.Marshal(symbolsRequest)
	if peer != nil {
		fmt.Println(peer)
	}

	nessage := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	if verr != nil {
		fmt.Println("verr:", verr)
	}
}

func SendProtoOAsymbolConversion(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_FOR_CONVERSION_REQ)
	nmess := "Subscribe_request"
	id := int64(25675710)
	firstasset := int64(6)
	lastasset := int64(4)

	conversionreq := &messages.ProtoOASymbolsForConversionReq{
		CtidTraderAccountId: &id,
		FirstAssetId:        &firstasset,
		LastAssetId:         &lastasset,
	}
	convBytes, peer := proto.Marshal(conversionreq)
	if peer != nil {
		fmt.Println(peer)
	}
	nessage := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     convBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	if verr != nil {
		fmt.Println("verr:", verr)
	}

}

func SendProtoAssetListReq(conn *websocket.Conn) {
	//PROTO_OA_ASSET_LIST_REQ
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_ASSET_LIST_REQ)
	fmt.Println(payloadtype)
	id := int64(25675710)
	nmess := "asset_req"
	assetReq := &messages.ProtoOAAssetListReq{
		CtidTraderAccountId: &id,
	}
	assetBytes, peer := proto.Marshal(assetReq)
	if peer != nil {
		fmt.Println(peer)
	}

	nessage := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     assetBytes,
		ClientMsgId: &nmess,
	}

	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	//_, verr := writer.Write(ttooo)
	if verr != nil {
		fmt.Println("assetrequesterr:", verr)
	}

}

// Request for all the Symbol's available in a given Trading account
func SendSymbolListRequest(conn *websocket.Conn) {
	var payloadtype = uint32(messages.ProtoOAPayloadType_PROTO_OA_SYMBOLS_LIST_REQ)
	id := int64(25675710)
	nmess := "symbols_req"
	allSymbolsReq := &messages.ProtoOASymbolsListReq{
		CtidTraderAccountId: &id,
	}
	symbolBytes, peer := proto.Marshal(allSymbolsReq)
	if peer != nil {
		fmt.Println(peer)
	}
	nessage := &messages.ProtoMessage{
		PayloadType: &payloadtype,
		Payload:     symbolBytes,
		ClientMsgId: &nmess,
	}
	ttooo, _ := proto.Marshal(nessage)
	verr := conn.WriteMessage(2, ttooo)
	//_, verr := writer.Write(ttooo)
	if verr != nil {
		fmt.Println("symbolreqerr:", verr)
	}
}
