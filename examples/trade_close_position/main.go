package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to close positions; set OKX_CONFIRM=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	instId := os.Getenv("OKX_INST_ID")
	mgnMode := os.Getenv("OKX_MGN_MODE")
	if instId == "" || mgnMode == "" {
		log.Fatal("missing env: OKX_INST_ID / OKX_MGN_MODE")
	}

	posSide := os.Getenv("OKX_POS_SIDE")
	ccy := os.Getenv("OKX_CCY")
	clOrdId := os.Getenv("OKX_CL_ORD_ID")
	tag := os.Getenv("OKX_TAG")

	var autoCxl *bool
	if v := os.Getenv("OKX_AUTO_CXL"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("invalid env OKX_AUTO_CXL=%q: %v", v, err)
		}
		autoCxl = &b
	}

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

	svc := c.NewClosePositionsService().InstId(instId).MgnMode(mgnMode)
	if posSide != "" {
		svc.PosSide(posSide)
	}
	if ccy != "" {
		svc.Ccy(ccy)
	}
	if autoCxl != nil {
		svc.AutoCxl(*autoCxl)
	}
	if clOrdId != "" {
		svc.ClOrdId(clOrdId)
	}
	if tag != "" {
		svc.Tag(tag)
	}

	acks, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_close_position: instId=%s mgnMode=%s acks=%d", instId, mgnMode, len(acks))
	for i := 0; i < len(acks) && i < 5; i++ {
		ack := acks[i]
		log.Printf("ack[%d]: instId=%s posSide=%s clOrdId=%s tag=%s", i, ack.InstId, ack.PosSide, ack.ClOrdId, ack.Tag)
	}
}
