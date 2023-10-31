package persistence

import (
	"ctrader_events/mappers"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {

	Instance, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if dbError != nil {
		log.Fatal(dbError)
	}
	log.Println("Connected to Database!")
}

func Migrate() {

	err := Instance.AutoMigrate(
		&mappers.SwingLightSymbol{},
		&mappers.SwingSymbol{},
		&mappers.SwingAsset{},
		&SymbolModel{},
	)

	if err != nil {
		log.Println(err)
	}
	log.Println("Database Migration Completed!")
}
