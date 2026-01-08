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
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/ws_public_tickers
```

订阅 WS 公共成交（首次收到消息后退出）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/ws_public_trades
```

订阅 WS 公共深度 books5（首次收到消息后退出）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/ws_public_books5
```

订阅 WS business K线（首次收到消息后退出）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_CANDLE_BAR="1m" # 可选：1m/5m/1H/1D...
go run ./examples/ws_business_candles
```

订阅 WS business 全部成交 trades-all（首次收到消息后退出）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/ws_business_trades_all
```

订阅 WS 私有订单频道（验证 subscribe/unsubscribe ack，不触发下单；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="ANY" # 可选：SPOT/MARGIN/SWAP/FUTURES/OPTION/ANY
export OKX_INST_ID="" # 可选
go run ./examples/ws_private_orders
```

订阅 WS 私有成交频道（仅验证 subscribe ack；需要 API Key，支持模拟盘；VIP6+ 才可用）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/ws_private_fills
```

订阅 WS 私有账户频道（收到首条 account 推送后退出；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_CCY="BTC" # 可选
export OKX_WS_EXTRA_PARAMS="{\"updateInterval\":\"0\"}" # 可选
go run ./examples/ws_private_account
```

订阅 WS 私有持仓频道（收到首条 positions 推送后退出；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="ANY" # 可选：MARGIN/SWAP/FUTURES/OPTION/ANY
export OKX_INST_FAMILY="" # 可选
export OKX_INST_ID="" # 可选
export OKX_WS_EXTRA_PARAMS="{\"updateInterval\":\"0\"}" # 可选
go run ./examples/ws_private_positions
```

订阅 WS 账户余额和持仓频道（收到首条 balance_and_position 推送后退出；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
go run ./examples/ws_private_balance_and_position
```

监听 WS 私有订单推送（收到第一条订单更新后退出；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="ANY" # 可选：SPOT/MARGIN/SWAP/FUTURES/OPTION/ANY
export OKX_INST_FAMILY="" # 可选
export OKX_INST_ID="" # 可选
export OKX_TIMEOUT="60s" # 可选：等待订单更新的超时
go run ./examples/ws_private_orders_stream
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

查询订单信息（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_ID="BTC-USDT" # 必填
export OKX_ORD_ID="123" # OKX_ORD_ID/OKX_CL_ORD_ID 二选一
export OKX_CL_ORD_ID="" # OKX_ORD_ID/OKX_CL_ORD_ID 二选一
go run ./examples/trade_get_order
```

批量下单（会实际下单；需要 API Key，支持模拟盘；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_BATCH_ORDERS='[{"instId":"BTC-USDT","tdMode":"cash","clOrdId":"b15","side":"buy","ordType":"limit","px":"2.15","sz":"2"}]'
export OKX_CONFIRM="YES"
go run ./examples/trade_batch_orders
```

批量撤单（会实际撤单；需要 API Key，支持模拟盘；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_BATCH_CANCEL_ORDERS='[{"instId":"BTC-USDT","ordId":"590908157585625111"}]'
export OKX_CONFIRM="YES"
go run ./examples/trade_cancel_batch_orders
```

批量改单（会实际改单；需要 API Key，支持模拟盘；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_BATCH_AMEND_ORDERS='[{"instId":"BTC-USDT","ordId":"590909145319051111","newSz":"2"}]'
export OKX_CONFIRM="YES"
go run ./examples/trade_amend_batch_orders
```

## 设计文档

- `docs/design.md`
