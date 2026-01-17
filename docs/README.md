# 文档导航

建议阅读顺序：

1. [`guide.md`](guide.md)：快速上手 + 常用入口 + 错误处理 + 如何定位某个接口怎么用
2. [`ws.md`](ws.md)：WebSocket（心跳/重连/typed handler/orderbook store/交易 op）
3. [`coverage.md`](coverage.md)：覆盖矩阵（Endpoint -> Service/Test/Example）
4. [`roadmap.md`](roadmap.md)：迭代计划
5. [`design.md`](design.md)：设计记录（更偏贡献者/维护者视角）
6. [`runbook.md`](runbook.md)：运行手册（生产故障处置与演练）

## 常用入口（API Quick Index）

### Client

- `okx.NewClient(...)`
- 常用 Option：`okx.WithCredentials(...)`、`okx.WithDemoTrading(true)`
- 建议：WS 登录前先 `c.SyncTime(ctx)`（降低时间偏差导致的登录失败）

### REST

- 统一风格：`c.NewXXXService().<Setters...>().Do(ctx)`
- 如何找 Service：优先查 [`coverage.md`](coverage.md)（每行直达 Service/Test/Example）

### WebSocket

- 端点选择：`c.NewWSPublic()` / `c.NewWSPrivate()` / `c.NewWSBusiness()` / `c.NewWSBusinessPrivate()`
- 订阅：`SubscribeAndWait` / `UnsubscribeAndWait`
- typed handler（推荐）：`ws.OnTickers/OnTrades/OnOrderBook/OnOrders/...`
- handler 较重：`okx.WithWSTypedHandlerAsync(1024)`
- 深度合并：`okx.NewWSOrderBookStore(channel, instId)`（配合 `OnOrderBook`）

### Error

- REST：`*okx.APIError`（支持 `errors.As`）
- 常用判定：`okx.IsAuthError` / `okx.IsRateLimitError` / `okx.IsTimeSkewError`

## 常用示例（Examples Quick Index）

> 提示：需要鉴权的示例建议先使用模拟盘（`OKX_DEMO=1`）。

### Public（无需鉴权）

- [`examples/public_time`](../examples/public_time)：服务器时间（最小连通性检查）
- [`examples/time_sync`](../examples/time_sync)：同步服务器时间并设置本地偏移（建议私有 REST/WS 之前先跑）
- [`examples/market_ticker`](../examples/market_ticker)：单产品行情（默认 BTC-USDT）
- [`examples/market_books`](../examples/market_books)：深度（默认 5 档）

### Private（需要鉴权）

- [`examples/account_balance`](../examples/account_balance)：账户余额
- [`examples/account_positions`](../examples/account_positions)：持仓
- [`examples/trade_orders_pending`](../examples/trade_orders_pending)：未成交订单（只读）

### WebSocket

- [`examples/ws_public_tickers`](../examples/ws_public_tickers)：公共行情订阅（tickers）
- [`examples/ws_public_books_store_typed`](../examples/ws_public_books_store_typed)：深度合并（snapshot/update）+ seq/checksum 校验（推荐）
- [`examples/ws_private_orders_stream`](../examples/ws_private_orders_stream)：私有订单推送（orders，需要鉴权）
