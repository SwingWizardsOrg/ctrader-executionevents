package main

import (
	"fmt"
	"net"
	"time"
)

func main() {

	conn, err := net.DialTimeout("tcp", "demo-sg.ctraderapi.com:5035", time.Second*10)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to cTrader server")

}
