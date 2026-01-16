# 文档导航

建议阅读顺序：

1. [`guide.md`](guide.md)：快速上手 + 常用入口 + 错误处理 + 如何定位某个接口怎么用
2. [`ws.md`](ws.md)：WebSocket（心跳/重连/typed handler/orderbook store/交易 op）
3. [`coverage.md`](coverage.md)：覆盖矩阵（Endpoint -> Service/Test/Example）
4. [`roadmap.md`](roadmap.md)：迭代计划
5. [`design.md`](design.md)：设计记录（更偏贡献者/维护者视角）

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
