package okx

import (
	"context"
	"fmt"
)

type wsTypedKind int

const (
	wsTypedKindOrders wsTypedKind = iota + 1
	wsTypedKindFills
	wsTypedKindAccount
	wsTypedKindPositions
	wsTypedKindBalanceAndPosition
	wsTypedKindDepositInfo
	wsTypedKindWithdrawalInfo
	wsTypedKindSprdOrders
	wsTypedKindSprdTrades
	wsTypedKindTickers
	wsTypedKindTrades
	wsTypedKindTradesAll
	wsTypedKindOrderBook
	wsTypedKindOpenInterest
	wsTypedKindFundingRate
	wsTypedKindMarkPrice
	wsTypedKindIndexTickers
	wsTypedKindPriceLimit
	wsTypedKindOptSummary
	wsTypedKindLiquidationOrders
	wsTypedKindOptionTrades
	wsTypedKindCallAuctionDetails
	wsTypedKindCandles
	wsTypedKindPriceCandles
	wsTypedKindSprdPublicTrades
	wsTypedKindSprdTickers
	wsTypedKindOpReply
)

func (k wsTypedKind) String() string {
	switch k {
	case wsTypedKindOrders:
		return "orders"
	case wsTypedKindFills:
		return "fills"
	case wsTypedKindAccount:
		return "account"
	case wsTypedKindPositions:
		return "positions"
	case wsTypedKindBalanceAndPosition:
		return "balance_and_position"
	case wsTypedKindDepositInfo:
		return "deposit_info"
	case wsTypedKindWithdrawalInfo:
		return "withdrawal_info"
	case wsTypedKindSprdOrders:
		return "sprd_orders"
	case wsTypedKindSprdTrades:
		return "sprd_trades"
	case wsTypedKindTickers:
		return "tickers"
	case wsTypedKindTrades:
		return "trades"
	case wsTypedKindTradesAll:
		return "trades_all"
	case wsTypedKindOrderBook:
		return "order_book"
	case wsTypedKindOpenInterest:
		return "open_interest"
	case wsTypedKindFundingRate:
		return "funding_rate"
	case wsTypedKindMarkPrice:
		return "mark_price"
	case wsTypedKindIndexTickers:
		return "index_tickers"
	case wsTypedKindPriceLimit:
		return "price_limit"
	case wsTypedKindOptSummary:
		return "opt_summary"
	case wsTypedKindLiquidationOrders:
		return "liquidation_orders"
	case wsTypedKindOptionTrades:
		return "option_trades"
	case wsTypedKindCallAuctionDetails:
		return "call_auction_details"
	case wsTypedKindCandles:
		return "candles"
	case wsTypedKindPriceCandles:
		return "price_candles"
	case wsTypedKindSprdPublicTrades:
		return "sprd_public_trades"
	case wsTypedKindSprdTickers:
		return "sprd_tickers"
	case wsTypedKindOpReply:
		return "op_reply"
	default:
		return "unknown"
	}
}

type wsTypedTask struct {
	kind wsTypedKind

	orders         []TradeOrder
	fills          []WSFill
	balances       []AccountBalance
	positions      []AccountPosition
	balPos         []WSBalanceAndPosition
	depositInfo    []WSDepositInfo
	withdrawalInfo []WSWithdrawalInfo
	sprdOrders     []SprdOrder
	sprdTrades     []WSSprdTrade
	tickers        []MarketTicker
	trades         []MarketTrade
	tradesAll      []MarketTrade
	orderBooks     []WSData[WSOrderBook]
	openInterests  []OpenInterest
	fundingRates   []FundingRate
	markPrices     []MarkPrice
	indexTickers   []IndexTicker
	priceLimits    []PriceLimit
	optSummaries   []OptSummary

	liquidationOrders  []LiquidationOrder
	optionTrades       []WSOptionTrade
	callAuctionDetails []WSCallAuctionDetails

	candles      []WSCandle
	priceCandles []WSPriceCandle

	sprdPublicTrades []WSSprdPublicTrade
	sprdTickers      []MarketSprdTicker

	op    WSOpReply
	opRaw []byte
}

func (w *WSClient) typedDispatchLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-w.typedQueue:
			w.handleTyped(task)
		}
	}
}

