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
	expTime := os.Getenv("OKX_EXP_TIME")
	if expTime == "" {
		expTime = "20210901"
	}
	period := os.Getenv("OKX_PERIOD")

	c := okx.NewClient()

	svc := c.NewRubikOptionOpenInterestVolumeStrikeService().Ccy(ccy).ExpTime(expTime)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Printf("rubik_option_open_interest_volume_strike: empty")
		return
	}
	log.Printf("rubik_option_open_interest_volume_strike: ccy=%s expTime=%s n=%d ts=%d strike=%s callOI=%s putOI=%s", ccy, expTime, len(data), data[0].TS, data[0].Strike, data[0].CallOI, data[0].PutOI)
}
