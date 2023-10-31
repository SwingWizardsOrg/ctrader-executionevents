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

<<<<<<< HEAD
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
=======
	var amount int

	databaseUrl := "postgres://swinngdata_user:4nZcOypBKc8E6RU96BftsnBMgClMGxqn@dpg-ci805l98g3n3vm2k9ifg-a.oregon-postgres.render.com/swinngdata"
	database.Connect(databaseUrl)
	//database.Migrate()

	// Open our jsonFile
	jsonFile, err := os.Open("symbols.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened symbol's json file")

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var symbols models.Symbols

	json.Unmarshal(byteValue, &symbols)

	client.ConnectToCtrader("demo.ctraderapi.com", 5035)
	fmt.Scanf("%s ", &amount)

>>>>>>> 96c65c03dc14288996da5e4a6b2ee606824552ee
}
