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
	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instId == "" && instFamily == "" {
		instFamily = "BTC-USD"
	}

	svc := okx.NewClient().NewPublicOptionTradesService()
	if instId != "" {
		svc.InstId(instId)
	} else {
		svc.InstFamily(instFamily)
	}
	if v := os.Getenv("OKX_OPT_TYPE"); v != "" {
		svc.OptType(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count=%d\n", len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first instId=%s tradeId=%s px=%s sz=%s side=%s optType=%s ts=%d\n", it.InstId, it.TradeId, it.Px, it.Sz, it.Side, it.OptType, it.TS)
}
