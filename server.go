package main

import (
	"crypto/tls"
	"ctraderapi/github.com/Carlosokumu/messages"
	"encoding/binary"
	"fmt"
	"net"
	"os"

	"google.golang.org/protobuf/proto"
)

func main() {

	clientId := "4082_lBzetXS3g9sTx3XROfYcZ6g3vzYk4c6k87WG1EpjTdeMFJQlsT"
	clientSecret := "hVYWOjeVAd9eIfuoqGUV0sBJ4l3jMsotFYNBY22R2ieYmNbKes"

	typepayload := uint32(2100)

	authReq := &messages.ProtoOAApplicationAuthReq{
		ClientId:     &clientId,
		ClientSecret: &clientSecret,
	}
	messageId := "myid"
	// change the authReqmessage to a byte array
	authReqBytes, err := proto.Marshal(authReq)
	if err != nil {
		fmt.Println("error marshaling message:", err)
		return
	}
	message := &messages.ProtoMessage{
		PayloadType: &typepayload,
		Payload:     authReqBytes,
		ClientMsgId: &messageId,
	}

	//listener, err := net.Listen("tcp", "5035")
	conn, err := net.Dial("tcp", "demo-sg-1.ctraderapi.com:5035")

	// create a new TLS client with the default settings
	tlsConn := tls.Client(conn, &tls.Config{
		ServerName: "demo-sg-1.ctraderapi.com",
	})

	// perform the TLS handshake with the server
	if err := tlsConn.Handshake(); err != nil {
		// handle error
		fmt.Println("Handshake Error", err)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	// change the authReqmessage to a byte array
	messageReqBytes, err := proto.Marshal(message)
	if err != nil {
		fmt.Println("error marshaling message:", err)
		return
	}

	fmt.Println("messagebytes:", messageReqBytes)
	// get the length of the byte array
	length := uint32(len(messageReqBytes))
	fmt.Println("Bytelenght:", length)

	lengthBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBytes, length)
	fmt.Println("Bytes:", lengthBytes)

	// reverse the length bytes
	for i, j := 0, len(lengthBytes)-1; i < j; i, j = i+1, j-1 {
		lengthBytes[i], lengthBytes[j] = lengthBytes[j], lengthBytes[i]
	}
	fmt.Println("reversed:", lengthBytes)

	// concatenate the length and message byte slices
	messageBytes := append(lengthBytes, messageReqBytes...)
	fmt.Println("combined", messageBytes)

	// send the message to the API
	_, err = conn.Write(messageReqBytes)
	if err != nil {
		fmt.Println("error sending message:", err)
		return
	}

	c := make(chan *messages.ProtoMessage)

	go handleProtoMessage(conn, c)

	fmt.Println(<-c)

	//go handleProtoMessage(conn, c)

	//go handleProtoMessage(conn, c)

}

func reverseBytes(arr []byte) {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func handleProtoMessage(conn net.Conn, c chan *messages.ProtoMessage) {

	buf := make([]byte, 4)
	_, readerr := conn.Read(buf)
	if readerr != nil {
		// handle error
		fmt.Println(readerr)
	}
	fmt.Println("received:", buf)
	reverseBytes(buf)
	fmt.Println("receivedrev:", buf)

	// convert the byte array to an integer
	num := binary.BigEndian.Uint32(buf)
	fmt.Println("intreceived:", num)

	bufread := make([]byte, num)

	code, readerrr := conn.Read(bufread)
	if readerrr != nil {
		fmt.Println(readerrr)
	}
	//fmt.Println(bufread)
	fmt.Println(code)
	pdata := new(messages.ProtoMessage)
	///rotoMessage := &messages.ProtoMessage{}
	berr := proto.Unmarshal(bufread, pdata)
	if berr != nil {
		fmt.Println("error unmarshaling message:", berr)
		return
	}
	c <- pdata

}
