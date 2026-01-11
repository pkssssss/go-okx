package okx

// WithWSOrdersHandler 设置 orders 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSOrdersHandler(handler func(order TradeOrder)) WSOption {
	return func(c *WSClient) {
		c.OnOrders(handler)
	}
}

// WithWSFillsHandler 设置 fills 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSFillsHandler(handler func(fill WSFill)) WSOption {
	return func(c *WSClient) {
		c.OnFills(handler)
	}
}

// WithWSAccountHandler 设置 account 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSAccountHandler(handler func(balance AccountBalance)) WSOption {
	return func(c *WSClient) {
		c.OnAccount(handler)
	}
}

// WithWSPositionsHandler 设置 positions 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSPositionsHandler(handler func(position AccountPosition)) WSOption {
	return func(c *WSClient) {
		c.OnPositions(handler)
	}
}

// WithWSBalanceAndPositionHandler 设置 balance_and_position 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSBalanceAndPositionHandler(handler func(data WSBalanceAndPosition)) WSOption {
	return func(c *WSClient) {
		c.OnBalanceAndPosition(handler)
	}
}

// WithWSDepositInfoHandler 设置 deposit-info 推送的逐条回调（business WS，需要登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSDepositInfoHandler(handler func(info WSDepositInfo)) WSOption {
	return func(c *WSClient) {
		c.OnDepositInfo(handler)
	}
}

// WithWSWithdrawalInfoHandler 设置 withdrawal-info 推送的逐条回调（business WS，需要登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSWithdrawalInfoHandler(handler func(info WSWithdrawalInfo)) WSOption {
	return func(c *WSClient) {
		c.OnWithdrawalInfo(handler)
	}
}

// WithWSSprdOrdersHandler 设置 sprd-orders 推送的逐条回调（business WS，需要登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSSprdOrdersHandler(handler func(order SprdOrder)) WSOption {
	return func(c *WSClient) {
		c.OnSprdOrders(handler)
	}
}

// WithWSSprdTradesHandler 设置 sprd-trades 推送的逐条回调（business WS，需要登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSSprdTradesHandler(handler func(trade WSSprdTrade)) WSOption {
	return func(c *WSClient) {
		c.OnSprdTrades(handler)
	}
}

// WithWSTickersHandler 设置 tickers 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSTickersHandler(handler func(ticker MarketTicker)) WSOption {
	return func(c *WSClient) {
		c.OnTickers(handler)
	}
}

// WithWSTradesHandler 设置 trades 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSTradesHandler(handler func(trade MarketTrade)) WSOption {
	return func(c *WSClient) {
		c.OnTrades(handler)
	}
}

// WithWSTradesAllHandler 设置 trades-all 推送的逐条回调（business WS）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSTradesAllHandler(handler func(trade MarketTrade)) WSOption {
	return func(c *WSClient) {
		c.OnTradesAll(handler)
	}
}

// WithWSOrderBookHandler 设置深度（order book）推送回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
// 建议配合 WSOrderBookStore 合并 snapshot/update 并校验 seqId/checksum。
func WithWSOrderBookHandler(handler func(data WSData[WSOrderBook])) WSOption {
	return func(c *WSClient) {
		c.OnOrderBook(handler)
	}
}

// WithWSOpenInterestHandler 设置 open-interest 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSOpenInterestHandler(handler func(oi OpenInterest)) WSOption {
	return func(c *WSClient) {
		c.OnOpenInterest(handler)
	}
}

// WithWSFundingRateHandler 设置 funding-rate 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSFundingRateHandler(handler func(rate FundingRate)) WSOption {
	return func(c *WSClient) {
		c.OnFundingRate(handler)
	}
}

// WithWSMarkPriceHandler 设置 mark-price 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSMarkPriceHandler(handler func(price MarkPrice)) WSOption {
	return func(c *WSClient) {
		c.OnMarkPrice(handler)
	}
}

