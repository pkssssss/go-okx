package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

	tdMode := os.Getenv("OKX_TD_MODE")
	if tdMode == "" {
		tdMode = "cross"
	}

	svc := okx.NewClient().NewPublicPositionTiersService().InstType(instType).TdMode(tdMode)

	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
	}
	if v := os.Getenv("OKX_ULY"); v != "" {
		svc.Uly(v)
	}
	if v := os.Getenv("OKX_INST_ID"); v != "" {
		svc.InstId(v)
	}
	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}
	if v := os.Getenv("OKX_TIER"); v != "" {
		svc.Tier(v)
	}

	if instType != "MARGIN" && os.Getenv("OKX_INST_FAMILY") == "" && os.Getenv("OKX_ULY") == "" {
		svc.InstFamily("BTC-USDT")
	}
	if instType == "MARGIN" && os.Getenv("OKX_INST_ID") == "" && os.Getenv("OKX_CCY") == "" {
		svc.InstId("BTC-USDT")
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("instType=%s tdMode=%s count=%d\n", instType, tdMode, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first tier=%s instFamily=%s maxSz=%s maxLever=%s mmr=%s imr=%s\n", it.Tier, it.InstFamily, it.MaxSz, it.MaxLever, it.MMR, it.IMR)
}
