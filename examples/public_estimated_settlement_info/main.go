package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "XRP-USDT-250307"
	}

	items, err := okx.NewClient().NewPublicEstimatedSettlementInfoService().InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s count=%d\n", instId, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first estSettlePx=%s nextSettleTime=%d ts=%d\n", it.EstSettlePx, it.NextSettleTime, it.TS)
}
