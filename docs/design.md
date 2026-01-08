# Go OKX SDK 设计草案（v0.1）

> 仓库：`github.com/pkssssss/go-okx`  
> Go 最低版本：`go1.25`  
> 目标：正确性、稳定性优先，其次才是易用性与性能；避免过度设计（KISS / YAGNI）。

## 1. 背景与目标

我们要做的是一个**工程级可用**的 OKX Go SDK，对标成熟的 `go-binance` 的“骨架能力”：

- **REST + WebSocket** 双通道齐备，并能形成交易闭环（下单/撤单/查单 + WS 推送回报）。
- **现货 + 合约（永续/交割）+ 期权**：v0.1 做到“能跑通主链路”，期权的专项能力按需增量补齐。
- **正确性**：签名、时间戳、参数编码、错误处理必须可测、可复现、可判定。
- **高可用/稳定性**：WS 自动重连、自动重订阅、心跳与升级通知处理；REST 超时、可控重试（仅幂等 GET 默认开启）。
- **可维护**：端点实现保持“薄”，通用能力沉淀为“厚”的管线与组件（DRY）。

## 2. 明确非目标（v0.1 不做）

- **SBE（二进制）行情**解码与专用 WS（`/ws/v5/public-sbe`）不进 v0.1（除非你明确需要极致性能）。
- 资产高级域（充提、子账户、借贷、申购赎回等）暂不覆盖。
- 过度抽象的“统一交易所接口层”（SDK 做通用 OKX 能力，量化内核自行适配 SDK）。

## 3. 关键差异点（OKX vs Binance，需要在架构里吸收）

### 3.1 鉴权/签名差异（必须统一封装）

OKX V5 鉴权的关键点（来自官方文档）：

- REST 私有请求头必须包含：
  - `OK-ACCESS-KEY`
  - `OK-ACCESS-SIGN`
  - `OK-ACCESS-TIMESTAMP`（文档示例为 **ISO UTC**，如 `2020-03-28T12:21:41.274Z`）
  - `OK-ACCESS-PASSPHRASE`
  - 模拟盘额外：`x-simulated-trading: 1`
- REST 签名串：`timestamp + method + requestPath + body`
  - `method`：大写 `GET/POST/...`
  - `requestPath`：包含 query string（GET 参数算在 requestPath，不算 body）
  - `body`：JSON 字符串；GET 通常为空字符串
- WS 登录（JSON 消息 `op=login`）：
  - `timestamp` 为 **Unix 秒**
  - 签名串：`timestamp + "GET" + "/users/self/verify"`
- WS-SBE（握手 header 登录）：
  - header 里的 `OK-ACCESS-TIMESTAMP` 为 **Unix 秒**
  - 签名串同上：`timestamp + "GET" + "/users/self/verify"`

> 结论：签名模块必须支持“多种 timestamp 格式 + 多种 prehash 规则”，并保证 **用于签名的 requestPath 与实际发送的 URL 完全一致**（同一处生成，避免漂移）。

### 3.2 错误形态差异（HTTP 200 也可能是业务失败）

OKX 常见返回 envelope：`code/msg/data`：

- HTTP 可能是 200，但 `code != "0"` 即失败。
- 交易类接口（尤其 batch）可能出现：
  - `code == "0"`，但 `data[i].sCode != "0"` 表示该条失败（需要把“部分失败”也当成 error 或至少可判定地返回）。

> 结论：错误模型必须把 “HTTP 错误 / 顶层 code 错误 / data 内部 sCode 错误”统一起来，且支持 `errors.As` 拿到结构化信息。

### 3.3 WebSocket 心跳与升级通知

- OKX 服务器会发 opcode=9 的 ping，客户端需尽快以 opcode=10 pong 响应，并**复制 ping payload**。
- 会推送 `event=notice, code=64008`，提示 60 秒后升级断线，建议主动重连。

> 结论：WS 连接管理必须内置“心跳处理 + notice 触发重连 + 自动重订阅”。

## 4. API 形态（对标 go-binance 的简洁体验）

v0.1 采用与 `go-binance` 类似的“Service + Do(ctx)”风格，原因：

- 调用方熟悉度高；
- 端点实现天然“薄”，便于 DRY；
- 可逐步扩展而不破坏兼容。

