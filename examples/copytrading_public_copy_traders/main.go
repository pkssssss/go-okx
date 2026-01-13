package main

import (
	"context"
	"log"
	"os"
	"strconv"

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

	svc := okx.NewClient().NewCopyTradingPublicCopyTradersService().
		InstType(instType).
		UniqueCode(uniqueCode)

	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	res, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_copy_traders: instType=%s uniqueCode=%s copyTotalPnl=%s ccy=%s copyTraders=%d", instType, uniqueCode, res.CopyTotalPnl, res.Ccy, len(res.CopyTraders))
	for i := 0; i < len(res.CopyTraders) && i < 3; i++ {
		it := res.CopyTraders[i]
		log.Printf("trader[%d]: beginCopyTime=%d nickName=%s pnl=%s", i, it.BeginCopyTime, it.NickName, it.Pnl)
	}
}
