package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	cfg, err := okx.NewClient().NewCopyTradingPublicConfigService().InstType(instType).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_config: instType=%s maxCopyAmt=%s minCopyAmt=%s maxCopyTotalAmt=%s", instType, cfg.MaxCopyAmt, cfg.MinCopyAmt, cfg.MaxCopyTotalAmt)
}
