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

	svc := c.NewRubikOptionOpenInterestVolumeRatioService().Ccy(ccy)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Printf("rubik_option_open_interest_volume_ratio: empty")
		return
	}
	log.Printf("rubik_option_open_interest_volume_ratio: ccy=%s n=%d ts=%d oiRatio=%s volRatio=%s", ccy, len(data), data[0].TS, data[0].OIRatio, data[0].VolRatio)
}