示意（非代码承诺，仅描述风格）：

- `client.NewMarketTickersService().InstType(...).Do(ctx)`
- `client.NewTradePlaceOrderService().InstId(...).Side(...).Sz(...).Do(ctx)`
- `ws := client.NewWSPrivate()` / `ws.SubscribeOrders(...)` / `ws.OnOrder(...)`

## 5. 包结构（建议）

保持 OKX v5 域的划分，避免按“现货/合约/期权”重复分层（因为 OKX 大多通过 `instType` 区分）：

```
go-okx/
  docs/
    design.md
  internal/
    rest/              # REST 通用管线（请求构建/签名/发送/解包/错误）
    ws/                # WS 连接管理（心跳/重连/重订阅/路由）
    sign/              # 签名与时间源（REST/WS 差异收敛）
    jsonx/             # 无损数值/时间戳解析工具（可选）
  okx.go               # package okx：Client、Option、公共类型入口
  public_*.go           # Public 服务（时间、产品、funding/mark/open-interest 等）
  market_*.go           # Market 服务（ticker/books/candles/trades 等）
  trade_*.go            # Trade 服务（下单/撤单/改单/查单/成交）
  account_*.go          # Account 服务（余额/仓位/杠杆等）
  ws_public_*.go        # WS public（订阅/事件模型）
  ws_private_*.go       # WS private（登录/订阅/事件模型）
  examples/             # 最小可运行示例（v0.1 必备）
```

> 说明：`internal/` 下放“厚管线”，外层文件放“薄端点”。这样新增端点时基本只写参数和返回类型。

## 6. 核心抽象与配置（v0.1 必备）

### 6.1 Client

`Client` 负责聚合配置与复用组件（不承载业务逻辑）：

- REST：
  - BaseURL：实盘 `https://www.okx.com`
  - Demo：同 BaseURL，但额外 header `x-simulated-trading: 1`
- WS：
  - public/private/business：实盘 `wss://ws.okx.com:8443/ws/v5/...`
  - demo：`wss://wspap.okx.com:8443/ws/v5/...`
- HTTPClient / Dialer：可注入，便于自定义代理、超时、监控与测试。
- Logger：可注入，默认安静；Debug 模式下输出请求/响应（注意脱敏）。
- Clock/TimeSource：可注入；支持 server time 校准（避免签名过期）。
- Credentials：`apiKey/secret/passphrase` + 可选的模拟盘开关。

### 6.2 Option 模式（KISS）

采用函数式 Option 配置（避免构造函数爆炸）：

- `WithDemoTrading(true)`：demo 交易（REST header + WS host 切换）
- `WithHTTPClient(*http.Client)`
- `WithUserAgent(string)`
- `WithLogger(Logger)`
- `WithTimeOffset(time.Duration)` 或 `WithTimeSync(...)`

## 7. REST 通用管线设计（正确性关键）

### 7.1 requestPath 生成必须“单一来源”

REST 签名依赖 `requestPath`（含 query），因此必须：

1) 在构建请求时生成 `requestPath`（`/api/v5/...` + `?k=v&...`）  
2) **同一份字符串**同时用于：
   - 实际请求 URL 拼接（BaseURL + requestPath）
   - 签名 prehash（timestamp + method + requestPath + body）

这样可避免“签名用的 query 顺序/编码”与“实际请求 URL”不一致导致的隐性 401。

### 7.2 Body 统一 JSON

OKX REST 请求体通常为 JSON：

- POST/PUT：body 为 JSON 字符串（用于签名）
- GET：body 为空字符串（用于签名），参数都走 query

### 7.3 统一响应 envelope 解包

统一结构（示意）：

- `code`（string）
- `msg`（string）
- `data`（json）

并提供两条解包路径：

1) 顶层 `code != "0"`：直接返回 `APIError`  
2) 顶层成功但存在 `data[].sCode != "0"`：  
   - 若是 batch 端点：返回 `BatchError`（携带每条失败项），或提供配置决定“部分失败是否视为 error”。

### 7.4 错误模型（必须可判定）

建议：

