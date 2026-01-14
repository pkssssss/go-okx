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
	runType := os.Getenv("OKX_RUN_TYPE")
	if runType == "" {
		runType = "1"
	}
	algoOrdType := os.Getenv("OKX_ALGO_ORD_TYPE")
	if algoOrdType == "" {
		algoOrdType = "grid"
	}
	maxPx := os.Getenv("OKX_MAX_PX")
	if maxPx == "" {
		maxPx = "50000"
	}
	minPx := os.Getenv("OKX_MIN_PX")
	if minPx == "" {
		minPx = "40000"
	}

	svc := okx.NewClient().NewTradingBotGridGridQuantityService().
		InstId(instId).
		RunType(runType).
		AlgoOrdType(algoOrdType).
		MaxPx(maxPx).
		MinPx(minPx)

	if v := os.Getenv("OKX_LEVER"); v != "" {
		svc.Lever(v)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s algoOrdType=%s runType=%s minPx=%s maxPx=%s maxGridQty=%s\n", instId, algoOrdType, runType, minPx, maxPx, res.MaxGridQty)
}
