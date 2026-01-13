package main

import (
	"context"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()

	res, err := c.NewRubikSupportCoinService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("rubik_support_coin: contract=%d option=%d spot=%d", len(res.Contract), len(res.Option), len(res.Spot))
}
