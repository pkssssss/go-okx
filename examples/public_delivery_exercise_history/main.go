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
		instType = "FUTURES"
	}

	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		instFamily = "BTC-USD"
	}

	svc := okx.NewClient().NewPublicDeliveryExerciseHistoryService().InstType(instType).InstFamily(instFamily)

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

	fmt.Printf("instType=%s instFamily=%s count=%d\n", instType, instFamily, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	if len(it.Details) == 0 {
		fmt.Printf("first ts=%d details=0\n", it.TS)
		return
	}
	fmt.Printf("first ts=%d type=%s insId=%s px=%s\n", it.TS, it.Details[0].Type, it.Details[0].InsId, it.Details[0].Px)
}
