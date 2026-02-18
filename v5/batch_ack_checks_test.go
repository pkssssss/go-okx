package okx

import (
	"errors"
	"net/http"
	"testing"
)

func TestTradeCheckBatchAcks_EmptyAcksFailClose(t *testing.T) {
	err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/batch-orders", "rid-trade-empty", 1, nil)
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

func TestTradeCheckBatchAcks_EmptySCodeFailClose(t *testing.T) {
	err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/batch-orders", "rid-trade-empty-scode", 1, []TradeOrderAck{{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradeBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradeBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-trade-empty-scode"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestTradeCheckAlgoAcks_EmptySCodeFailClose(t *testing.T) {
	err := tradeCheckAlgoAcks(http.MethodPost, "/api/v5/trade/cancel-algos", "rid-trade-algo-empty-scode", []TradeAlgoOrderAck{{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradeAlgoBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradeAlgoBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-trade-algo-empty-scode"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestTradingBotCheckBatchAcks_EmptyAcksFailClose(t *testing.T) {
	err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/signal/stop-order-algo", "rid-bot-empty", 1, nil)
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

func TestTradingBotCheckBatchAcks_EmptySCodeFailClose(t *testing.T) {
	err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/signal/stop-order-algo", "rid-bot-empty-scode", 1, []TradingBotOrderAck{{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradingBotBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradingBotBatchError", err)
	}
	if got, want := batchErr.RequestID, "rid-bot-empty-scode"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestRFQCheckCancelBatchRFQs_EmptyAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchRFQs(http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", "rid-rfq-empty", 1, nil)
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

func TestRFQCheckCancelBatchRFQs_EmptySCodeFailClose(t *testing.T) {
	err := rfqCheckCancelBatchRFQs(http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", "rid-rfq-empty-scode", 1, []RFQCancelAck{{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchRFQsError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchRFQsError", err)
	}
	if got, want := batchErr.RequestID, "rid-rfq-empty-scode"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestRFQCheckCancelBatchQuotes_EmptyAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchQuotes(http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", "rid-quote-empty", 1, nil)
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

func TestRFQCheckCancelBatchQuotes_EmptySCodeFailClose(t *testing.T) {
	err := rfqCheckCancelBatchQuotes(http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", "rid-quote-empty-scode", 1, []RFQCancelQuoteAck{{}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchQuotesError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchQuotesError", err)
	}
	if got, want := batchErr.RequestID, "rid-quote-empty-scode"; got != want {
		t.Fatalf("RequestID = %q, want %q", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestTradeCheckBatchAcks_ShortAcksFailClose(t *testing.T) {
	err := tradeCheckBatchAcks(http.MethodPost, "/api/v5/trade/batch-orders", "rid-trade-short", 2, []TradeOrderAck{{SCode: "0"}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradeBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradeBatchError", err)
	}
	if got, want := batchErr.Expected, 2; got != want {
		t.Fatalf("Expected = %d, want %d", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestTradingBotCheckBatchAcks_ShortAcksFailClose(t *testing.T) {
	err := tradingBotCheckBatchAcks(http.MethodPost, "/api/v5/tradingBot/signal/stop-order-algo", "rid-bot-short", 2, []TradingBotOrderAck{{SCode: "0"}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *TradingBotBatchError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *TradingBotBatchError", err)
	}
	if got, want := batchErr.Expected, 2; got != want {
		t.Fatalf("Expected = %d, want %d", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestRFQCheckCancelBatchRFQs_ShortAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchRFQs(http.MethodPost, "/api/v5/rfq/cancel-batch-rfqs", "rid-rfq-short", 2, []RFQCancelAck{{SCode: "0"}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchRFQsError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchRFQsError", err)
	}
	if got, want := batchErr.Expected, 2; got != want {
		t.Fatalf("Expected = %d, want %d", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}

func TestRFQCheckCancelBatchQuotes_ShortAcksFailClose(t *testing.T) {
	err := rfqCheckCancelBatchQuotes(http.MethodPost, "/api/v5/rfq/cancel-batch-quotes", "rid-quote-short", 2, []RFQCancelQuoteAck{{SCode: "0"}})
	if err == nil {
		t.Fatalf("expected error")
	}

	var batchErr *RFQCancelBatchQuotesError
	if !errors.As(err, &batchErr) {
		t.Fatalf("error = %T, want *RFQCancelBatchQuotesError", err)
	}
	if got, want := batchErr.Expected, 2; got != want {
		t.Fatalf("Expected = %d, want %d", got, want)
	}
	if got, want := len(batchErr.Acks), 1; got != want {
		t.Fatalf("Acks len = %d, want %d", got, want)
	}
}
