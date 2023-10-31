package persistence

import (
	"ctrader_events/mappers"
	"log"
)

func InsertSwingSymbol(swingsymbol mappers.SwingSymbol) {
	record := Instance.Create(swingsymbol)

	if record.Error != nil {
		log.Fatal(record.Error)
	}
}

func InsertSwingLightSymbol(lightsymbol mappers.SwingLightSymbol) {
	record := Instance.Create(lightsymbol)
	if record.Error != nil {
		log.Fatal(record.Error)
	}

}

func InsertSwingAsset(swingAsset mappers.SwingAsset) {
	record := Instance.Create(swingAsset)
	if record.Error != nil {
		log.Fatal(record.Error)
	}
}

func GetSwingLightSymbol(symbolId int64) mappers.SwingLightSymbol {
	var lightSymbolSchema mappers.SwingLightSymbol
	lightSymbol := Instance.First(&lightSymbolSchema, "symbol_id = ?", symbolId)

	if (lightSymbol.Error) != nil {
		log.Fatal(lightSymbol.Error)
	}
	return lightSymbolSchema
}

func GetAllSwingLightSymmbol() []mappers.SwingLightSymbol {
	var allLight []mappers.SwingLightSymbol
	result := Instance.Find(&allLight)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return allLight
}

func GetSwingAsset(assetId int64) mappers.SwingAsset {
	var swingAssetSchema mappers.SwingAsset
	lightSymbol := Instance.First(&swingAssetSchema, "asset_id = ?", assetId)

	if (lightSymbol.Error) != nil {
		log.Fatal(lightSymbol.Error)
	}
	return swingAssetSchema
}

func GetAllSwingAssets() []mappers.SwingAsset {
	var allassets []mappers.SwingAsset
	result := Instance.Find(&allassets)

	if result.Error != nil {
		log.Fatal(result.Error)
	}

	return allassets
}
