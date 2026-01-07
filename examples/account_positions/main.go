package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	instType := os.Getenv("OKX_INST_TYPE")
	instId := os.Getenv("OKX_INST_ID")

	c := okx.NewClient(
		okx.WithCredentials(okx.Credentials{
			APIKey:     apiKey,
			SecretKey:  secretKey,
			Passphrase: passphrase,
		}),
		okx.WithDemoTrading(demo),
	)

	if _, err := c.SyncTime(context.Background()); err != nil {
		log.Fatal(err)
	}

	svc := c.NewAccountPositionsService()
	if instType != "" {
		svc.InstType(instType)
	}
	if instId != "" {
		svc.InstId(instId)
	}

	positions, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("positions=%d", len(positions))
	for _, p := range positions {
		log.Printf("instType=%s instId=%s posSide=%s pos=%s avgPx=%s upl=%s", p.InstType, p.InstId, p.PosSide, p.Pos, p.AvgPx, p.Upl)
	}
}
