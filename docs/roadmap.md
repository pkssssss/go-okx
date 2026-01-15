# v0.1 路线图（P0-P3）

> 目标：简洁 / 高效 / 稳定；API 体验对标 `go-binance` 的 `Service + Do(ctx)` 风格。  
> Go Module：`github.com/pkssssss/go-okx/v5`（规则 A：对齐 OKX API v5）。

## P0（基础骨架）✅

- 目录结构：根目录 + `v5/` 主模块 + `examples/` 独立 module + `go.work`
- REST 通用管线：`requestPath` 单一来源、签名、envelope 解包、统一 `APIError`
- 时间校准：`SyncTime` + `TimeOffset`
- 工程自检：CI + `./check.sh`（`gofmt`/`vet`/`test`/`race`/`examples`）

## P1（行情 & WS Public）✅

- REST 公共/行情：`public_time`、`public_instruments`、`public/block-trades`、`public/discount-rate-interest-free-quota/interest-rate-loan-quota/economic-calendar/market-data-history/option-trades`、`funding/funding-rate-history/open-interest/mark-price/price-limit/opt-summary/instrument-tick-bands/insurance-fund/convert-contract-coin/underlying/estimated-price/estimated-settlement-info/settlement-history/delivery-exercise-history/premium-history/position-tiers`、`market_(books/books-full/books-sbe/candles/tickers/trades/index-tickers)`
- WS public：`books/books5/tickers/trades` + 关键频道解析（`open-interest`/`funding-rate`/`mark-price`/`index-tickers`/`price-limit`/`opt-summary`/`liquidation-orders`）
- WS business：`mark-price-candle*`、`index-candle*`（按需）

## P2（交易主链路 & 资产/账单）✅

- REST 交易主链路：下单/改单/撤单/批量 + 查单 + 历史委托 + 成交
- WS private：orders/fills/positions/account/balance_and_position（含自动重连/重订阅与异步 handler 方案）
- REST 资产/账单：`account/bills`、`account/bills-archive`、`asset/balances`、`asset/transfer`、`asset/transfer-state`

## P3（扩展 & 工程化）✅

- Market：
  - ✅ 历史行情：`market/history-candles`、`market/history-trades`
  - ✅ 指数K线：`market/index-candles`、`market/history-index-candles`
  - ✅ 标记价格K线：`market/mark-price-candles`、`market/history-mark-price-candles`

- Asset：
  - ✅ 估值/资产：`asset/asset-valuation`、`asset/non-tradable-assets`
  - ✅ 币种/辅助：`asset/currencies`（充提状态/手续费/精度）、`asset/exchange-list`（交易所列表）
  - ✅ 资金流水：`asset/bills`、`asset/bills-history`；`account/bills`、`account/bills-archive`
  - ✅ 划转：`asset/transfer`、`asset/transfer-state`；子账户：`asset/subaccount/balances`、`asset/subaccount/bills`、`asset/subaccount/managed-subaccount-bills`、`asset/subaccount/transfer`
  - ✅ 充提链路：`asset/deposit-address`、`asset/deposit-history`、`asset/withdrawal`、`asset/cancel-withdrawal`、`asset/withdrawal-history`、`asset/deposit-withdraw-status`
  - ✅ 月结单：`asset/monthly-statement`（apply/get）
  - ✅ 闪兑：`asset/convert/currencies`、`asset/convert/currency-pair`、`asset/convert/estimate-quote`、`asset/convert/trade`、`asset/convert/history`
  - ✅ WS（business）：`deposit-info`、`withdrawal-info`
- Account：✅ `account/config`、✅ `account/instruments`、✅ `account/adjust-leverage-info`、✅ `account/greeks`、✅ `account/set-greeks`、✅ `account/set-riskOffset-amt`、✅ `account/set-fee-type`、✅ `account/set-isolated-mode`、✅ `account/set-auto-earn`、✅ `account/set-settle-currency`、✅ `account/set-trading-config`、✅ `account/activate-option`、✅ `account/precheck-set-delta-neutral`、✅ `account/bills-history-archive`（apply/get）、✅ `account/set-position-mode`、✅ `account/set-leverage`、✅ `account/leverage-info`、✅ `account/max-size`、✅ `account/max-avail-size`、✅ `account/max-loan`、✅ `account/trade-fee`、✅ `account/interest-accrued`、✅ `account/interest-rate`、✅ `account/max-withdrawal`、✅ `account/subaccount/balances`、✅ `account/subaccount/max-withdrawal`、✅ `account/interest-limits`、✅ `account/position/margin-balance`、✅ `account/spot-manual-borrow-repay`、✅ `account/risk-state`、✅ `account/set-auto-repay`、✅ `account/set-auto-loan`、✅ `account/account-level-switch-preset`、✅ `account/set-account-switch-precheck`、✅ `account/set-account-level`、✅ `account/set-collateral-assets`、✅ `account/collateral-assets`、✅ `account/mmp-reset`、✅ `account/mmp-config`、✅ `account/move-positions`、✅ `account/move-positions-history`、✅ `account/spot-borrow-repay-history`、✅ `account/positions-history`、✅ `account/account-position-risk`、✅ `account/position-tiers`、✅ `account/position-builder`、✅ `account/position-builder-graph`
- Users：✅ `users/subaccount/list`、✅ `users/subaccount/create-subaccount`、✅ `users/subaccount/apikey`（create/query）、✅ `users/subaccount/modify-apikey`、✅ `users/subaccount/delete-apikey`、✅ `users/subaccount/set-transfer-out`、✅ `users/entrust-subaccount-list`
- Trade（策略委托）：✅ `trade/order-precheck`、✅ `trade/order-algo`（place/get）、✅ `trade/cancel-algos`、✅ `trade/amend-algos`、✅ `trade/orders-algo-pending`、✅ `trade/orders-algo-history`
- Trade（平仓）：✅ `trade/close-position`
- Trade（风控/限速）：✅ `trade/cancel-all-after`、✅ `trade/account-rate-limit`
- Trade（MMP）：✅ `trade/mass-cancel`
- Trade（一键还债 v2）：✅ `trade/one-click-repay-currency-list-v2`、✅ `trade/one-click-repay-v2`、✅ `trade/one-click-repay-history-v2`
- Trade（一键还债）：✅ `trade/one-click-repay-currency-list`、✅ `trade/one-click-repay`、✅ `trade/one-click-repay-history`
- Trade（一键兑换）：✅ `trade/easy-convert-currency-list`、✅ `trade/easy-convert`、✅ `trade/easy-convert-history`
- 工程化：
  - ✅ 错误分类：`IsAuthError` / `IsRateLimitError` / `IsTimeSkewError`
  - ✅ 可控重试（仅幂等 GET）：`WithRetry(RetryConfig{...})`
  - ✅ README/docs：补齐高频行情 examples 运行指引

## P4（TradingBot 交易机器人）✅

- ✅ REST：Grid / Recurring / Signal（含 public 端点）
- ✅ 单测：聚合用例 `v5/trading_bot_services_test.go`
- ✅ Examples：已补齐（交易/变更类示例强制 `OKX_CONFIRM=YES` 防误触）
