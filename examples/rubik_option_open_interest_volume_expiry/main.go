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

	svc := c.NewRubikOptionOpenInterestVolumeExpiryService().Ccy(ccy)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Printf("rubik_option_open_interest_volume_expiry: empty")
		return
	}
	log.Printf("rubik_option_open_interest_volume_expiry: ccy=%s n=%d ts=%d expTime=%s callOI=%s putOI=%s", ccy, len(data), data[0].TS, data[0].ExpTime, data[0].CallOI, data[0].PutOI)
}
