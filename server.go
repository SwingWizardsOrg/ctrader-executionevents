package main

import (
	"ctrader_events/client"
	"ctrader_events/database"
	"ctrader_events/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

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

}
