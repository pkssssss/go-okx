package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	c := okx.NewClient()

	v, err := c.NewMarketPlatform24VolumeService().Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("volUsd=%s volCny=%s ts=%d\n", v.VolUsd, v.VolCny, v.TS)
}
