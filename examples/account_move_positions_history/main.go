package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	blockTdId := os.Getenv("OKX_BLOCK_TD_ID")
	clientId := os.Getenv("OKX_CLIENT_ID")
	beginTs := os.Getenv("OKX_BEGIN_TS")
	endTs := os.Getenv("OKX_END_TS")
	limit := os.Getenv("OKX_LIMIT")
	state := os.Getenv("OKX_STATE") // filled/pending

	demo := os.Getenv("OKX_DEMO") == "1"

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

	svc := c.NewAccountMovePositionsHistoryService()
	if blockTdId != "" {
		svc.BlockTdId(blockTdId)
	}
	if clientId != "" {
		svc.ClientId(clientId)
	}
	if beginTs != "" {
		svc.BeginTs(beginTs)
	}
	if endTs != "" {
		svc.EndTs(endTs)
	}
	if limit != "" {
		svc.Limit(limit)
	}
	if state != "" {
		svc.State(state)
	}

	data, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_move_positions_history: items=%d", len(data))
	for _, it := range data {
		log.Printf("clientId=%s blockTdId=%s state=%s ts=%d from=%s to=%s legs=%d", it.ClientId, it.BlockTdId, it.State, it.TS, it.FromAcct, it.ToAcct, len(it.Legs))
	}
}
