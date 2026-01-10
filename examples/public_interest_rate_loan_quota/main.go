package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	items, err := okx.NewClient().NewPublicInterestRateLoanQuotaService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("count=%d\n", len(items))
	if len(items) == 0 {
		return
	}

	it := items[0]
	fmt.Printf("basic=%d vip=%d regular=%d configCcyList=%d config=%d\n", len(it.Basic), len(it.VIP), len(it.Regular), len(it.ConfigCcyList), len(it.Config))
}
