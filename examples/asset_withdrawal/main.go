package main

import (
	"context"
	"encoding/json"
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
	if os.Getenv("OKX_CONFIRM_WITHDRAWAL") != "YES" {
		log.Fatal("refusing to withdraw; set OKX_CONFIRM_WITHDRAWAL=YES to continue")
	}

	demo := os.Getenv("OKX_DEMO") == "1"

	ccy := os.Getenv("OKX_CCY")
	amt := os.Getenv("OKX_WITHDRAWAL_AMT")
	dest := os.Getenv("OKX_WITHDRAWAL_DEST")
	toAddr := os.Getenv("OKX_WITHDRAWAL_TO_ADDR")
	if ccy == "" || amt == "" || dest == "" || toAddr == "" {
		log.Fatal("missing env: OKX_CCY / OKX_WITHDRAWAL_AMT / OKX_WITHDRAWAL_DEST / OKX_WITHDRAWAL_TO_ADDR")
	}

	toAddrType := os.Getenv("OKX_WITHDRAWAL_TO_ADDR_TYPE")
	chain := os.Getenv("OKX_WITHDRAWAL_CHAIN")
	areaCode := os.Getenv("OKX_WITHDRAWAL_AREA_CODE")
	clientId := os.Getenv("OKX_WITHDRAWAL_CLIENT_ID")

	var rcvrInfo json.RawMessage
	if v := os.Getenv("OKX_WITHDRAWAL_RCVR_INFO_JSON"); v != "" {
		rcvrInfo = json.RawMessage(v)
		if !json.Valid(rcvrInfo) {
			log.Fatal("invalid OKX_WITHDRAWAL_RCVR_INFO_JSON")
		}
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

	svc := c.NewAssetWithdrawalService().
		Ccy(ccy).
		Amt(amt).
		Dest(dest).
		ToAddr(toAddr)
	if toAddrType != "" {
		svc.ToAddrType(toAddrType)
	}
	if chain != "" {
		svc.Chain(chain)
	}
	if areaCode != "" {
		svc.AreaCode(areaCode)
	}
	if len(rcvrInfo) > 0 {
		svc.RcvrInfoJSON(rcvrInfo)
	}
	if clientId != "" {
		svc.ClientId(clientId)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("withdrawal: wdId=%s ccy=%s chain=%s amt=%s clientId=%s", ack.WdId, ack.Ccy, ack.Chain, ack.Amt, ack.ClientId)
}