- `type APIError struct {`
  - `HTTPStatus int`
  - `Code string`
  - `Message string`
  - `RequestID string`（若响应头有）
  - `Raw []byte`
  - `Endpoint string` / `Method string` / `RequestPath string`（便于定位）
  - `}`
- 支持：
  - `errors.As(err, *APIError)`
  - 分类函数：`IsAuthError / IsRateLimit / IsTimestampSkew`（v0.1 先做最小集合）

## 8. WebSocket 设计（稳定性关键）

### 8.1 连接形态

OKX WS 端点：

- public：`/ws/v5/public`
- private：`/ws/v5/private`
- business：`/ws/v5/business`（按需）

demo 对应 host：`wspap.okx.com`

### 8.2 心跳与断线策略

- 设置 PingHandler：收到 ping(opcode=9) 后，立即回 pong(opcode=10)，payload 原样复制。
- 收到 `event=notice code=64008`：触发“主动重连”，并在新连接上恢复订阅（避免被动断线导致数据空窗）。

### 8.3 自动重连 + 自动重订阅（状态机）

核心要求：

- 订阅请求必须幂等化：SDK 维护 `desiredSubscriptions` 集合；
- 断线后重连成功：
  - public：直接重发 subscribe
  - private：先 login 成功，再重发 subscribe
- 对外暴露：
  - `OnConnect/OnDisconnect/OnReconnected` 钩子
  - `Close()` 明确停止（不再自动重连）

### 8.4 消息路由（最小但够用）

OKX 消息大致分两类：

- 控制类：`event=login/subscribe/unsubscribe/error/notice`
- 数据类：包含 `arg` 与 `data`

v0.1 路由策略：

- 先按 `event` 分发控制消息；
- 数据消息按 `arg.channel + arg.instId/arg.instType/...` 作为 key 路由给对应订阅的 handler。

为了降低上层（量化内核）对 raw message 的耦合，SDK 在 WSClient 内置最小的 typed handler：

- `ws.OnOrders(func(TradeOrder){...})`：逐条处理 orders 推送
- `ws.OnFills(func(WSFill){...})`：逐条处理 fills 推送
- `ws.OnAccount(func(AccountBalance){...})`：逐条处理 account 推送
- `ws.OnPositions(func(AccountPosition){...})`：逐条处理 positions 推送
- `ws.OnBalanceAndPosition(func(WSBalanceAndPosition){...})`：逐条处理 balance_and_position 推送
- `ws.OnOpReply(func(WSOpReply, []byte){...})`：观测业务 op 回包（含 raw 便于日志/审计）

高可用建议：

- 若 handler 逻辑较重，建议启用 `WithWSTypedHandlerAsync(buffer)` 将 typed handler 移到独立 worker goroutine 执行（避免阻塞 read goroutine）。
- 队列满时会丢弃该条 typed 回调，并通过 `errHandler` 报错；调用方可调大 buffer 或优化 handler。

### 8.5 业务 op 请求/响应（交易链路闭环）

OKX WS 除了 `event` 与 `arg+data` 推送外，还有一类“业务操作回包”：

- 请求：`{"id":"...","op":"order|cancel-order|amend-order", "args":[...]}`  
- 响应：`{"id":"...","op":"...","code":"0|...","msg":"...","data":[...],"inTime":"...","outTime":"..."}`

v0.1 设计要点：

- `id` 是关联请求与响应的唯一键：SDK 在发送 op 时注册 waiter，以 `id` 匹配回包并唤醒调用方。
- `event=error` 且带 `id` 时，应直接失败对应 waiter（避免调用方超时等待）。
- 断线时应立即失败所有未完成 waiter（避免调用方一直挂到 ctx 超时）。
- 交易 op 的错误需要“可判定”：
  - 顶层 `code != "0"` 视为失败；
  - 顶层成功但 `data[0].sCode != "0"` 仍视为失败（交易接口常见）。

对外 API（v0.1）：

- `(*WSClient).PlaceOrder/CancelOrder/AmendOrder(ctx, arg)`：返回 `*TradeOrderAck` 或 `error`
- `(*WSClient).PlaceOrders/CancelOrders/AmendOrders(ctx, args...)`：返回 `[]TradeOrderAck` 或 `error`
- 错误类型：`*WSTradeOpError`（携带 `op/id/code/msg/sCode/sMsg/inTime/outTime/raw`，便于上层诊断与告警）
- 批量部分失败：`*WSTradeOpBatchError`（顶层成功但存在 `sCode!=0`，返回 `acks` 并携带 `raw` 便于诊断）

