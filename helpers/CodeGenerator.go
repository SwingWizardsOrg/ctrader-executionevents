package helpers

import (
	"ctrader_events/messages/github.com/Carlosokumu/messages"
	"math/rand"
	"time"
)

func GenerateCode() int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	for {
		// Generate a random 4-digit code
		code := rand.Intn(10000)
		return code
	}
}

func DetermineTradeSide(tradeSide *messages.ProtoOATradeSide) int32 {
	if *tradeSide == messages.ProtoOATradeSide_BUY {
		return 1
	} else {
		return 2
	}
}

func GetStopLossPips(relativeStopLoss *int64, volume *int64) float32 {
	switch *volume {
	case 100000:
		{
			return float32(*relativeStopLoss) / 1000.0
		}
	case 3000000:
		{
			return float32(*relativeStopLoss) / 1000.0
		}
	default:
		panic("unimplemented")
	}
}

func GetTakeProfitPips(relativeTakeProfit *int64, volume *int64) float32 {

	switch *volume {
	case 100000:
		{
			return float32(*relativeTakeProfit) / 1000.0
		}
	case 3000000:
		{
			return float32(*relativeTakeProfit) / 1000.0
		}
	default:
		panic("unimplemented")
	}

}
