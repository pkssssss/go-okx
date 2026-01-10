# go-okx

OKX V5 API 的 Go SDK（REST + WebSocket），目标是提供工程级的正确性与稳定性。

## 状态

开发中（尚未发布稳定版本）。

## 安装

```bash
go get github.com/pkssssss/go-okx/v5
```

## 快速开始

仓库使用 Go workspace（`go.work`），可直接从根目录运行 `examples/`；SDK 主模块位于 `v5/`。

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

订阅 WS 公共深度 books（本地合并 + seq/checksum 校验，收到两条推送后退出）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_BOOKS_CHANNEL="books" # 可选：books/books-elp/books5/bbo-tbt/books-l2-tbt/books50-l2-tbt
go run ./examples/ws_public_books
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

监听 WS 私有订单推送（异步 typed handler；适合在 handler 内做较重逻辑；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="ANY" # 可选：SPOT/MARGIN/SWAP/FUTURES/OPTION/ANY
export OKX_INST_FAMILY="" # 可选
export OKX_INST_ID="" # 可选
export OKX_WS_ASYNC_BUFFER="1024" # 可选：typed handler 异步队列大小
export OKX_HANDLER_SLEEP="500ms" # 可选：模拟 handler 耗时（建议 <= timeout）
export OKX_TIMEOUT="60s" # 可选：等待订单更新的超时
go run ./examples/ws_private_orders_stream_async
```

通过 WS 下单（会真实下单；需要 API Key；建议先用模拟盘；必须显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 强烈建议：1=模拟盘
export OKX_CONFIRM=YES # 必填：防止误触真实下单
export OKX_INST_ID="BTC-USDT" # 必填
export OKX_TD_MODE="cash" # 必填：cash/cross/isolated...
export OKX_SIDE="buy" # 必填：buy/sell
export OKX_ORD_TYPE="market" # 必填：market/limit/...
export OKX_SZ="0.001" # 必填
export OKX_PX="..." # 可选：limit/post_only/... 需要 px/pxUsd/pxVol 之一
export OKX_TIMEOUT="10s" # 可选
go run ./examples/ws_private_trade_order
```

通过 WS 撤单（会真实撤单；需要 API Key；建议先用模拟盘；必须显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 强烈建议：1=模拟盘
export OKX_CONFIRM=YES # 必填：防止误触真实撤单
export OKX_INST_ID="BTC-USDT" # 必填
export OKX_ORD_ID="..." # OKX_ORD_ID 与 OKX_CL_ORD_ID 二选一
export OKX_CL_ORD_ID="" # OKX_ORD_ID 与 OKX_CL_ORD_ID 二选一
export OKX_TIMEOUT="10s" # 可选
go run ./examples/ws_private_trade_cancel
```

通过 WS 改单（会真实改单；需要 API Key；建议先用模拟盘；必须显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 强烈建议：1=模拟盘
export OKX_CONFIRM=YES # 必填：防止误触真实改单
export OKX_INST_ID="BTC-USDT" # 必填
export OKX_ORD_ID="..." # OKX_ORD_ID 与 OKX_CL_ORD_ID 二选一
export OKX_CL_ORD_ID="" # OKX_ORD_ID 与 OKX_CL_ORD_ID 二选一
export OKX_NEW_SZ="..." # OKX_NEW_SZ / OKX_NEW_PX / OKX_NEW_PX_USD / OKX_NEW_PX_VOL 至少一个
export OKX_NEW_PX="" # 可选
export OKX_NEW_PX_USD="" # 可选
export OKX_NEW_PX_VOL="" # 可选
export OKX_REQ_ID="" # 可选
export OKX_TIMEOUT="10s" # 可选
go run ./examples/ws_private_trade_amend
```

通过 WS 批量交易 op（会真实下单/撤单/改单；需要 API Key；建议先用模拟盘；必须显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 强烈建议：1=模拟盘
export OKX_CONFIRM=YES # 必填：防止误触真实交易
export OKX_WS_BATCH_OP="order" # 必填：order/cancel-order/amend-order
export OKX_WS_BATCH_ARGS='[{"instId":"BTC-USDT","tdMode":"cash","side":"buy","ordType":"market","sz":"0.001","clOrdId":"c1"}]' # 必填：JSON 数组
export OKX_TIMEOUT="10s" # 可选
go run ./examples/ws_private_trade_batch_ops
```

获取全量产品行情（market/tickers；默认 SPOT）：

```bash
export OKX_INST_TYPE="SPOT" # 可选：SPOT/SWAP/FUTURES/OPTION
export OKX_INST_FAMILY="" # 可选：仅衍生品/期权有效
go run ./examples/market_tickers
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

获取产品完整深度（market/books-full；默认 BTC-USDT，默认 20 档）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_BOOKS_SZ="20" # 可选：最大 5000（买卖深度共 10000 条）
go run ./examples/market_books_full
```

获取 SBE 订单簿快照（market/books-sbe；成功返回二进制，失败返回 JSON；默认 SPOT/BTC-USDT）：

```bash
export OKX_INST_TYPE="SPOT" # 可选：SPOT/SWAP/FUTURES/OPTION
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_SOURCE="0" # 可选：当前仅 0
go run ./examples/market_books_sbe
```

获取 K 线（默认 BTC-USDT，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
go run ./examples/market_candles
```

