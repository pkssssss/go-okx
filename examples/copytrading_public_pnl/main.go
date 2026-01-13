package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	uniqueCode := os.Getenv("OKX_UNIQUE_CODE")
	lastDays := os.Getenv("OKX_LAST_DAYS")
	if uniqueCode == "" || lastDays == "" {
		log.Fatal("missing env: OKX_UNIQUE_CODE / OKX_LAST_DAYS")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	items, err := okx.NewClient().NewCopyTradingPublicPnlService().
		InstType(instType).
		UniqueCode(uniqueCode).
		LastDays(lastDays).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_pnl: instType=%s uniqueCode=%s lastDays=%s items=%d", instType, uniqueCode, lastDays, len(items))
	if len(items) == 0 {
		return
	}
	it := items[0]
	log.Printf("first: beginTs=%d pnl=%s pnlRatio=%s", it.BeginTs, it.Pnl, it.PnlRatio)
}
