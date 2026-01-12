package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instFamily := os.Getenv("OKX_INST_FAMILY")
	if instFamily == "" {
		instFamily = "BTC-USD"
	}

	c := okx.NewClient()

	items, err := c.NewMarketOptionInstrumentFamilyTradesService().InstFamily(instFamily).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instFamily=%s groups=%d\n", instFamily, len(items))
	if len(items) == 0 || len(items[0].TradeInfo) == 0 {
		return
	}

	g := items[0]
	it := g.TradeInfo[0]
	fmt.Printf("first optType=%s vol24h=%s instId=%s px=%s sz=%s side=%s tradeId=%s ts=%d\n",
		g.OptType, g.Vol24h, it.InstId, it.Px, it.Sz, it.Side, it.TradeId, it.TS)
}
