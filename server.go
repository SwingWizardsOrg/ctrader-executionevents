package main

import (
<<<<<<< HEAD
	"ctraderapi/messagehandler"
	"ctraderapi/middlewares"

	"github.com/gin-gonic/gin"
=======
	"ctrader_events/client"
	"ctrader_events/database"
	"ctrader_events/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
>>>>>>> 96c65c03dc14288996da5e4a6b2ee606824552ee
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
