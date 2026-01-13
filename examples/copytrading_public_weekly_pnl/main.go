package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	uniqueCode := os.Getenv("OKX_UNIQUE_CODE")
	if uniqueCode == "" {
		log.Fatal("missing env: OKX_UNIQUE_CODE")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	items, err := okx.NewClient().NewCopyTradingPublicWeeklyPnlService().
		InstType(instType).
		UniqueCode(uniqueCode).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_weekly_pnl: instType=%s uniqueCode=%s items=%d", instType, uniqueCode, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	log.Printf("first: beginTs=%d pnl=%s pnlRatio=%s", it.BeginTs, it.Pnl, it.PnlRatio)
}
