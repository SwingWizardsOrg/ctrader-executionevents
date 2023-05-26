package main

import (
	"ctraderapi/messagehandler"
	"ctraderapi/middlewares"
	"ctraderapi/persistence"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	dsn := os.Getenv("DATABASE_URL")
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
