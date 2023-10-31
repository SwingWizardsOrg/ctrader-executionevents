package main

import (
	"ctraderapi/messagehandler"
	"ctraderapi/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	//dsn := os.Getenv("DATABASE_URL")

	// dsn := "postgres://carlos:vUUYROlx74jmAdnvVunkNqiNdxAvZI32@dpg-cg9il4pmbg54mbfbrte0-a.oregon-postgres.render.com/swings"
	// persistence.Connect(dsn)
	//persistence.Migrate()
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