// WithWSIndexTickersHandler 设置 index-tickers 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSIndexTickersHandler(handler func(ticker IndexTicker)) WSOption {
	return func(c *WSClient) {
		c.OnIndexTickers(handler)
	}
}

// WithWSPriceLimitHandler 设置 price-limit 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSPriceLimitHandler(handler func(limit PriceLimit)) WSOption {
	return func(c *WSClient) {
		c.OnPriceLimit(handler)
	}
}

// WithWSOptSummaryHandler 设置 opt-summary 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSOptSummaryHandler(handler func(summary OptSummary)) WSOption {
	return func(c *WSClient) {
		c.OnOptSummary(handler)
	}
}

// WithWSLiquidationOrdersHandler 设置 liquidation-orders 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSLiquidationOrdersHandler(handler func(order LiquidationOrder)) WSOption {
	return func(c *WSClient) {
		c.OnLiquidationOrders(handler)
	}
}

// WithWSOptionTradesHandler 设置 option-trades 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSOptionTradesHandler(handler func(trade WSOptionTrade)) WSOption {
	return func(c *WSClient) {
		c.OnOptionTrades(handler)
	}
}

// WithWSCallAuctionDetailsHandler 设置 call-auction-details 推送的逐条回调。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSCallAuctionDetailsHandler(handler func(detail WSCallAuctionDetails)) WSOption {
	return func(c *WSClient) {
		c.OnCallAuctionDetails(handler)
	}
}

// WithWSCandlesHandler 设置 K线推送的逐条回调（candle* / sprd-candle*，business WS）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSCandlesHandler(handler func(candle WSCandle)) WSOption {
	return func(c *WSClient) {
		c.OnCandles(handler)
	}
}

// WithWSPriceCandlesHandler 设置指数/标记价格K线推送的逐条回调（mark-price-candle* / index-candle*，business WS）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSPriceCandlesHandler(handler func(candle WSPriceCandle)) WSOption {
	return func(c *WSClient) {
		c.OnPriceCandles(handler)
	}
}

// WithWSSprdPublicTradesHandler 设置 sprd-public-trades 推送的逐条回调（business WS，无需登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSSprdPublicTradesHandler(handler func(trade WSSprdPublicTrade)) WSOption {
	return func(c *WSClient) {
		c.OnSprdPublicTrades(handler)
	}
}

// WithWSSprdTickersHandler 设置 sprd-tickers 推送的逐条回调（business WS，无需登录）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSSprdTickersHandler(handler func(ticker MarketSprdTicker)) WSOption {
	return func(c *WSClient) {
		c.OnSprdTickers(handler)
	}
}

// WithWSOpReplyHandler 设置 WS 业务 op 回包回调（order/cancel-order/amend-order 等）。
// 注意：默认在 WS read goroutine 中执行；若启用 WithWSTypedHandlerAsync，则在独立 worker goroutine 中执行。
func WithWSOpReplyHandler(handler func(reply WSOpReply, raw []byte)) WSOption {
	return func(c *WSClient) {
		c.OnOpReply(handler)
	}
}

// OnOrders 设置 orders 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOrders(handler func(order TradeOrder)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.ordersHandler = handler
	w.typedMu.Unlock()
}

// OnFills 设置 fills 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnFills(handler func(fill WSFill)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.fillsHandler = handler
	w.typedMu.Unlock()
}

// OnAccount 设置 account 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnAccount(handler func(balance AccountBalance)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.accountHandler = handler
	w.typedMu.Unlock()
}

// OnPositions 设置 positions 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnPositions(handler func(position AccountPosition)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.positionsHandler = handler
	w.typedMu.Unlock()
}

// OnBalanceAndPosition 设置 balance_and_position 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnBalanceAndPosition(handler func(data WSBalanceAndPosition)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.balanceAndPositionHandler = handler
	w.typedMu.Unlock()
}

// OnDepositInfo 设置 deposit-info 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnDepositInfo(handler func(info WSDepositInfo)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.depositInfoHandler = handler
	w.typedMu.Unlock()
}

