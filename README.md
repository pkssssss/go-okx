# go-okx

OKX V5 API 的 Go SDK（REST + WebSocket），目标是提供工程级的正确性与稳定性。

## 状态

开发中（尚未发布稳定版本）。

## 快速开始

获取系统时间：

```bash
go run ./examples/public_time
```

同步服务器时间并设置本地偏移：

```bash
go run ./examples/time_sync
```

读取账户余额（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_CCY="BTC,ETH" # 可选
go run ./examples/account_balance
```

查看持仓（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SWAP" # 可选
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/account_positions
```

订阅 WS 公共行情（首次收到消息后退出）：

```bash
go run ./examples/ws_public_tickers
```

订阅 WS 私有订单频道（仅验证 subscribe ack，不触发下单；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="ANY" # 可选：SPOT/MARGIN/SWAP/FUTURES/OPTION/ANY
export OKX_INST_ID="" # 可选
go run ./examples/ws_private_orders
```

获取单个产品行情（默认 BTC-USDT）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/market_ticker
```

获取产品深度（默认 BTC-USDT，默认 5 档）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_BOOKS_SZ="5" # 可选
go run ./examples/market_books
```

获取 K 线（默认 BTC-USDT，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
go run ./examples/market_candles
```

获取最近成交（默认 BTC-USDT，默认 20 条）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_TRADES_LIMIT="20" # 可选
go run ./examples/market_trades
```

查询产品信息（默认 SPOT）：

```bash
export OKX_INST_TYPE="SPOT" # 可选：SPOT/SWAP/FUTURES/OPTION
export OKX_ULY="BTC-USD" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/public_instruments
```

查询标记价格（默认 SWAP）：

```bash
export OKX_INST_TYPE="SWAP" # 可选：SWAP/FUTURES/OPTION
export OKX_ULY="BTC-USD" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/public_mark_price
```

查询资金费率（默认 BTC-USDT-SWAP）：

```bash
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/public_funding_rate
```

查询持仓总量（默认 SWAP）：

```bash
export OKX_INST_TYPE="SWAP" # 可选：SWAP/FUTURES/OPTION
export OKX_ULY="BTC-USD" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/public_open_interest
```

查询期权行情概要（默认 BTC-USD）：

```bash
export OKX_ULY="BTC-USD" # 可选
export OKX_OPT_EXP_TIME="260123" # 可选：YYMMDD
go run ./examples/public_opt_summary
```

查询未成交订单列表（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SPOT" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_ORD_TYPE="post_only" # 可选（可用逗号分隔多个）
export OKX_ORDER_STATE="live" # 可选
export OKX_LIMIT="100" # 可选
go run ./examples/trade_orders_pending
```

查询历史订单（近七天，需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SPOT" # 必填（示例默认 SPOT）
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_ORD_TYPE="limit" # 可选（可用逗号分隔多个）
export OKX_ORDER_STATE="filled" # 可选
export OKX_LIMIT="100" # 可选
go run ./examples/trade_orders_history
```

查询成交明细（近三天，需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SPOT" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_ORD_ID="123" # 可选
export OKX_LIMIT="100" # 可选
go run ./examples/trade_fills
```

## 设计文档

- `docs/design.md`
