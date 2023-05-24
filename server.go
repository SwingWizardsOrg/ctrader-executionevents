package main

import (
	"ctraderapi/messagehandler"
	"ctraderapi/middlewares"
	"ctraderapi/persistence"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// var then variable name then variable type
	//var first string
	// initServer()
	dsn := "host=localhost user=postgres password=Agent047 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	persistence.Connect(dsn)
	// persistence.Migrate()
	// len := len(persistence.GetAllSwingAssets())
	// fmt.Println("len:", len)

	//badgerconnection := persistence.CreateBadgerConnection()

	//persistence.InsertSymbolData()
	// symboldata, err := persistence.ReadSymbolData(int64(1))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("symbolhere", symboldata)

	//lightsymbol, _ := persistence.ReadLightSymbolData(int64(4))
	// symboldata, _ := persistence.ReadSymbolData(int64(2))

	// // // if err != nil {
	// // // 	fmt.Println(err)
	// // // }
	//fmt.Println("LightHere:", lightsymbol)
	// fmt.Println("symboldata:", symboldata)

	router := initServer()
	router.Run(":" + "8080")
	// 	var second string
	// 	fmt.Scanln(&second)
	// 	fmt.Scanln(&first)
	// hub := middlewares.NewHub()
	// go hub.Run()

	// messagehandler.ConnectToOpen("demo.ctraderapi.com", 5035, hub)

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

// 	// return router

// }
