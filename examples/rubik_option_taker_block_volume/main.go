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
	period := os.Getenv("OKX_PERIOD")

	c := okx.NewClient()

	svc := c.NewRubikOptionTakerBlockVolumeService().Ccy(ccy)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rubik_option_taker_block_volume: ccy=%s ts=%d callBuy=%s callSell=%s putBuy=%s putSell=%s callBlock=%s putBlock=%s", ccy, data.TS, data.CallBuyVol, data.CallSellVol, data.PutBuyVol, data.PutSellVol, data.CallBlockVol, data.PutBlockVol)
}