func (w *WSClient) dispatchTyped(task wsTypedTask) {
	if w == nil {
		return
	}

	if !w.typedAsync || w.typedQueue == nil {
		w.handleTyped(task)
		return
	}

	select {
	case w.typedQueue <- task:
	default:
		w.onError(fmt.Errorf("okx: ws typed handler queue full; dropped kind=%s", task.kind.String()))
	}
}

func (w *WSClient) handleTyped(task wsTypedTask) {
	if w == nil {
		return
	}

	switch task.kind {
	case wsTypedKindOrders:
		w.typedMu.RLock()
		h := w.ordersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.orders) == 0 {
			return
		}
		for _, order := range task.orders {
			o := order
			w.safeTypedCall(task.kind, func() { h(o) })
		}
	case wsTypedKindFills:
		w.typedMu.RLock()
		h := w.fillsHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.fills) == 0 {
			return
		}
		for _, fill := range task.fills {
			f := fill
			w.safeTypedCall(task.kind, func() { h(f) })
		}
	case wsTypedKindAccount:
		w.typedMu.RLock()
		h := w.accountHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.balances) == 0 {
			return
		}
		for _, balance := range task.balances {
			b := balance
			w.safeTypedCall(task.kind, func() { h(b) })
		}
	case wsTypedKindPositions:
		w.typedMu.RLock()
		h := w.positionsHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.positions) == 0 {
			return
		}
		for _, position := range task.positions {
			p := position
			w.safeTypedCall(task.kind, func() { h(p) })
		}
	case wsTypedKindBalanceAndPosition:
		w.typedMu.RLock()
		h := w.balanceAndPositionHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.balPos) == 0 {
			return
		}
		for _, data := range task.balPos {
			d := data
			w.safeTypedCall(task.kind, func() { h(d) })
		}
	case wsTypedKindDepositInfo:
		w.typedMu.RLock()
		h := w.depositInfoHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.depositInfo) == 0 {
			return
		}
		for _, info := range task.depositInfo {
			i := info
			w.safeTypedCall(task.kind, func() { h(i) })
		}
	case wsTypedKindWithdrawalInfo:
		w.typedMu.RLock()
		h := w.withdrawalInfoHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.withdrawalInfo) == 0 {
			return
		}
		for _, info := range task.withdrawalInfo {
			i := info
			w.safeTypedCall(task.kind, func() { h(i) })
		}
	case wsTypedKindSprdOrders:
		w.typedMu.RLock()
		h := w.sprdOrdersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.sprdOrders) == 0 {
			return
		}
		for _, order := range task.sprdOrders {
			o := order
			w.safeTypedCall(task.kind, func() { h(o) })
		}
	case wsTypedKindSprdTrades:
		w.typedMu.RLock()
		h := w.sprdTradesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.sprdTrades) == 0 {
			return
		}
		for _, trade := range task.sprdTrades {
			tr := trade
			w.safeTypedCall(task.kind, func() { h(tr) })
		}
	case wsTypedKindTickers:
		w.typedMu.RLock()
		h := w.tickersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.tickers) == 0 {
			return
		}
		for _, ticker := range task.tickers {
			tk := ticker
			w.safeTypedCall(task.kind, func() { h(tk) })
		}
	case wsTypedKindTrades:
		w.typedMu.RLock()
		h := w.tradesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.trades) == 0 {
			return
		}
		for _, trade := range task.trades {
			tr := trade
			w.safeTypedCall(task.kind, func() { h(tr) })
		}
	case wsTypedKindTradesAll:
		w.typedMu.RLock()
		h := w.tradesAllHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.tradesAll) == 0 {
			return
		}
		for _, trade := range task.tradesAll {
			tr := trade
			w.safeTypedCall(task.kind, func() { h(tr) })
		}
	case wsTypedKindOrderBook:
		w.typedMu.RLock()
		h := w.orderBookHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.orderBooks) == 0 {
			return
		}
		for _, data := range task.orderBooks {
			d := data
			w.safeTypedCall(task.kind, func() { h(d) })
		}
	case wsTypedKindOpenInterest:
		w.typedMu.RLock()
		h := w.openInterestHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.openInterests) == 0 {
			return
		}
		for _, oi := range task.openInterests {
			v := oi
			w.safeTypedCall(task.kind, func() { h(v) })
		}
	case wsTypedKindFundingRate:
		w.typedMu.RLock()
		h := w.fundingRateHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.fundingRates) == 0 {
			return
		}
		for _, rate := range task.fundingRates {
			r := rate
			w.safeTypedCall(task.kind, func() { h(r) })
		}
	case wsTypedKindMarkPrice:
		w.typedMu.RLock()
		h := w.markPriceHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.markPrices) == 0 {
			return
		}
		for _, price := range task.markPrices {
			p := price
			w.safeTypedCall(task.kind, func() { h(p) })
		}
	case wsTypedKindIndexTickers:
		w.typedMu.RLock()
		h := w.indexTickersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.indexTickers) == 0 {
			return
		}
		for _, ticker := range task.indexTickers {
			tk := ticker
			w.safeTypedCall(task.kind, func() { h(tk) })
		}
	case wsTypedKindPriceLimit:
		w.typedMu.RLock()
		h := w.priceLimitHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.priceLimits) == 0 {
			return
		}
		for _, limit := range task.priceLimits {
			pl := limit
			w.safeTypedCall(task.kind, func() { h(pl) })
		}
	case wsTypedKindOptSummary:
		w.typedMu.RLock()
		h := w.optSummaryHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.optSummaries) == 0 {
			return
		}
		for _, summary := range task.optSummaries {
			s := summary
			w.safeTypedCall(task.kind, func() { h(s) })
		}
	case wsTypedKindLiquidationOrders:
		w.typedMu.RLock()
		h := w.liquidationOrdersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.liquidationOrders) == 0 {
			return
		}
		for _, order := range task.liquidationOrders {
			o := order
			w.safeTypedCall(task.kind, func() { h(o) })
		}
	case wsTypedKindOptionTrades:
		w.typedMu.RLock()
		h := w.optionTradesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.optionTrades) == 0 {
			return
		}
		for _, trade := range task.optionTrades {
			tr := trade
			w.safeTypedCall(task.kind, func() { h(tr) })
		}
	case wsTypedKindCallAuctionDetails:
		w.typedMu.RLock()
		h := w.callAuctionDetailsHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.callAuctionDetails) == 0 {
			return
		}
		for _, detail := range task.callAuctionDetails {
			d := detail
			w.safeTypedCall(task.kind, func() { h(d) })
		}
	case wsTypedKindCandles:
		w.typedMu.RLock()
		h := w.candlesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.candles) == 0 {
			return
		}
		for _, candle := range task.candles {
			c := candle
			w.safeTypedCall(task.kind, func() { h(c) })
		}
	case wsTypedKindPriceCandles:
		w.typedMu.RLock()
		h := w.priceCandlesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.priceCandles) == 0 {
			return
		}
		for _, candle := range task.priceCandles {
			c := candle
			w.safeTypedCall(task.kind, func() { h(c) })
		}
	case wsTypedKindSprdPublicTrades:
		w.typedMu.RLock()
		h := w.sprdPublicTradesHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.sprdPublicTrades) == 0 {
			return
		}
		for _, trade := range task.sprdPublicTrades {
			tr := trade
			w.safeTypedCall(task.kind, func() { h(tr) })
		}
	case wsTypedKindSprdTickers:
		w.typedMu.RLock()
		h := w.sprdTickersHandler
		w.typedMu.RUnlock()
		if h == nil || len(task.sprdTickers) == 0 {
			return
		}
		for _, ticker := range task.sprdTickers {
			tk := ticker
			w.safeTypedCall(task.kind, func() { h(tk) })
		}
	case wsTypedKindOpReply:
		w.typedMu.RLock()
		h := w.opReplyHandler
		w.typedMu.RUnlock()
		if h == nil {
			return
		}
		w.safeTypedCall(task.kind, func() { h(task.op, task.opRaw) })
	default:
		return
	}
}

func (w *WSClient) safeTypedCall(kind wsTypedKind, fn func()) {
	defer func() {
		if r := recover(); r != nil {
			w.onError(fmt.Errorf("okx: ws typed handler panic kind=%s: %v", kind.String(), r))
		}
	}()
	fn()
}
