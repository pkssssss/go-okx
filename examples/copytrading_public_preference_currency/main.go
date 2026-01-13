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

	items, err := okx.NewClient().NewCopyTradingPublicPreferenceCurrencyService().
		InstType(instType).
		UniqueCode(uniqueCode).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_public_preference_currency: instType=%s uniqueCode=%s items=%d", instType, uniqueCode, len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		it := items[i]
		log.Printf("item[%d]: ccy=%s ratio=%s", i, it.Ccy, it.Ratio)
	}
}
