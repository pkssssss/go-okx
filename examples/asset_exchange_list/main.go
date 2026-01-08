package main

import (
	"context"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()

	exchanges, err := c.NewAssetExchangeListService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("exchanges=%d", len(exchanges))
	for i := 0; i < len(exchanges) && i < 5; i++ {
		x := exchanges[i]
		log.Printf("exchange[%d]: name=%s id=%s", i, x.ExchName, x.ExchId)
	}
}
