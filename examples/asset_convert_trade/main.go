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

	baseCcy := os.Getenv("OKX_BASE_CCY")
	quoteCcy := os.Getenv("OKX_QUOTE_CCY")
	side := os.Getenv("OKX_SIDE") // buy/sell
	sz := os.Getenv("OKX_SZ")
	szCcy := os.Getenv("OKX_SZ_CCY")
	if baseCcy == "" || quoteCcy == "" || side == "" || sz == "" || szCcy == "" {
		log.Fatal("missing env: OKX_BASE_CCY / OKX_QUOTE_CCY / OKX_SIDE / OKX_SZ / OKX_SZ_CCY")
	}

	demo := os.Getenv("OKX_DEMO") == "1"
	tag := os.Getenv("OKX_TAG")
	clTReqId := os.Getenv("OKX_CL_T_REQ_ID")
	confirm := os.Getenv("OKX_CONFIRM_CONVERT_TRADE") == "YES"

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

	quoteId := os.Getenv("OKX_QUOTE_ID")
	if quoteId == "" {
		q, err := c.NewAssetConvertEstimateQuoteService().
			BaseCcy(baseCcy).
			QuoteCcy(quoteCcy).
			Side(side).
			RfqSz(sz).
			RfqSzCcy(szCcy).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		quoteId = q.QuoteId
		log.Printf("estimated quote: quoteId=%s ttlMs=%s cnvtPx=%s", q.QuoteId, q.TtlMs, q.CnvtPx)
	}

	if !confirm {
		log.Fatalf("refusing to convert trade; set OKX_CONFIRM_CONVERT_TRADE=YES to continue (quoteId=%s)", quoteId)
	}

	svc := c.NewAssetConvertTradeService().
		QuoteId(quoteId).
		BaseCcy(baseCcy).
		QuoteCcy(quoteCcy).
		Side(side).
		Sz(sz).
		SzCcy(szCcy)
	if clTReqId != "" {
		svc.ClTReqId(clTReqId)
	}
	if tag != "" {
		svc.Tag(tag)
	}

	trade, err := svc.Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("trade: tradeId=%s state=%s instId=%s fillPx=%s fillBaseSz=%s fillQuoteSz=%s ts=%d", trade.TradeId, trade.State, trade.InstId, trade.FillPx, trade.FillBaseSz, trade.FillQuoteSz, trade.TS)
}
