package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	ccy := os.Getenv("OKX_CCY")
	if ccy == "" {
		ccy = "BTC"
	}
	instType := os.Getenv("OKX_INST_TYPE") // SPOT/CONTRACTS
	if instType == "" {
		instType = "SPOT"
	}
	period := os.Getenv("OKX_PERIOD")

	c := okx.NewClient()

	svc := c.NewRubikTakerVolumeService().Ccy(ccy).InstType(instType)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Printf("rubik_taker_volume: empty")
		return
	}
	log.Printf("rubik_taker_volume: ccy=%s instType=%s n=%d ts=%d sellVol=%s buyVol=%s", ccy, instType, len(data), data[0].TS, data[0].SellVol, data[0].BuyVol)
}
