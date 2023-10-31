package main

import (
	"ctrader_events/messagehandler"
	"ctrader_events/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	router := initServer()
	router.Run(":" + "8070")

}

func initServer() *gin.Engine {
	hub := middlewares.NewHub()
	go hub.Run()
	router := gin.Default()

	//[Websocket] Echo Endpoint ------
	router.GET("/ws", func(c *gin.Context) {
		messagehandler.ConnectToOpen("live.ctraderapi.com", 5035, hub, c.Writer, c.Request)
	})
	return router
}
