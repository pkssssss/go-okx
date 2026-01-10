package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	svc := okx.NewClient().NewPublicInsuranceFundService().InstType(instType)

	if v := os.Getenv("OKX_TYPE"); v != "" {
		svc.Type(v)
	}
	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
	}
	if v := os.Getenv("OKX_ULY"); v != "" {
		svc.Uly(v)
	}
	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}
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

	if instType != "MARGIN" && os.Getenv("OKX_INST_FAMILY") == "" && os.Getenv("OKX_ULY") == "" {
		svc.Uly("BTC-USD")
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
	fmt.Printf("first instFamily=%s total=%s details=%d\n", it.InstFamily, it.Total, len(it.Details))
	if len(it.Details) > 0 {
		d := it.Details[0]
		fmt.Printf("first detail: ccy=%s balance=%s amt=%s type=%s ts=%d\n", d.Ccy, d.Balance, d.Amt, d.Type, d.TS)
	}
}
