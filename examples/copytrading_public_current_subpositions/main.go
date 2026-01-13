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

	svc := okx.NewClient().NewCopyTradingPublicCurrentSubpositionsService().
		InstType(instType).
		UniqueCode(uniqueCode)

	if v := os.Getenv("OKX_AFTER"); v != "" {
		svc.After(v)
	}
	if v := os.Getenv("OKX_BEFORE"); v != "" {
		svc.Before(v)
	}
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		svc.Limit(n)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_current_subpositions: instType=%s uniqueCode=%s items=%d", instType, uniqueCode, len(items))
	for i := 0; i < len(items) && i < 3; i++ {
		it := items[i]
		log.Printf("item[%d]: instId=%s subPosId=%s posSide=%s subPos=%s openTime=%d upl=%s", i, it.InstId, it.SubPosId, it.PosSide, it.SubPos, it.OpenTime, it.Upl)
	}
}