获取交易产品历史 K 线（market/history-candles；默认 BTC-USDT，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
go run ./examples/market_history_candles
```

获取指数行情（market/index-tickers；instId 与 quoteCcy 二选一）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_QUOTE_CCY="" # 可选：USD/USDT/BTC/USDC
go run ./examples/market_index_tickers
```

获取指数 K 线（market/index-candles；默认 BTC-USD，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USD" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
go run ./examples/market_index_candles
```

获取指数历史 K 线（market/history-index-candles；默认 BTC-USD，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USD" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
go run ./examples/market_history_index_candles
```

获取标记价格 K 线（market/mark-price-candles；默认 BTC-USD-240628，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USD-240628" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
go run ./examples/market_mark_price_candles
```

获取标记价格历史 K 线（market/history-mark-price-candles；默认 BTC-USD-240628，默认 1m，默认 5 根）：

```bash
export OKX_INST_ID="BTC-USD-240628" # 可选
export OKX_CANDLE_BAR="1m" # 可选
export OKX_CANDLE_LIMIT="5" # 可选
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
go run ./examples/market_history_mark_price_candles
```

获取最近成交（默认 BTC-USDT，默认 20 条）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_TRADES_LIMIT="20" # 可选
go run ./examples/market_trades
```

获取公共历史成交（market/history-trades；默认 BTC-USDT，默认 20 条）：

```bash
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_TYPE="" # 可选：1=tradeId 分页，2=时间戳分页
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
export OKX_TRADES_LIMIT="20" # 可选
go run ./examples/market_history_trades
```

查询产品信息（默认 SPOT）：

```bash
export OKX_INST_TYPE="SPOT" # 可选：SPOT/SWAP/FUTURES/OPTION
export OKX_ULY="BTC-USD" # 可选
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
go run ./examples/public_instruments
```

获取衍生品标的指数（public/underlying；默认 FUTURES）：

```bash
export OKX_INST_TYPE="FUTURES" # 可选：SWAP/FUTURES/OPTION
go run ./examples/public_underlying
```

获取预估交割/行权价格（public/estimated-price；交割/行权前一小时才有值；默认 BTC-USD-200214）：

```bash
export OKX_INST_ID="BTC-USD-200214" # 可选：仅适用于交割/期权
go run ./examples/public_estimated_price
```

查询限价（public/price-limit；默认 BTC-USDT-SWAP）：

```bash
export OKX_INST_ID="BTC-USDT-SWAP" # 可选
go run ./examples/public_price_limit
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

查询永续合约历史资金费率（public/funding-rate-history；默认 BTC-USD-SWAP）：

```bash
export OKX_INST_ID="BTC-USD-SWAP" # 可选
export OKX_AFTER="" # 可选：fundingTime（Unix 毫秒）
export OKX_BEFORE="" # 可选：fundingTime（Unix 毫秒）
export OKX_LIMIT="400" # 可选
go run ./examples/public_funding_rate_history
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

查询期权价格梯度（public/instrument-tick-bands；默认 OPTION）：

```bash
export OKX_INST_TYPE="OPTION" # 必填：当前仅支持 OPTION
export OKX_INST_FAMILY="BTC-USD" # 可选
go run ./examples/public_instrument_tick_bands
```

查询风险保证金余额（public/insurance-fund；默认 SWAP + BTC-USD）：

```bash
export OKX_INST_TYPE="SWAP" # 必填：MARGIN/SWAP/FUTURES/OPTION
export OKX_TYPE="" # 可选：regular_update/liquidation_balance_deposit/bankruptcy_loss/platform_revenue/adl
export OKX_INST_FAMILY="" # 可选（instType=SWAP/FUTURES/OPTION 时，instFamily/uly 至少传一个）
export OKX_ULY="" # 可选（instType=SWAP/FUTURES/OPTION 时，instFamily/uly 至少传一个）
export OKX_CCY="" # 可选（instType=MARGIN 时必填）
export OKX_AFTER="" # 可选
export OKX_BEFORE="" # 可选
export OKX_LIMIT="100" # 可选
go run ./examples/public_insurance_fund
```

张/币转换（public/convert-contract-coin；默认 BTC-USD-SWAP）：

