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
		instType = "OPTION"
	}

	svc := okx.NewClient().NewPublicInstrumentTickBandsService().InstType(instType)
	if v := os.Getenv("OKX_INST_FAMILY"); v != "" {
		svc.InstFamily(v)
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
	fmt.Printf("first instFamily=%s tickBands=%d\n", it.InstFamily, len(it.TickBands))
	if len(it.TickBands) > 0 {
		b := it.TickBands[0]
		fmt.Printf("first band: minPx=%s maxPx=%s tickSz=%s\n", b.MinPx, b.MaxPx, b.TickSz)
	}
}
