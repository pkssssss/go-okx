# WebSocket 使用指南（OKX API v5）

> 模块：`github.com/pkssssss/go-okx/v5`  
> 目标：简洁 / 高效 / 稳定；默认优先“正确性与可判定行为”。

## 1. 选择 WS 端点（public / private / business）

- `client.NewWSPublic()`：`/ws/v5/public`，无需登录，适合行情/深度等公共数据。
- `client.NewWSPrivate()`：`/ws/v5/private`，需要登录，适合账户/订单/仓位等私有频道，且可用于 WS 交易 op（下单/撤单/改单）。
- `client.NewWSBusiness()`：`/ws/v5/business`，是否需要登录取决于频道（如 K 线无需登录；资金推送需要登录）。
- `client.NewWSBusinessPrivate()`：`/ws/v5/business` + 强制登录（如 `deposit-info` / `withdrawal-info` / `orders-algo` / `algo-advance` 等）。

> 建议：凡是需要登录的 WS，都先调用一次 `client.SyncTime(ctx)`，减少时间偏移导致的登录失败。

## 2. 生命周期与订阅

### 2.1 Start / Close / Done

- `Start(ctx, rawHandler, errHandler)`：启动后台 goroutine 读写 WS。
- `Close()`：主动关闭（会中断 `ReadMessage`，让 `Done()` 尽快返回）。
- `Done()`：等待 WS 退出。

### 2.2 SubscribeAndWait / UnsubscribeAndWait

订阅参数以 `WSArg` 为核心，不同频道的必填字段以 OKX 文档为准（常见：`InstType`/`InstId`/`InstFamily`/`AlgoId` 等）。  
例如：`grid-positions` / `grid-sub-orders` 频道订阅时必须提供 `AlgoId`。

推荐使用 `SubscribeAndWait(ctx, args...)`：

- 会等待 OKX 的 `event=subscribe` / `event=error`，用于判定订阅成功/失败；
- 断线后自动重连时，会自动重订阅已记录的 `desired` 订阅集合。

`Subscribe(args...)` 仅记录并尝试发送（不等待 ACK），适合你自己另有“订阅确认机制”的场景。

## 3. 心跳与断线（SDK 内置）

### 3.1 控制帧 ping/pong

OKX 会发送 opcode=9 的 ping，SDK 会自动用 opcode=10 的 pong 立即响应，并复制 payload。

### 3.2 文本 ping/pong（OKX 推荐）

OKX 文档建议：若 N 秒内未收到新消息，发送文本 `"ping"` 并期待 `"pong"`（N < 30s）。

SDK 默认启用（25s），可配置：

- `WithWSHeartbeat(15*time.Second)`：调小心跳间隔
- `WithWSHeartbeat(0)`：关闭文本心跳（不推荐，除非你在外部统一管理心跳/重连）

> 注意：`"ping"`/`"pong"` 文本消息会被 SDK 消费，不会进入 raw handler，避免业务侧误解析。

## 4. Typed handler（推荐）

### 4.1 为什么推荐 typed handler

- 业务侧无需重复 `json.Unmarshal` 与频道分发；
- 可与 `WithWSTypedHandlerAsync` 配合，把 handler 从 WS read goroutine 解耦，降低阻塞导致断线的概率；
- SDK 会捕获 typed handler 的 panic，并通过 `errHandler` 上报（避免进程崩溃）。

### 4.2 启用异步分发（可选）

SDK 默认启用异步分发（typed/raw 都是：buffer=1024），把 handler 从 WS read goroutine 解耦，降低“回调阻塞导致断线”的概率。

你仍可根据场景显式选择：

- typed handler：
  - `WithWSTypedHandlerAsync(n)`：异步（默认）
  - `WithWSTypedHandlerInline()`：inline（仅适合极轻 handler，避免额外队列/协程开销）
- raw handler：
  - `WithWSRawHandlerAsync(n)`：异步（默认）
  - `WithWSRawHandlerInline()`：inline

**背压语义（重要）**：

- 当队列满时，SDK 会通过 `errHandler` 上报 `"queue full; blocking"` 并阻塞等待入队（默认策略：`WSQueueFullBlock`，尽量不丢关键事件）。
- 若你订阅的是公开高频行情且对延迟更敏感，可使用 `WithWSQueueFullPolicy(WSQueueFullDrop)` 改为丢弃，或使用 `WithWSQueueFullPolicy(WSQueueFullDisconnect)` 在队列满时主动断线重连（Fail-Fast）。

建议：

- handler 只做轻量分发，把重逻辑交给你自己的 worker；
- 结合 `ws.Stats()` 监控 `TypedQueueLen/Cap`、`RawQueueLen/Cap`，并据此调大 buffer / 拆分 worker / 降载；
- 若你使用了 drop/disconnect 策略，请同时监控 `TypedDropped/RawDropped` 并在异常时触发 REST 对账或重建本地状态机。

## 5. 深度（Order Book）的正确用法

### 5.1 typed 收到的是“解析后的 WSData”

深度推送包含 `action=snapshot/update`、`seqId/prevSeqId`、`checksum` 等字段，正确用法通常需要：

- 合并 snapshot/update；
- 校验 prevSeqId 连续性；
- 校验 checksum（可选但建议开启）。

SDK 提供 `WSOrderBookStore` 来做这件事。

> 并发说明：`WSOrderBookStore` 并发安全；`Apply/ApplyMessage/Reset` 与 `Snapshot/Ready` 可并发调用（内部带锁与快照深拷贝）。为减少锁竞争，建议由单一 goroutine 串行 `Apply`，其他 goroutine 只读 `Snapshot`。

### 5.2 典型模式

1) 订阅 `books/books5/bbo-tbt/books-l2-tbt/...`（以及对应的 `sprd-*` 频道）  
2) `OnOrderBook` 收到 `WSData[WSOrderBook]`，直接 `store.Apply(&dm)`  
3) 通过 `store.Snapshot()` 读取最新快照

示例：见 `examples/ws_public_books_store_typed`。

## 6. K 线（Candles）的正确用法

OKX 的 K 线数据本身不带 `instId/sprdId`，需要通过订阅参数 `arg` 来携带上下文。

- `OnCandles(func(c okx.WSCandle))`：`candle*` / `sprd-candle*`，回调携带 `c.Arg` + `c.Candle`
- `OnPriceCandles(func(c okx.WSPriceCandle))`：`mark-price-candle*` / `index-candle*`，回调携带 `c.Arg` + `c.Candle`

示例：见 `examples/ws_business_candles_typed`。

## 7. raw handler / event handler / err handler

- raw handler：拿到原始 message bytes（注意 `"ping"/"pong"` 文本消息会被消费）
- event handler：`WithWSEventHandler` 获取 `subscribe/login/error/notice/channel-conn-count/...`
- err handler：运行时错误回调（断线、心跳超时、队列满、回调 panic 等）

SDK 会对 raw/event/err handler 的 panic 做保护（不会让 WS 主循环崩溃）。

## 8. 运行监控（Stats）

`ws.Stats()` 可返回一份并发安全的运行状态快照，便于你做指标化与告警（例如：最后收包时间、重连次数、订阅成功/失败计数、handler 队列堆积等）。
