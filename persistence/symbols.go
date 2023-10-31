package persistence

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"ctraderapi/messages/github.com/Carlosokumu/messages"

	"gorm.io/gorm"

	badger "github.com/dgraph-io/badger/v3"
)

func CreateBadgerConnection() *badger.DB {
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func InsertSymbolData(Symbol messages.ProtoOASymbol) error {
	db, err := badger.Open(badger.DefaultOptions("/swingwizards/symbols"))
	if err != nil {
		log.Fatal(err)
	}
	errupdate := db.Update(func(txn *badger.Txn) error {
		json, _ := json.Marshal(&Symbol)
		err := txn.Set([]byte(strconv.Itoa(int(*Symbol.SymbolId))), []byte(json))
		if err != nil {
			return err
		}
		return nil
	})
	if errupdate != nil {
		fmt.Println("Insererror:", err)
		return err
	}
	defer db.Close()
	return nil
}

func InsertLightSymbol(symbolId int64, protoLightSymbol messages.ProtoOALightSymbol) error {
	db, err := badger.Open(badger.DefaultOptions("/swingwizards/lightsymbols"))
	if err != nil {
		log.Fatal(err)
	}
	errupdate := db.Update(func(txn *badger.Txn) error {
		json, _ := json.Marshal(&protoLightSymbol)
		err := txn.Set([]byte(strconv.Itoa(int(symbolId))), []byte(json))
		if err != nil {
			return err
		}
		return nil
	})
	if errupdate != nil {
		fmt.Println("Insererror:", err)
		return err
	}
	defer db.Close()
	return nil
}

func ReadSymbolData(SymbolId int64) (*messages.ProtoOASymbol, error) {
	var symbolEntity messages.ProtoOASymbol
	db, err := badger.Open(badger.DefaultOptions("/swingwizards/symbols"))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	readerr := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(int(SymbolId))))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {

			valuehere := json.Unmarshal(val, &symbolEntity)
			if valuehere != nil {
				fmt.Println("unmarsherr:", valuehere)
			}
			fmt.Println("itemhere:", *symbolEntity.Digits)
			return nil
		})

		return nil
	})
	if readerr != nil {
		return nil, readerr
	}
	defer db.Close()
	return &symbolEntity, nil
}

func ReadLightSymbolData(SymbolId int64) (*messages.ProtoOALightSymbol, error) {
	var symbolEntity messages.ProtoOALightSymbol
	db, err := badger.Open(badger.DefaultOptions("/swingwizards/lightsymbols"))
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	readerr := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strconv.Itoa(int(SymbolId))))
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {

			valuehere := json.Unmarshal(val, &symbolEntity)
			if valuehere != nil {
				fmt.Println("unmarsherr:", valuehere)
			}
			fmt.Println("LightSymbol:", symbolEntity)
			return nil
		})

		return nil
	})
	if readerr != nil {
		return nil, readerr
	}
	defer db.Close()
	return &symbolEntity, nil
}

type SymbolData struct {
	Name string `json:"name"`
}

type SymbolModel struct {
	gorm.Model
	ConversionSymbols []SymbolModel `gorm:"foreignkey:ParentID"`
	ID                int64
	ParentID          int64
}
