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
		instId = "BTC-USDT"
	}
	algoOrdType := os.Getenv("OKX_ALGO_ORD_TYPE")
	if algoOrdType == "" {
		algoOrdType = "grid"
	}

	svc := okx.NewClient().NewTradingBotGridAIParamService().
		AlgoOrdType(algoOrdType).
		InstId(instId)

	if v := os.Getenv("OKX_DIRECTION"); v != "" {
		svc.Direction(v)
	}
	if v := os.Getenv("OKX_DURATION"); v != "" {
		svc.Duration(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s algoOrdType=%s count=%d\n", instId, algoOrdType, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	fmt.Printf("first duration=%s gridNum=%s minPx=%s maxPx=%s annualizedRate=%s minInvestment=%s ccy=%s\n",
		it.Duration, it.GridNum, it.MinPx, it.MaxPx, it.AnnualizedRate, it.MinInvestment, it.Ccy,
	)
}
