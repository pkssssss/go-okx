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
	runType := os.Getenv("OKX_RUN_TYPE")
	if runType == "" {
		runType = "1"
	}
	gridNum := os.Getenv("OKX_GRID_NUM")
	if gridNum == "" {
		gridNum = "10"
	}
	maxPx := os.Getenv("OKX_MAX_PX")
	if maxPx == "" {
		maxPx = "50000"
	}
	minPx := os.Getenv("OKX_MIN_PX")
	if minPx == "" {
		minPx = "40000"
	}

	svc := okx.NewClient().NewTradingBotGridMinInvestmentService().
		InstId(instId).
		AlgoOrdType(algoOrdType).
		RunType(runType).
		GridNum(gridNum).
		MaxPx(maxPx).
		MinPx(minPx)

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s algoOrdType=%s runType=%s gridNum=%s minPx=%s maxPx=%s singleAmt=%s investmentData=%d\n",
		instId, algoOrdType, runType, gridNum, minPx, maxPx, res.SingleAmt, len(res.MinInvestmentData),
	)
}
