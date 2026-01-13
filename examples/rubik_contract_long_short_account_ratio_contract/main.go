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

	c := okx.NewClient()

	svc := c.NewRubikLongShortAccountRatioContractService().InstId(instId)
	if period != "" {
		svc.Period(period)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(data) == 0 {
		log.Printf("rubik_contract_long_short_account_ratio_contract: empty")
		return
	}
	log.Printf("rubik_contract_long_short_account_ratio_contract: instId=%s n=%d ts=%d ratio=%s", instId, len(data), data[0].TS, data[0].Ratio)
}
