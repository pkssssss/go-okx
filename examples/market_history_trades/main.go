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
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	typ := os.Getenv("OKX_TYPE") // 1: tradeId 分页；2: 时间戳分页
	after := os.Getenv("OKX_AFTER")
	before := os.Getenv("OKX_BEFORE")

	limit := 20
	if v := os.Getenv("OKX_TRADES_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_TRADES_LIMIT: %v", err)
		}
		limit = n
	}

	c := okx.NewClient()

	svc := c.NewMarketHistoryTradesService().InstId(instId).Limit(limit)
	if typ != "" {
		svc.Type(typ)
	}
	if after != "" {
		svc.After(after)
	}
	if before != "" {
		svc.Before(before)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instId=%s type=%s count=%d\n", instId, typ, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first tradeId=%s px=%s sz=%s side=%s ts=%d\n", it.TradeId, it.Px, it.Sz, it.Side, it.TS)
}
