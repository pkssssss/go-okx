package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/pkssssss/go-okx/v5"
)

func main() {
	if os.Getenv("OKX_CONFIRM") != "YES" {
		log.Fatal("refusing to amend order; set OKX_CONFIRM=YES to continue")
	}

	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_API_SECRET")
	passphrase := os.Getenv("OKX_API_PASSPHRASE")
	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("missing env: OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE")
	}

	instId := os.Getenv("OKX_INST_ID")
	if instId == "" {
		log.Fatal("missing env: OKX_INST_ID")
	}

	ordId, hasOrdId := os.LookupEnv("OKX_ORD_ID")
	clOrdId, hasClOrdId := os.LookupEnv("OKX_CL_ORD_ID")
	if (hasOrdId && hasClOrdId) || (!hasOrdId && !hasClOrdId) {
		log.Fatal("require exactly one of OKX_ORD_ID or OKX_CL_ORD_ID")
	}

	newSz := os.Getenv("OKX_NEW_SZ")
	newPx := os.Getenv("OKX_NEW_PX")
	newPxUsd := os.Getenv("OKX_NEW_PX_USD")
	newPxVol := os.Getenv("OKX_NEW_PX_VOL")
	pxCount := 0
	if newPx != "" {
		pxCount++
	}
	if newPxUsd != "" {
		pxCount++
	}
	if newPxVol != "" {
		pxCount++
	}
	if pxCount > 1 {
		log.Fatal("too many price envs: require at most one of OKX_NEW_PX / OKX_NEW_PX_USD / OKX_NEW_PX_VOL")
	}
	if newSz == "" && newPx == "" && newPxUsd == "" && newPxVol == "" {
		log.Fatal("missing env: OKX_NEW_SZ or one of OKX_NEW_PX / OKX_NEW_PX_USD / OKX_NEW_PX_VOL")
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

	svc := c.NewAmendOrderService().InstId(instId)
	if hasOrdId {
		svc.OrdId(ordId)
	} else {
		svc.ClOrdId(clOrdId)
	}

	if v := os.Getenv("OKX_REQ_ID"); v != "" {
		svc.ReqId(v)
	}

	if newSz != "" {
		svc.NewSz(newSz)
	}
	if newPx != "" {
		svc.NewPx(newPx)
	}
	if newPxUsd != "" {
		svc.NewPxUsd(newPxUsd)
	}
	if newPxVol != "" {
		svc.NewPxVol(newPxVol)
	}

	if v := os.Getenv("OKX_PX_AMEND_TYPE"); v != "" {
		svc.PxAmendType(v)
	}
	if v := os.Getenv("OKX_CXL_ON_FAIL"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			log.Fatalf("invalid env OKX_CXL_ON_FAIL=%q: %v", v, err)
		}
		svc.CxlOnFail(b)
	}
	if v := os.Getenv("OKX_EXP_TIME"); v != "" {
		svc.ExpTime(v)
	}

	ack, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("trade_amend_order: instId=%s reqId=%s clOrdId=%s ordId=%s sCode=%s sMsg=%s", instId, ack.ReqId, ack.ClOrdId, ack.OrdId, ack.SCode, ack.SMsg)
}
