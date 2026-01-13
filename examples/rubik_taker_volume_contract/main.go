package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT-SWAP"
	}
	period := os.Getenv("OKX_PERIOD")
	unit := os.Getenv("OKX_UNIT") // 0/1/2

	c := okx.NewClient()

	svc := c.NewRubikTakerVolumeContractService().InstId(instId)
	if period != "" {
		svc.Period(period)
	}
	if unit != "" {
		svc.Unit(unit)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Printf("rubik_taker_volume_contract: empty")
		return
	}
	log.Printf("rubik_taker_volume_contract: instId=%s n=%d ts=%d sellVol=%s buyVol=%s", instId, len(data), data[0].TS, data[0].SellVol, data[0].BuyVol)
}
