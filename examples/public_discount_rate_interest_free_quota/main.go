package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	svc := okx.NewClient().NewPublicDiscountRateInterestFreeQuotaService()

	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count=%d\n", len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("first ccy=%s amt=%s colRes=%s discountLv=%s details=%d\n", it.Ccy, it.Amt, it.ColRes, it.DiscountLv, len(it.Details))
}
