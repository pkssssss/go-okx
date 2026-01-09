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

	instId := os.Getenv("OKX_INST_ID")
	tdMode := os.Getenv("OKX_TD_MODE")
	side := os.Getenv("OKX_SIDE")
	ordType := os.Getenv("OKX_ORD_TYPE")
	sz := os.Getenv("OKX_SZ")
	if instId == "" || tdMode == "" || side == "" || ordType == "" || sz == "" {
		log.Fatal("missing env: OKX_INST_ID / OKX_TD_MODE / OKX_SIDE / OKX_ORD_TYPE / OKX_SZ")
	}

	px := os.Getenv("OKX_PX") // required for limit/post_only/fok/ioc

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

	svc := c.NewOrderPrecheckService().
		InstId(instId).
		TdMode(tdMode).
		Side(side).
		OrdType(ordType).
		Sz(sz)
	if px != "" {
		svc.Px(px)
	}

	if v := os.Getenv("OKX_CL_ORD_ID"); v != "" {
		svc.ClOrdId(v)
	}
	if v := os.Getenv("OKX_POS_SIDE"); v != "" {
		svc.PosSide(v)
	}
	if v := os.Getenv("OKX_REDUCE_ONLY"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("invalid env OKX_REDUCE_ONLY=%q: %v", v, err)
		}
		svc.ReduceOnly(b)
	}
	if v := os.Getenv("OKX_TGT_CCY"); v != "" {
		svc.TgtCcy(v)
	}

	items, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_order_precheck: instId=%s tdMode=%s side=%s ordType=%s sz=%s px=%s items=%d", instId, tdMode, side, ordType, sz, px, len(items))
	for i := 0; i < len(items) && i < 5; i++ {
		x := items[i]
		log.Printf("item[%d]: adjEq=%s adjEqChg=%s imr=%s imrChg=%s mmr=%s mmrChg=%s mgnRatio=%s", i, x.AdjEq, x.AdjEqChg, x.Imr, x.ImrChg, x.Mmr, x.MmrChg, x.MgnRatio)
	}
}
