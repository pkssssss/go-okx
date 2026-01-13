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

	svc := c.NewRubikOptionOpenInterestVolumeService().Ccy(ccy)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Printf("rubik_option_open_interest_volume: empty")
		return
	}
	log.Printf("rubik_option_open_interest_volume: ccy=%s n=%d ts=%d oi=%s vol=%s", ccy, len(data), data[0].TS, data[0].OI, data[0].Vol)
}