## 9. 精度与类型策略（正确性优先）

- 金额/数量/费率/价格等小数：SDK 层保持 `string`（无损），不使用 `float64`。
  - 可选：定义 `type Decimal string`，提高语义清晰度。
- 时间戳：常见为字符串的毫秒/微秒时间戳，解析为 `int64`（整数不丢精度）。
- 枚举：用 `type Xxx string` + 常量集，避免 magic string。

## 10. v0.1 MVP 范围（按“最基础功能完善”定义）

### 10.1 REST（最低集合）

- 校时：`GET /api/v5/public/time`
- 产品：`GET /api/v5/public/instruments`（覆盖 `SPOT/SWAP/FUTURES/OPTION`）
- 行情：`/api/v5/market/ticker`、`/api/v5/market/tickers`、`/api/v5/market/books`、`/api/v5/market/candles`、`/api/v5/market/trades`
- 合约增强：`/api/v5/public/mark-price`、`/api/v5/public/funding-rate`、`/api/v5/public/open-interest`
- 期权最小增强：`/api/v5/public/opt-summary`
- 账户：`/api/v5/account/balance`、`/api/v5/account/positions`
- 交易：`/api/v5/trade/order`、`/api/v5/trade/cancel-order`、`/api/v5/trade/amend-order`、`/api/v5/trade/orders-pending`、`/api/v5/trade/orders-history`、`/api/v5/trade/fills`

### 10.2 WebSocket（最低集合）

- public：至少支持一个核心频道（建议先做 `books` 或 `tickers`，再扩到 `trades/candle`）
- private：订单与成交推送（形成闭环）；仓位/余额推送作为强烈建议项

## 11. 测试策略（SDK 正确性的“地基”）

### 11.1 单元测试（必须）

- 签名：用固定的 timestamp、method、requestPath、body，断言 `OK-ACCESS-SIGN` 与预期一致（可直接对照官方公式自行生成测试向量）。
- requestPath：断言 query 编码与签名使用的 requestPath 完全一致（同一生成逻辑）。
- 错误解包：覆盖 `HTTP!=200`、`code!=0`、`code==0 but sCode!=0`。

### 11.2 组件测试（建议）

- REST：用 `httptest.Server` 模拟 OKX 返回，验证 header、签名串拼接、body JSON 等。
- WS：用本地 websocket server 模拟：
  - ping/pong payload 回显
  - `notice(64008)` 触发重连
  - 断线重连后重订阅

### 11.3 集成测试（可选、默认关闭）

使用环境变量开启（避免 CI/本地默认打真实交易所）：

- `OKX_API_KEY / OKX_API_SECRET / OKX_API_PASSPHRASE`
- `OKX_DEMO=1`

## 12. 迭代方式（Git 驱动、风险最小化）

- 先落地“通用管线 + 最小端点 + 最小示例”，做到可运行。
- 每次新增端点都遵循：
  1) 端点薄封装（只组参数/类型）
  2) 复用统一管线（签名/请求/解包/错误）
  3) 增补测试向量（至少覆盖签名/错误）

## 13. 与 go-binance 的映射（我们要“继承”的不是代码，是骨架）

我们借鉴的关键点：

- `Client` 聚合配置；端点以 `NewXXXService().Do(ctx)` 形式暴露；
- “薄 service + 厚 callAPI”的复用组织；
- WS keepalive 的工程化细节（但 OKX 的 ping/pong 与 notice 需要我们增强为自动重连/重订阅状态机）；
- 测试驱动的正确性保障（签名、解析、错误都能单测）。

OKX 需要我们补齐/改造的点：

- REST/WS/SBE 的 timestamp 与 prehash 规则不同；
- 错误不止 HTTP 4xx/5xx，还包含 `code/msg` 与 `sCode/sMsg`；
- demo trading 的 header 与 WS host 切换；
- WS 侧必须处理 ping payload 与升级 notice。
