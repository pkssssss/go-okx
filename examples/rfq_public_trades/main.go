package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	beginId := os.Getenv("OKX_BEGIN_ID") // 可选
	endId := os.Getenv("OKX_END_ID")     // 可选

	var limit *int
	if v := os.Getenv("OKX_LIMIT"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_LIMIT: %v", err)
		}
		limit = &n
	}

	c := okx.NewClient()

	svc := c.NewRFQPublicTradesService()
	if beginId != "" {
		svc = svc.BeginId(beginId)
	}
	if endId != "" {
		svc = svc.EndId(endId)
	}
	if limit != nil {
		svc = svc.Limit(*limit)
	}

	trades, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rfq public trades count=%d", len(trades))
	if len(trades) > 0 {
		log.Printf("first: blockTdId=%s strategy=%s cTime=%d legs=%d", trades[0].BlockTdId, trades[0].Strategy, trades[0].CTime, len(trades[0].Legs))
		if len(trades[0].Legs) > 0 {
			leg := trades[0].Legs[0]
			log.Printf("first leg: instId=%s side=%s sz=%s px=%s tradeId=%s", leg.InstId, leg.Side, leg.Sz, leg.Px, leg.TradeId)
		}
	}
}