```bash
export OKX_INST_ID="BTC-USD-SWAP" # 可选
export OKX_SZ="0.888" # 必填：币转张时为币数量；张转币时为张数量
export OKX_CONVERT_TYPE="1" # 可选：1=币转张，2=张转币
export OKX_PX="35000" # 可选（某些场景必填，按 OKX 规则）
export OKX_UNIT="coin" # 可选：coin/usds
export OKX_OP_TYPE="close" # 可选：open/close
go run ./examples/public_convert_contract_coin
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

查询历史订单（近三个月，需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SPOT" # 可选（示例默认 SPOT）
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_ORD_TYPE="limit" # 可选（可用逗号分隔多个）
export OKX_ORDER_STATE="filled" # 可选
export OKX_CATEGORY="twap" # 可选
export OKX_BEGIN="1695190491421" # 可选：Unix 毫秒时间戳
export OKX_END="1695190491421" # 可选：Unix 毫秒时间戳
export OKX_LIMIT="100" # 可选
go run ./examples/trade_orders_history_archive
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

查询成交明细（近三个月，需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_TYPE="SPOT" # 可选（示例默认 SPOT）
export OKX_INST_FAMILY="BTC-USD" # 可选
export OKX_INST_ID="BTC-USDT" # 可选
export OKX_ORD_ID="123" # 可选
export OKX_SUB_TYPE="1" # 可选
export OKX_BEGIN="1695190491421" # 可选：Unix 毫秒时间戳
export OKX_END="1695190491421" # 可选：Unix 毫秒时间戳
export OKX_LIMIT="100" # 可选
go run ./examples/trade_fills_history
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

查询账户限速信息（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
go run ./examples/trade_account_rate_limit
```

查询小币一键兑换主流币币种列表（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_SOURCE="1" # 可选：1=交易账户，2=资金账户
go run ./examples/trade_easy_convert_currency_list
```

查询小币一键兑换主流币历史记录（需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_AFTER="" # 可选：Unix 毫秒时间戳
export OKX_BEFORE="" # 可选：Unix 毫秒时间戳
export OKX_LIMIT="100" # 可选
go run ./examples/trade_easy_convert_history
```

小币一键兑换主流币交易（会实际交易；需要 API Key；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_FROM_CCY_LIST='["ADA","USDC"]' # 必填：JSON array（最多 5 个）
export OKX_TO_CCY="OKB" # 必填
export OKX_SOURCE="1" # 可选：1=交易账户，2=资金账户
export OKX_CONFIRM="YES"
go run ./examples/trade_easy_convert
```

查询一键还债币种列表（新，仅适用于现货模式；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
go run ./examples/trade_one_click_repay_currency_list_v2
```

查询一键还债历史记录（新，仅适用于现货模式；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_AFTER="" # 可选：Unix 毫秒时间戳
export OKX_BEFORE="" # 可选：Unix 毫秒时间戳
export OKX_LIMIT="100" # 可选
go run ./examples/trade_one_click_repay_history_v2
```

一键还债交易（新，仅适用于现货模式；会实际交易；需要 API Key；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEBT_CCY="USDC" # 必填
export OKX_REPAY_CCY_LIST='["USDC","BTC"]' # 必填：JSON array（排序代表偿还优先级）
export OKX_CONFIRM="YES"
go run ./examples/trade_one_click_repay_v2
```

查询一键还债币种列表（跨币种保证金/组合保证金；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_DEBT_TYPE="" # 可选：cross/isolated
go run ./examples/trade_one_click_repay_currency_list
```

查询一键还债历史记录（跨币种保证金/组合保证金；需要 API Key，支持模拟盘）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_AFTER="" # 可选：Unix 毫秒时间戳
export OKX_BEFORE="" # 可选：Unix 毫秒时间戳
export OKX_LIMIT="100" # 可选
go run ./examples/trade_one_click_repay_history
```

一键还债交易（跨币种保证金/组合保证金；会实际交易；需要 API Key；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEBT_CCY_LIST='["ETH","BTC"]' # 必填：JSON array（最多 5 个）
export OKX_REPAY_CCY="USDT" # 必填
export OKX_CONFIRM="YES"
go run ./examples/trade_one_click_repay
```

市价仓位全平（会实际平仓；需要 API Key，支持模拟盘；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_INST_ID="BTC-USDT-SWAP" # 必填
export OKX_MGN_MODE="cross" # 必填：cross/isolated
export OKX_POS_SIDE="" # 可选：long/short/net
export OKX_CCY="" # 可选
export OKX_AUTO_CXL="true" # 可选：true/false
export OKX_CL_ORD_ID="" # 可选
export OKX_TAG="" # 可选
export OKX_CONFIRM="YES"
go run ./examples/trade_close_position
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

设置倒计时全部撤单（会实际撤单；需要 API Key，支持模拟盘；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1 # 可选：1=模拟盘
export OKX_TIME_OUT="60" # 必填：0 或 10-120
export OKX_TAG="my-bot" # 可选
export OKX_CONFIRM="YES"
go run ./examples/trade_cancel_all_after
```

撤销 MMP 订单（会实际撤单；需要 API Key；仅适用于组合保证金期权账户且有 MMP 权限；需显式确认）：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_INST_FAMILY="BTC-USD" # 必填
export OKX_LOCK_INTERVAL="0" # 可选：0-10000（毫秒）
export OKX_CONFIRM="YES"
go run ./examples/trade_mass_cancel
```

## 设计文档

- `docs/design.md`
- `docs/roadmap.md`
