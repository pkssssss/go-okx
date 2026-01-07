package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	svc := okx.NewClient().NewPublicMarkPriceService().InstType(instType)

	if v := os.Getenv("OKX_ULY"); v != "" {
		svc.Uly(v)
	}
	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
	}
	if v := os.Getenv("OKX_INST_ID"); v != "" {
		svc.InstId(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instType=%s count=%d\n", instType, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first instId=%s markPx=%s ts=%d\n", it.InstId, it.MarkPx, it.TS)
}
