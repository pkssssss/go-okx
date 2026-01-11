package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	baseCcy := os.Getenv("OKX_SPRD_BASE_CCY") // 可选，例如 BTC
	instId := os.Getenv("OKX_SPRD_INST_ID")   // 可选，例如 BTC-USDT
	if instId == "" {
		instId = "BTC-USDT"
	}
	sprdId := os.Getenv("OKX_SPRD_ID")   // 可选，例如 BTC-USDT_BTC-USDT-SWAP
	state := os.Getenv("OKX_SPRD_STATE") // 可选：live/suspend/expired

	c := okx.NewClient()

	svc := c.NewSprdSpreadsService().InstId(instId)
	if baseCcy != "" {
		svc = svc.BaseCcy(baseCcy)
	}
	if sprdId != "" {
		svc = svc.SprdId(sprdId)
	}
	if state != "" {
		svc = svc.State(state)
	}

	spreads, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("sprd spreads count=%d", len(spreads))
	if len(spreads) > 0 {
		s := spreads[0]
		log.Printf("first: sprdId=%s state=%s baseCcy=%s quoteCcy=%s tickSz=%s minSz=%s lotSz=%s legs=%d", s.SprdId, s.State, s.BaseCcy, s.QuoteCcy, s.TickSz, s.MinSz, s.LotSz, len(s.Legs))
	}
}
