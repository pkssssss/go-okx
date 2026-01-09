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
	if os.Getenv("OKX_CONFIRM_ACCOUNT_MOVE_POSITIONS") != "YES" {
		log.Fatal("refusing to move positions; set OKX_CONFIRM_ACCOUNT_MOVE_POSITIONS=YES to continue")
	}

	fromAcct := os.Getenv("OKX_FROM_ACCT")
	toAcct := os.Getenv("OKX_TO_ACCT")
	clientId := os.Getenv("OKX_CLIENT_ID")
	posId := os.Getenv("OKX_POS_ID")
	side := os.Getenv("OKX_SIDE") // buy/sell
	sz := os.Getenv("OKX_SZ")
	if fromAcct == "" || toAcct == "" || clientId == "" || posId == "" || side == "" || sz == "" {
		log.Fatal("missing env: OKX_FROM_ACCT / OKX_TO_ACCT / OKX_CLIENT_ID / OKX_POS_ID / OKX_SIDE / OKX_SZ")
	}

	tdMode := os.Getenv("OKX_TD_MODE")   // optional
	posSide := os.Getenv("OKX_POS_SIDE") // optional
	ccy := os.Getenv("OKX_CCY")          // optional

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

	leg := okx.AccountMovePositionsLeg{
		From: okx.AccountMovePositionsLegFrom{PosId: posId, Sz: sz, Side: side},
		To:   okx.AccountMovePositionsLegTo{TdMode: tdMode, PosSide: posSide, Ccy: ccy},
	}

	ack, err := c.NewAccountMovePositionsService().
		FromAcct(fromAcct).
		ToAcct(toAcct).
		ClientId(clientId).
		Legs([]okx.AccountMovePositionsLeg{leg}).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("account_move_positions: clientId=%s blockTdId=%s state=%s ts=%d legs=%d", ack.ClientId, ack.BlockTdId, ack.State, ack.TS, len(ack.Legs))
}
