package helpers

import (
	"ctraderapi/mappers"
	"ctraderapi/messages/github.com/Carlosokumu/messages"
	"ctraderapi/models"
	"math"
)

// Calculate the Tick Size
func GetTickSize(protoOAsymbol *messages.ProtoOASymbol) float64 {
	tickSize := 1 / math.Pow(10, float64(*protoOAsymbol.Digits))
	return tickSize
}

// Calculating  a Symbol's PipSize
func GetPipSize(protoOAsymbol messages.ProtoOASymbol) float64 {
	pipSize := 1 / math.Pow(10, float64(*protoOAsymbol.PipPosition))
	return pipSize
}

func GetPipValue(protoOAsymbol messages.ProtoOASymbol, tickValue float64) float64 {
	pipSize := GetPipSize(protoOAsymbol)
	tickSize := GetTickSize(&protoOAsymbol)
	pipValue := tickValue * (pipSize / tickSize)
	return pipValue
}

// Get TickValue
func GetTickValue(protoOAsymbol *messages.ProtoOASymbol, symbolQuoteAsset mappers.SwingAsset, accountDepositAssetID int64, tickValues []models.Tuple) float64 {
	var tickValue float64

	symbolTickSize := GetTickSize(protoOAsymbol)

	if symbolQuoteAsset.AssetId == accountDepositAssetID {
		tickValue = symbolTickSize
	} else {
		tickValue = symbolTickSize * Convert(symbolQuoteAsset, tickValues)
	}
	return tickValue
}

func Convert(symbolQuoteAsset mappers.SwingAsset, conversionAssets []models.Tuple) float64 {
	if conversionAssets == nil {
		panic("conversionAssets cannot be nil")
	}

	if len(conversionAssets) == 0 {
		panic("conversionAssets cannot be empty")
	}

	for _, asset := range conversionAssets {
		if asset.Bid == 0 {
			panic("some conversionAssets have a price of zero")
		}
	}

	result := 1.0

	var asset mappers.SwingAsset = symbolQuoteAsset

	for _, conversion := range conversionAssets {
		if asset.AssetId == conversion.BaseAsset.AssetId {
			result *= conversion.Bid
			asset = conversion.QuoteAsset
		} else {
			result /= conversion.Bid
			asset = conversion.BaseAsset
		}
	}

	return result
}

func GetPriceRelative(symbol *messages.ProtoOASymbol, relative float64) float64 {
	value := roundFloat(relative/100000, uint(*symbol.Digits))
	return value
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func GetPipsFromPrice(symbol *messages.ProtoOASymbol, price float64) float64 {
	pow := math.Pow(10, float64(*symbol.PipPosition))
	result := price * pow
	digits := float64(*symbol.Digits - *symbol.PipPosition)
	return math.Round(result*math.Pow(10, -digits)) * math.Pow(10, digits)
}

func FromMonetary(monetary float64) float64 {
	return monetary / 100.0
}
