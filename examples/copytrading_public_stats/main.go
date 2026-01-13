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

	stats, err := okx.NewClient().NewCopyTradingPublicStatsService().
		InstType(instType).
		UniqueCode(uniqueCode).
		LastDays(lastDays).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_stats: instType=%s uniqueCode=%s lastDays=%s winRatio=%s profitDays=%s lossDays=%s curPnl=%s investAmt=%s", instType, uniqueCode, lastDays, stats.WinRatio, stats.ProfitDays, stats.LossDays, stats.CurCopyTraderPnl, stats.InvestAmt)
}
