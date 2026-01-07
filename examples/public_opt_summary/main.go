package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx"
)

func main() {
	uly := os.Getenv("OKX_ULY")
	if uly == "" {
		uly = "BTC-USD"
	}

	svc := okx.NewClient().NewPublicOptSummaryService().Uly(uly)

	if v := os.Getenv("OKX_OPT_EXP_TIME"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_OPT_EXP_TIME: %v", err)
		}
		svc.ExpTime(n)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("uly=%s count=%d\n", uly, len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first instId=%s markVol=%s delta=%s gamma=%s theta=%s vega=%s ts=%d\n", it.InstId, it.MarkVol, it.Delta, it.Gamma, it.Theta, it.Vega, it.TS)
}
