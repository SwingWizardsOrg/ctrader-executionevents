package models

type ResourceId struct {
	ResourceName string `json:"resourcename"`
}

type AccountModelUseCase struct {
	Positions    []MarketOrderModel
	Equity       float64
	Balance      float64
	Symbols      []SymbolModel
	DepositAsset int64
}

type SpotEventUseCase struct {
	Bid      uint64
	Ask      uint64
	SymbolId int64
}
