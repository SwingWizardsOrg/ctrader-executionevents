package database

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName              string            `gorm:"size:255;not null" form:"firstname"`
	LastName               string            `gorm:"size:255;not null" form:"lastname"`
	Username               string            `gorm:"size:150;not null;unique" form:"username"`
	Email                  string            `gorm:"size:100;not null;unique" form:"email"`
	Password               string            `gorm:"size:100;not null" form:"password"`
	PhoneNumber            string            `gorm:"size:50;default:null;unique" form:"phonenumber,omitempty"`
	Balance                *float64          `gorm:"default:0" form:"balance"`
	PercentageContribution *float64          `gorm:"default:0" form:"contribution,omitempty"`
	FloatingProfit         *float64          `gorm:"default:0" form:"floatingprofit,omitempty"`
	Equity                 *float64          `gorm:"default:0" form:"equity,omitempty"`
	Positions              []RunningPosition `gorm:"foreignkey:UserID" form:"positions,omitempty"`
}

type RunningPosition struct {
	gorm.Model
	UserID          uint
	Volume          *int64
	Price           *float64
	TradeSide       *int32
	SymbolId        *int64
	OpenTime        *int64
	Commission      *int64
	Swap            *int64
	MoneyDigits     *uint32
	PositionRisk    *float32
	PositionsReward *float32
}
type MasterAccount struct {
	gorm.Model
	AccountLogin uint
	Balance      *int64
}

type SymbolEntity struct {
	gorm.Model
	SymbolId       string
	SymbolName     string
	Lot_ZeroOne    *float32
	Lot_ZeroFive   *float32
	Lot_ZeroTen    *float32
	Lot_ZeroTwenty *float32
	Lot_ZeroThirty *float32
}
