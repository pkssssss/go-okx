package okx

import (
	"errors"
	"net/http"
	"testing"
)

func TestTradeCheckBatchAcks_EmptyAcksFailClose(t *testing.T) {
	err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/batch-orders", "rid-trade-empty", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradeBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradeBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-trade-empty"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if len(batchErr.Acks) != 0 {
		t.Fatalf("Acks len = %d, want 0", len(batchErr.Acks))
	}
}

func TestTradeCheckAlgoAcks_EmptyAcksFailClose(t *testing.T) {
	err := tradeCheckAlgoAcks(http.MethodPost, "/api/v5/trade/cancel-algos", "rid-trade-algo-empty", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradeAlgoBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradeAlgoBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-trade-algo-empty"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if len(batchErr.Acks) != 0 {
		t.Fatalf("Acks len = %d, want 0", len(batchErr.Acks))
	}
}

func TestTradingBotCheckBatchAcks_EmptyAcksFailClose(t *testing.T) {
	err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/signal/stop-order-algo", "rid-bot-empty", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradingBotBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradingBotBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-bot-empty"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if len(batchErr.Acks) != 0 {
		t.Fatalf("Acks len = %d, want 0", len(batchErr.Acks))
	}
}

func TestRFQCheckCancelBatchRFQs_EmptyAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchRFQs(http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", "rid-rfq-empty", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchRFQsError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchRFQsError", err)
	}
	if got, want := batchErr.RequestID, "rid-rfq-empty"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if len(batchErr.Acks) != 0 {
		t.Fatalf("Acks len = %d, want 0", len(batchErr.Acks))
	}
}

func TestRFQCheckCancelBatchQuotes_EmptyAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchQuotes(http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", "rid-quote-empty", nil)
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchQuotesError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchQuotesError", err)
	}
	if got, want := batchErr.RequestID, "rid-quote-empty"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if len(batchErr.Acks) != 0 {
		t.Fatalf("Acks len = %d, want 0", len(batchErr.Acks))
	}
}
