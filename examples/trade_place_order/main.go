package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to place order; set OKX_CONFIRM=YES to continue")
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
	sz := os.Getenv("OKX_SZ")
	if instId == "" || tdMode == "" || side == "" || ordType == "" || sz == "" {
		log.Fatal("missing env: OKX_INST_ID / OKX_TD_MODE / OKX_SIDE / OKX_ORD_TYPE / OKX_SZ")
	}

	px := os.Getenv("OKX_PX")
	pxUsd := os.Getenv("OKX_PX_USD")
	pxVol := os.Getenv("OKX_PX_VOL")
	pxCount := 0
	if px != "" {
		pxCount++
	}
	if pxUsd != "" {
		pxCount++
	}
	if pxVol != "" {
		pxCount++
	}
	if pxCount > 1 {
		log.Fatal("too many price envs: require at most one of OKX_PX / OKX_PX_USD / OKX_PX_VOL")
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

	syncCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := c.SyncTime(syncCtx); err != nil {
		log.Fatal(err)
	}

	clOrdId := os.Getenv("OKX_CL_ORD_ID")
	if clOrdId == "" {
		// 默认生成一个可复用的 clOrdId（便于幂等重试与排查）。如果你需要人工重试，请复用日志中打印的 clOrdId。
		clOrdId = "gookx" + strconv.FormatInt(time.Now().UnixNano(), 36)
		log.Printf("generated OKX_CL_ORD_ID=%s (reuse it for idempotent retry)", clOrdId)
	}

	svc := c.NewPlaceOrderService().
		InstId(instId).
		TdMode(tdMode).
		ClOrdId(clOrdId).
		Side(side).
		OrdType(ordType).
		Sz(sz)

	if v := os.Getenv("OKX_CCY"); v != "" {
		svc.Ccy(v)
	}
	if v := os.Getenv("OKX_TAG"); v != "" {
		svc.Tag(v)
	}
	if v := os.Getenv("OKX_POS_SIDE"); v != "" {
		svc.PosSide(v)
	}

	if px != "" {
		svc.Px(px)
	}
	if pxUsd != "" {
		svc.PxUsd(pxUsd)
	}
	if pxVol != "" {
		svc.PxVol(pxVol)
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
	if v := os.Getenv("OKX_BAN_AMEND"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("invalid env OKX_BAN_AMEND=%q: %v", v, err)
		}
		svc.BanAmend(b)
	}
	if v := os.Getenv("OKX_PX_AMEND_TYPE"); v != "" {
		svc.PxAmendType(v)
	}
	if v := os.Getenv("OKX_TRADE_QUOTE_CCY"); v != "" {
		svc.TradeQuoteCcy(v)
	}
	if v := os.Getenv("OKX_STP_MODE"); v != "" {
		svc.StpMode(v)
	}
	if v := os.Getenv("OKX_EXP_TIME"); v != "" {
		svc.ExpTime(v)
	}

	orderCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ack, err := svc.Do(orderCtx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_place_order: clOrdId=%s ordId=%s sCode=%s sMsg=%s", ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg)
}
