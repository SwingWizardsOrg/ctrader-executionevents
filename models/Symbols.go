package models

type Symbols struct {
	Symbols []Symbol `json:"symbols"`
}

type Symbol struct {
	SymbolId       string  `json:"symbolId"`
	SymbolName     string  `json:"symbolname"`
	Lot_ZeroOne    float32 `json:"lot_0.01"`
	Lot_ZeroFive   float32 `json:"lot_0.05"`
	Lot_ZeroTen    float32 `json:"lot_0.10"`
	Lot_ZeroTwenty float32 `json:"lot_0.20"`
	Lot_ZeroThirty float32 `json:"lot_0.30"`
}
