package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()

	res, err := c.SyncTime(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("server=%s roundTrip=%s offset=%s\n", res.ServerTime.Format("2006-01-02T15:04:05.000Z"), res.RoundTrip, res.Offset)
}
