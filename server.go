package main

import (
	"ctraderapi/messagehandler"
	"ctraderapi/middlewares"
	"ctraderapi/persistence"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {

	dsn := "host=localhost user=postgres password=Agent047 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	persistence.Connect(dsn)

	router := initServer()
	router.Run(":" + "8080")

	var second string
	fmt.Scanln(&second)
}

func initServer() *gin.Engine {
	hub := middlewares.NewHub()
	go hub.Run()
	router := gin.Default()

	//[Websocket] Echo Endpoint ------
	router.GET("/ws", func(c *gin.Context) {
		messagehandler.ConnectToOpen("demo.ctraderapi.com", 5035, hub, c.Writer, c.Request)
	})
	return router
}
