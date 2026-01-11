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

// OnOpReply 设置 WS 业务 op 回包回调（可在 Start 前或运行中设置；传 nil 表示清空）。
func (w *WSClient) OnOpReply(handler func(reply WSOpReply, raw []byte)) {
	if w == nil {
		return
	}
	w.typedMu.Lock()
	w.opReplyHandler = handler
	w.typedMu.Unlock()
}
