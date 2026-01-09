package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		instId = "BTC-USDT"
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SPOT"
	}

	source := 0
	if v := os.Getenv("OKX_SOURCE"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid OKX_SOURCE: %v", err)
		}
		source = n
	}

	c := okx.NewClient()

	instruments, err := c.NewPublicInstrumentsService().InstType(instType).InstId(instId).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if len(instruments) == 0 {
		log.Fatalf("instrument not found: instType=%s instId=%s", instType, instId)
	}
	if instruments[0].InstIdCode == nil {
		log.Fatalf("missing instIdCode for: instType=%s instId=%s", instType, instId)
	}

	instIdCode := *instruments[0].InstIdCode

	data, err := c.NewMarketBooksSBEService().InstIdCode(instIdCode).Source(source).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	prefix := data
	if len(prefix) > 32 {
		prefix = prefix[:32]
	}

	fmt.Printf("instType=%s instId=%s instIdCode=%d source=%d bytes=%d prefix=%x\n", instType, instId, instIdCode, source, len(data), prefix)
}