// OnWithdrawalInfo 设置 withdrawal-info 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnWithdrawalInfo(handler func(info WSWithdrawalInfo)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.withdrawalInfoHandler = handler
	w.typedMu.Unlock()
}

// OnSprdOrders 设置 sprd-orders 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnSprdOrders(handler func(order SprdOrder)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.sprdOrdersHandler = handler
	w.typedMu.Unlock()
}

// OnSprdTrades 设置 sprd-trades 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnSprdTrades(handler func(trade WSSprdTrade)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.sprdTradesHandler = handler
	w.typedMu.Unlock()
}

// OnTickers 设置 tickers 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnTickers(handler func(ticker MarketTicker)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.tickersHandler = handler
	w.typedMu.Unlock()
}

// OnTrades 设置 trades 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnTrades(handler func(trade MarketTrade)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.tradesHandler = handler
	w.typedMu.Unlock()
}

// OnTradesAll 设置 trades-all 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnTradesAll(handler func(trade MarketTrade)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.tradesAllHandler = handler
	w.typedMu.Unlock()
}

// OnOrderBook 设置深度（order book）推送回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOrderBook(handler func(data WSData[WSOrderBook])) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.orderBookHandler = handler
	w.typedMu.Unlock()
}

// OnOpenInterest 设置 open-interest 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOpenInterest(handler func(oi OpenInterest)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.openInterestHandler = handler
	w.typedMu.Unlock()
}

// OnFundingRate 设置 funding-rate 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnFundingRate(handler func(rate FundingRate)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.fundingRateHandler = handler
	w.typedMu.Unlock()
}

// OnMarkPrice 设置 mark-price 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnMarkPrice(handler func(price MarkPrice)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.markPriceHandler = handler
	w.typedMu.Unlock()
}

// OnIndexTickers 设置 index-tickers 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnIndexTickers(handler func(ticker IndexTicker)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.indexTickersHandler = handler
	w.typedMu.Unlock()
}

// OnPriceLimit 设置 price-limit 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnPriceLimit(handler func(limit PriceLimit)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.priceLimitHandler = handler
	w.typedMu.Unlock()
}

// OnOptSummary 设置 opt-summary 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOptSummary(handler func(summary OptSummary)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.optSummaryHandler = handler
	w.typedMu.Unlock()
}

// OnLiquidationOrders 设置 liquidation-orders 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnLiquidationOrders(handler func(order LiquidationOrder)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.liquidationOrdersHandler = handler
	w.typedMu.Unlock()
}

// OnOptionTrades 设置 option-trades 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOptionTrades(handler func(trade WSOptionTrade)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.optionTradesHandler = handler
	w.typedMu.Unlock()
}

// OnCallAuctionDetails 设置 call-auction-details 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnCallAuctionDetails(handler func(detail WSCallAuctionDetails)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.callAuctionDetailsHandler = handler
	w.typedMu.Unlock()
}

// OnCandles 设置 K线推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnCandles(handler func(candle WSCandle)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.candlesHandler = handler
	w.typedMu.Unlock()
}

// OnPriceCandles 设置指数/标记价格 K线推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnPriceCandles(handler func(candle WSPriceCandle)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.priceCandlesHandler = handler
	w.typedMu.Unlock()
}

// OnSprdPublicTrades 设置 sprd-public-trades 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnSprdPublicTrades(handler func(trade WSSprdPublicTrade)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.sprdPublicTradesHandler = handler
	w.typedMu.Unlock()
}

// OnSprdTickers 设置 sprd-tickers 推送的逐条回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnSprdTickers(handler func(ticker MarketSprdTicker)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.sprdTickersHandler = handler
	w.typedMu.Unlock()
}

// OnOpReply 设置 WS 业务 op 回包回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOpReply(handler func(reply WSOpReply, raw []byte)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.opReplyHandler = handler
	w.typedMu.Unlock()
}
