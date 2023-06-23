package main

import (
	"ctrader_events/client"
	"fmt"
)

func main() {

	var amount int

	client.ConnectToCtrader("demo.ctraderapi.com", 5035)
	// taking input and storing in variable using the buffer string
	fmt.Scanf("%s ", &amount)

}
