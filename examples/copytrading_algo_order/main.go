package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to place copytrading algo order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	subPosId := os.Getenv("OKX_SUB_POS_ID")
	if subPosId == "" {
		log.Fatal("missing env: OKX_SUB_POS_ID")
	}

	tpTriggerPx := os.Getenv("OKX_TP_TRIGGER_PX")
	slTriggerPx := os.Getenv("OKX_SL_TRIGGER_PX")
	if tpTriggerPx == "" && slTriggerPx == "" {
		log.Fatal("missing env: OKX_TP_TRIGGER_PX or OKX_SL_TRIGGER_PX")
	}

	instType := os.Getenv("OKX_INST_TYPE")
	if instType == "" {
		instType = "SWAP"
	}

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

	svc := c.NewCopyTradingAlgoOrderService().
		InstType(instType).
		SubPosId(subPosId)

	if tpTriggerPx != "" {
		svc.TpTriggerPx(tpTriggerPx)
	}
	if slTriggerPx != "" {
		svc.SlTriggerPx(slTriggerPx)
	}

	if v := os.Getenv("OKX_TP_ORD_PX"); v != "" {
		svc.TpOrdPx(v)
	}
	if v := os.Getenv("OKX_SL_ORD_PX"); v != "" {
		svc.SlOrdPx(v)
	}
	if v := os.Getenv("OKX_TP_TRIGGER_PX_TYPE"); v != "" {
		svc.TpTriggerPxType(v)
	}
	if v := os.Getenv("OKX_SL_TRIGGER_PX_TYPE"); v != "" {
		svc.SlTriggerPxType(v)
	}
	if v := os.Getenv("OKX_TAG"); v != "" {
		svc.Tag(v)
	}
	if v := os.Getenv("OKX_SUB_POS_TYPE"); v != "" {
		svc.SubPosType(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("copytrading_algo_order: subPosId=%s tag=%s", ack.SubPosId, ack.Tag)
}
