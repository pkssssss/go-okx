package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()

	r, err := c.NewMarketExchangeRateService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("usdCny=%s\n", r.UsdCny)
}
