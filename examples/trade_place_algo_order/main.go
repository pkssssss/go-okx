package main

import (
	"context"
	"log"
	"os"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to place algo order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	instId := os.Getenv("OKX_INST_ID")
	tdMode := os.Getenv("OKX_TD_MODE")
	side := os.Getenv("OKX_SIDE")
	ordType := os.Getenv("OKX_ORD_TYPE")
	if instId == "" || tdMode == "" || side == "" || ordType == "" {
		log.Fatal("missing env: OKX_INST_ID / OKX_TD_MODE / OKX_SIDE / OKX_ORD_TYPE")
	}

	sz := os.Getenv("OKX_SZ")
	closeFraction := os.Getenv("OKX_CLOSE_FRACTION")
	if sz == "" && closeFraction == "" {
		log.Fatal("missing env: OKX_SZ or OKX_CLOSE_FRACTION")
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

	svc := c.NewPlaceAlgoOrderService().
		InstId(instId).
		TdMode(tdMode).
		Side(side).
		OrdType(ordType)

	if sz != "" {
		svc.Sz(sz)
	} else {
		svc.CloseFraction(closeFraction)
	}

	if v := os.Getenv("OKX_POS_SIDE"); v != "" {
		svc.PosSide(v)
	}
	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}
	if v := os.Getenv("OKX_TAG"); v != "" {
		svc.Tag(v)
	}
	if v := os.Getenv("OKX_ALGO_CL_ORD_ID"); v != "" {
		svc.AlgoClOrdId(v)
	}

	if v := os.Getenv("OKX_TP_TRIGGER_PX"); v != "" {
		svc.TpTriggerPx(v)
	}
	if v := os.Getenv("OKX_TP_TRIGGER_PX_TYPE"); v != "" {
		svc.TpTriggerPxType(v)
	}
	if v := os.Getenv("OKX_TP_ORD_PX"); v != "" {
		svc.TpOrdPx(v)
	}
	if v := os.Getenv("OKX_TP_ORD_KIND"); v != "" {
		svc.TpOrdKind(v)
	}
	if v := os.Getenv("OKX_SL_TRIGGER_PX"); v != "" {
		svc.SlTriggerPx(v)
	}
	if v := os.Getenv("OKX_SL_TRIGGER_PX_TYPE"); v != "" {
		svc.SlTriggerPxType(v)
	}
	if v := os.Getenv("OKX_SL_ORD_PX"); v != "" {
		svc.SlOrdPx(v)
	}
	if v := os.Getenv("OKX_SL_ORD_KIND"); v != "" {
		svc.SlOrdKind(v)
	}

	if v := os.Getenv("OKX_TRIGGER_PX"); v != "" {
		svc.TriggerPx(v)
	}
	if v := os.Getenv("OKX_TRIGGER_PX_TYPE"); v != "" {
		svc.TriggerPxType(v)
	}
	if v := os.Getenv("OKX_ORDER_PX"); v != "" {
		svc.OrderPx(v)
	}

	if v := os.Getenv("OKX_CALLBACK_RATIO"); v != "" {
		svc.CallbackRatio(v)
	}
	if v := os.Getenv("OKX_CALLBACK_SPREAD"); v != "" {
		svc.CallbackSpread(v)
	}
	if v := os.Getenv("OKX_ACTIVE_PX"); v != "" {
		svc.ActivePx(v)
	}

	if v := os.Getenv("OKX_PX_LIMIT"); v != "" {
		svc.PxLimit(v)
	}
	if v := os.Getenv("OKX_SZ_LIMIT"); v != "" {
		svc.SzLimit(v)
	}
	if v := os.Getenv("OKX_TIME_INTERVAL"); v != "" {
		svc.TimeInterval(v)
	}
	if v := os.Getenv("OKX_PX_SPREAD"); v != "" {
		svc.PxSpread(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_place_algo_order: algoId=%s algoClOrdId=%s sCode=%s sMsg=%s", ack.AlgoId, ack.AlgoClOrdId, ack.SCode, ack.SMsg)
}
