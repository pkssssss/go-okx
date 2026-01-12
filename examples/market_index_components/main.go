package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	index := os.Getenv("OKX_INDEX")
	if index == "" {
		index = "BTC-USD"
	}

	c := okx.NewClient()

	d, err := c.NewMarketIndexComponentsService().Index(index).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("index=%s last=%s ts=%d components=%d\n", d.Index, d.Last, d.TS, len(d.Components))
	if len(d.Components) == 0 {
		return
	}
	it := d.Components[0]
	fmt.Printf("first exch=%s symbol=%s symPx=%s wgt=%s cnvPx=%s\n", it.Exch, it.Symbol, it.SymPx, it.Wgt, it.CnvPx)
}
