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
		instType = "SPOT"
	}

	svc := okx.NewClient().NewPublicInstrumentsService().InstType(instType)

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
	fmt.Printf("first instId=%s base=%s quote=%s settle=%s tickSz=%s lotSz=%s minSz=%s state=%s\n", it.InstId, it.BaseCcy, it.QuoteCcy, it.SettleCcy, it.TickSz, it.LotSz, it.MinSz, it.State)
}
