package mappers

type SwingSymbol struct {
	Digits   int32
	SymbolId int64
}

type SwingLightSymbol struct {
	SymbolId     int64
	SymbolName   string
	BaseAssetId  int64
	QuoteAssetId int64
}

type SwingAsset struct {
	AssetId     int64
	Name        string
	DisplayName string
}
