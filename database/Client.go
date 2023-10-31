package database

import (
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
		&SymbolEntity{},
		&RunningPosition{},
		&MasterAccount{},
	)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database Migration Completed!")
}

func GetAllUsers() []User {
	var users []User
	result := Instance.Table("users").Find(&users).Error

	if result != nil {
		log.Fatal(result)
	}
	return users

}

func GetSymbolEntity(symbolId string) (*SymbolEntity, error) {
	var symbolEntity SymbolEntity
	if result := Instance.Table("symbol_entities").Where("symbol_id = ?", symbolId).First(&symbolEntity).Error; result != nil {
		return nil, result
	}
	return &symbolEntity, nil
}

func UpdateMasterAccountBalance(accountLogin, balance int64) error {

	if result := Instance.Table("master_accounts").Model(&MasterAccount{}).Where("account_login = ?", accountLogin).Update("balance", balance); result.Error != nil {
		return result.Error
	}
	return nil
}

func InsertMasterAccountDetails(masterAccount MasterAccount) {
	record := Instance.Table("master_accounts").Model(&MasterAccount{}).Create(&masterAccount)
	if record.Error != nil {
		log.Fatal(record.Error)
	}
}

func SaveUser(user User) {
	err := Instance.Save(&user).Error
	if err != nil {
		log.Fatal(err)
	}
}

func SaveSymbolEntity(symbolentity SymbolEntity) {
	err := Instance.Save(&symbolentity).Error
	if err != nil {
		log.Fatal(err)
	}
}
