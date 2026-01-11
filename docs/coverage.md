# 覆盖矩阵（OKX API v5）

> 主维度：按 OKX 官方 v5 路径前缀分组（REST/WS）。  
> 辅维度：用标签做“场景索引”（便于从需求反查能力）。

## 约定

- 状态符号：✅ 已实现；🟡 部分实现/有已知限制；❌ 未实现（按需增量）。
- 鉴权：`public`=无需签名；`private`=需要签名/凭证。
- 精度：金额/数量/费率/价格默认保持为 `string`（无损）。
- 说明：OKX 多数能力通过参数（如 `instType`/`instId`/`instFamily`）区分现货/合约/期权，本表不重复标注。

## REST

### Public（公共数据）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/public/block-trades` | `public` | [`public_block_trades_service.go`](../v5/public_block_trades_service.go) | ✅ [`public_block_trades_service_test.go`](../v5/public_block_trades_service_test.go) | [`examples/public_block_trades`](../examples/public_block_trades) | `rest, auth-public, public` |
| `GET /api/v5/public/convert-contract-coin` | `public` | [`public_convert_contract_coin_service.go`](../v5/public_convert_contract_coin_service.go) | ✅ [`public_convert_contract_coin_service_test.go`](../v5/public_convert_contract_coin_service_test.go) | [`examples/public_convert_contract_coin`](../examples/public_convert_contract_coin) | `rest, auth-public, public, convert` |
| `GET /api/v5/public/delivery-exercise-history` | `public` | [`public_delivery_exercise_history_service.go`](../v5/public_delivery_exercise_history_service.go) | ✅ [`public_delivery_exercise_history_service_test.go`](../v5/public_delivery_exercise_history_service_test.go) | [`examples/public_delivery_exercise_history`](../examples/public_delivery_exercise_history) | `rest, auth-public, public` |
| `GET /api/v5/public/discount-rate-interest-free-quota` | `public` | [`public_discount_rate_interest_free_quota_service.go`](../v5/public_discount_rate_interest_free_quota_service.go) | ✅ [`public_discount_rate_interest_free_quota_service_test.go`](../v5/public_discount_rate_interest_free_quota_service_test.go) | [`examples/public_discount_rate_interest_free_quota`](../examples/public_discount_rate_interest_free_quota) | `rest, auth-public, public, loan` |
| `GET /api/v5/public/economic-calendar` | `private` | [`public_economic_calendar_service.go`](../v5/public_economic_calendar_service.go) | ✅ [`public_economic_calendar_service_test.go`](../v5/public_economic_calendar_service_test.go) | [`examples/public_economic_calendar`](../examples/public_economic_calendar) | `rest, auth-private, public` |
| `GET /api/v5/public/estimated-price` | `public` | [`public_estimated_price_service.go`](../v5/public_estimated_price_service.go) | ✅ [`public_estimated_price_service_test.go`](../v5/public_estimated_price_service_test.go) | [`examples/public_estimated_price`](../examples/public_estimated_price) | `rest, auth-public, public` |
| `GET /api/v5/public/estimated-settlement-info` | `public` | [`public_estimated_settlement_info_service.go`](../v5/public_estimated_settlement_info_service.go) | ✅ [`public_estimated_settlement_info_service_test.go`](../v5/public_estimated_settlement_info_service_test.go) | [`examples/public_estimated_settlement_info`](../examples/public_estimated_settlement_info) | `rest, auth-public, public` |
| `GET /api/v5/public/funding-rate` | `public` | [`public_funding_rate_service.go`](../v5/public_funding_rate_service.go) | ✅ [`public_funding_rate_service_test.go`](../v5/public_funding_rate_service_test.go) | [`examples/public_funding_rate`](../examples/public_funding_rate) | `rest, auth-public, public` |
| `GET /api/v5/public/funding-rate-history` | `public` | [`public_funding_rate_history_service.go`](../v5/public_funding_rate_history_service.go) | ✅ [`public_funding_rate_history_service_test.go`](../v5/public_funding_rate_history_service_test.go) | [`examples/public_funding_rate_history`](../examples/public_funding_rate_history) | `rest, auth-public, public` |
| `GET /api/v5/public/instrument-tick-bands` | `public` | [`public_instrument_tick_bands_service.go`](../v5/public_instrument_tick_bands_service.go) | ✅ [`public_instrument_tick_bands_service_test.go`](../v5/public_instrument_tick_bands_service_test.go) | [`examples/public_instrument_tick_bands`](../examples/public_instrument_tick_bands) | `rest, auth-public, public` |
| `GET /api/v5/public/instruments` | `public` | [`public_instruments_service.go`](../v5/public_instruments_service.go) | ✅ [`public_instruments_service_test.go`](../v5/public_instruments_service_test.go) | [`examples/public_instruments`](../examples/public_instruments) | `rest, auth-public, public` |
| `GET /api/v5/public/insurance-fund` | `public` | [`public_insurance_fund_service.go`](../v5/public_insurance_fund_service.go) | ✅ [`public_insurance_fund_service_test.go`](../v5/public_insurance_fund_service_test.go) | [`examples/public_insurance_fund`](../examples/public_insurance_fund) | `rest, auth-public, public` |
| `GET /api/v5/public/interest-rate-loan-quota` | `public` | [`public_interest_rate_loan_quota_service.go`](../v5/public_interest_rate_loan_quota_service.go) | ✅ [`public_interest_rate_loan_quota_service_test.go`](../v5/public_interest_rate_loan_quota_service_test.go) | [`examples/public_interest_rate_loan_quota`](../examples/public_interest_rate_loan_quota) | `rest, auth-public, public, loan` |
| `GET /api/v5/public/mark-price` | `public` | [`public_mark_price_service.go`](../v5/public_mark_price_service.go) | ✅ [`public_mark_price_service_test.go`](../v5/public_mark_price_service_test.go) | [`examples/public_mark_price`](../examples/public_mark_price) | `rest, auth-public, public` |
| `GET /api/v5/public/market-data-history` | `public` | [`public_market_data_history_service.go`](../v5/public_market_data_history_service.go) | ✅ [`public_market_data_history_service_test.go`](../v5/public_market_data_history_service_test.go) | [`examples/public_market_data_history`](../examples/public_market_data_history) | `rest, auth-public, public` |
| `GET /api/v5/public/open-interest` | `public` | [`public_open_interest_service.go`](../v5/public_open_interest_service.go) | ✅ [`public_open_interest_service_test.go`](../v5/public_open_interest_service_test.go) | [`examples/public_open_interest`](../examples/public_open_interest) | `rest, auth-public, public, open-interest` |
| `GET /api/v5/public/opt-summary` | `public` | [`public_opt_summary_service.go`](../v5/public_opt_summary_service.go) | ✅ [`public_opt_summary_service_test.go`](../v5/public_opt_summary_service_test.go) | [`examples/public_opt_summary`](../examples/public_opt_summary) | `rest, auth-public, public, option` |
| `GET /api/v5/public/option-trades` | `public` | [`public_option_trades_service.go`](../v5/public_option_trades_service.go) | ✅ [`public_option_trades_service_test.go`](../v5/public_option_trades_service_test.go) | [`examples/public_option_trades`](../examples/public_option_trades) | `rest, auth-public, public, option` |
| `GET /api/v5/public/position-tiers` | `public` | [`public_position_tiers_service.go`](../v5/public_position_tiers_service.go) | ✅ [`public_position_tiers_service_test.go`](../v5/public_position_tiers_service_test.go) | [`examples/public_position_tiers`](../examples/public_position_tiers) | `rest, auth-public, public` |
| `GET /api/v5/public/premium-history` | `public` | [`public_premium_history_service.go`](../v5/public_premium_history_service.go) | ✅ [`public_premium_history_service_test.go`](../v5/public_premium_history_service_test.go) | [`examples/public_premium_history`](../examples/public_premium_history) | `rest, auth-public, public` |
| `GET /api/v5/public/price-limit` | `public` | [`public_price_limit_service.go`](../v5/public_price_limit_service.go) | ✅ [`public_price_limit_service_test.go`](../v5/public_price_limit_service_test.go) | [`examples/public_price_limit`](../examples/public_price_limit) | `rest, auth-public, public` |
| `GET /api/v5/public/settlement-history` | `public` | [`public_settlement_history_service.go`](../v5/public_settlement_history_service.go) | ✅ [`public_settlement_history_service_test.go`](../v5/public_settlement_history_service_test.go) | [`examples/public_settlement_history`](../examples/public_settlement_history) | `rest, auth-public, public` |
| `GET /api/v5/public/time` | `public` | [`public_time_service.go`](../v5/public_time_service.go) | ✅ [`public_time_service_test.go`](../v5/public_time_service_test.go) | [`examples/public_time`](../examples/public_time) | `rest, auth-public, public` |
| `GET /api/v5/public/underlying` | `public` | [`public_underlying_service.go`](../v5/public_underlying_service.go) | ✅ [`public_underlying_service_test.go`](../v5/public_underlying_service_test.go) | [`examples/public_underlying`](../examples/public_underlying) | `rest, auth-public, public` |

### Market（行情）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/market/books` | `public` | [`market_books_service.go`](../v5/market_books_service.go) | ✅ [`market_books_service_test.go`](../v5/market_books_service_test.go) | [`examples/market_books`](../examples/market_books) | `rest, auth-public, market, orderbook` |
| `GET /api/v5/market/books-full` | `public` | [`market_books_full_service.go`](../v5/market_books_full_service.go) | ✅ [`market_books_full_service_test.go`](../v5/market_books_full_service_test.go) | [`examples/market_books_full`](../examples/market_books_full) | `rest, auth-public, market, orderbook` |
| `GET /api/v5/market/candles` | `public` | [`market_candles_service.go`](../v5/market_candles_service.go) | ✅ [`market_candles_service_test.go`](../v5/market_candles_service_test.go) | [`examples/market_candles`](../examples/market_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/history-candles` | `public` | [`market_history_candles_service.go`](../v5/market_history_candles_service.go) | ✅ [`market_history_candles_service_test.go`](../v5/market_history_candles_service_test.go) | [`examples/market_history_candles`](../examples/market_history_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/history-index-candles` | `public` | [`market_history_index_candles_service.go`](../v5/market_history_index_candles_service.go) | ✅ [`market_history_index_candles_service_test.go`](../v5/market_history_index_candles_service_test.go) | [`examples/market_history_index_candles`](../examples/market_history_index_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/history-mark-price-candles` | `public` | [`market_history_mark_price_candles_service.go`](../v5/market_history_mark_price_candles_service.go) | ✅ [`market_history_mark_price_candles_service_test.go`](../v5/market_history_mark_price_candles_service_test.go) | [`examples/market_history_mark_price_candles`](../examples/market_history_mark_price_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/history-trades` | `public` | [`market_history_trades_service.go`](../v5/market_history_trades_service.go) | ✅ [`market_history_trades_service_test.go`](../v5/market_history_trades_service_test.go) | [`examples/market_history_trades`](../examples/market_history_trades) | `rest, auth-public, market` |
| `GET /api/v5/market/index-candles` | `public` | [`market_index_candles_service.go`](../v5/market_index_candles_service.go) | ✅ [`market_index_candles_service_test.go`](../v5/market_index_candles_service_test.go) | [`examples/market_index_candles`](../examples/market_index_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/index-tickers` | `public` | [`market_index_tickers_service.go`](../v5/market_index_tickers_service.go) | ✅ [`market_index_tickers_service_test.go`](../v5/market_index_tickers_service_test.go) | [`examples/market_index_tickers`](../examples/market_index_tickers) | `rest, auth-public, market, tickers` |
| `GET /api/v5/market/mark-price-candles` | `public` | [`market_mark_price_candles_service.go`](../v5/market_mark_price_candles_service.go) | ✅ [`market_mark_price_candles_service_test.go`](../v5/market_mark_price_candles_service_test.go) | [`examples/market_mark_price_candles`](../examples/market_mark_price_candles) | `rest, auth-public, market, candles` |
| `GET /api/v5/market/ticker` | `public` | [`market_ticker_service.go`](../v5/market_ticker_service.go) | ✅（聚合） [`market_services_test.go`](../v5/market_services_test.go) | [`examples/market_ticker`](../examples/market_ticker) | `rest, auth-public, market, tickers` |
| `GET /api/v5/market/tickers` | `public` | [`market_tickers_service.go`](../v5/market_tickers_service.go) | ✅（聚合） [`market_services_test.go`](../v5/market_services_test.go) | [`examples/market_tickers`](../examples/market_tickers) | `rest, auth-public, market, tickers` |
| `GET /api/v5/market/trades` | `public` | [`market_trades_service.go`](../v5/market_trades_service.go) | ✅ [`market_trades_service_test.go`](../v5/market_trades_service_test.go) | [`examples/market_trades`](../examples/market_trades) | `rest, auth-public, market, trades` |

### Trade（交易）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/trade/account-rate-limit` | `private` | [`trade_account_rate_limit`](../v5/trade_account_rate_limit_service.go) | ✅ [test](../v5/trade_account_rate_limit_service_test.go) | [ex](../examples/trade_account_rate_limit) | `auth-private, risk` |
| `POST /api/v5/trade/amend-algos` | `private` | [`trade_amend_algo_order`](../v5/trade_amend_algo_order_service.go) | ✅ [test](../v5/trade_amend_algo_order_service_test.go) | — | `auth-private, algo` |
| `POST /api/v5/trade/amend-batch-orders` | `private` | [`trade_batch_amend_orders`](../v5/trade_batch_amend_orders_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `POST /api/v5/trade/amend-order` | `private` | [`trade_amend_order`](../v5/trade_amend_order_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `POST /api/v5/trade/batch-orders` | `private` | [`trade_batch_place_orders`](../v5/trade_batch_place_orders_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `POST /api/v5/trade/cancel-algos` | `private` | [`trade_cancel_algo_orders`](../v5/trade_cancel_algo_orders_service.go) | ✅ [test](../v5/trade_cancel_algo_orders_service_test.go) | — | `auth-private, algo` |
| `POST /api/v5/trade/cancel-all-after` | `private` | [`trade_cancel_all_after`](../v5/trade_cancel_all_after_service.go) | ✅ [test](../v5/trade_cancel_all_after_service_test.go) | [ex](../examples/trade_cancel_all_after) | `auth-private, risk` |
| `POST /api/v5/trade/cancel-batch-orders` | `private` | [`trade_batch_cancel_orders`](../v5/trade_batch_cancel_orders_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `POST /api/v5/trade/cancel-order` | `private` | [`trade_cancel_order`](../v5/trade_cancel_order_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `POST /api/v5/trade/close-position` | `private` | [`trade_close_positions`](../v5/trade_close_positions_service.go) | ✅ [test](../v5/trade_close_positions_service_test.go) | — | `auth-private` |
| `POST /api/v5/trade/easy-convert` | `private` | [`trade_easy_convert`](../v5/trade_easy_convert_service.go) | ✅ [test](../v5/trade_easy_convert_service_test.go) | [ex](../examples/trade_easy_convert) | `auth-private, convert` |
| `GET /api/v5/trade/easy-convert-currency-list` | `private` | [`trade_easy_convert_currency_list`](../v5/trade_easy_convert_currency_list_service.go) | ✅ [test](../v5/trade_easy_convert_currency_list_service_test.go) | [ex](../examples/trade_easy_convert_currency_list) | `auth-private, convert` |
| `GET /api/v5/trade/easy-convert-history` | `private` | [`trade_easy_convert_history`](../v5/trade_easy_convert_history_service.go) | ✅ [test](../v5/trade_easy_convert_history_service_test.go) | [ex](../examples/trade_easy_convert_history) | `auth-private, convert` |
| `GET /api/v5/trade/fills` | `private` | [`trade_fills`](../v5/trade_fills_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | [ex](../examples/trade_fills) | `auth-private, fills` |
| `GET /api/v5/trade/fills-history` | `private` | [`trade_fills_history`](../v5/trade_fills_history_service.go) | ✅ [test](../v5/trade_fills_history_service_test.go) | [ex](../examples/trade_fills_history) | `auth-private, fills` |
| `POST /api/v5/trade/mass-cancel` | `private` | [`trade_mass_cancel`](../v5/trade_mass_cancel_service.go) | ✅ [test](../v5/trade_mass_cancel_service_test.go) | [ex](../examples/trade_mass_cancel) | `auth-private` |
| `POST /api/v5/trade/one-click-repay` | `private` | [`trade_one_click_repay`](../v5/trade_one_click_repay_service.go) | ✅ [test](../v5/trade_one_click_repay_service_test.go) | [ex](../examples/trade_one_click_repay) | `auth-private, loan` |
| `GET /api/v5/trade/one-click-repay-currency-list` | `private` | [`trade_one_click_repay_currency_list`](../v5/trade_one_click_repay_currency_list_service.go) | ✅ [test](../v5/trade_one_click_repay_currency_list_service_test.go) | [ex](../examples/trade_one_click_repay_currency_list) | `auth-private, loan` |
| `GET /api/v5/trade/one-click-repay-currency-list-v2` | `private` | [`trade_one_click_repay_currency_list_v2`](../v5/trade_one_click_repay_currency_list_v2_service.go) | ✅ [test](../v5/trade_one_click_repay_currency_list_v2_service_test.go) | [ex](../examples/trade_one_click_repay_currency_list_v2) | `auth-private, loan` |
| `GET /api/v5/trade/one-click-repay-history` | `private` | [`trade_one_click_repay_history`](../v5/trade_one_click_repay_history_service.go) | ✅ [test](../v5/trade_one_click_repay_history_service_test.go) | [ex](../examples/trade_one_click_repay_history) | `auth-private, loan` |
| `GET /api/v5/trade/one-click-repay-history-v2` | `private` | [`trade_one_click_repay_history_v2`](../v5/trade_one_click_repay_history_v2_service.go) | ✅ [test](../v5/trade_one_click_repay_history_v2_service_test.go) | [ex](../examples/trade_one_click_repay_history_v2) | `auth-private, loan` |
| `POST /api/v5/trade/one-click-repay-v2` | `private` | [`trade_one_click_repay_v2`](../v5/trade_one_click_repay_v2_service.go) | ✅ [test](../v5/trade_one_click_repay_v2_service_test.go) | [ex](../examples/trade_one_click_repay_v2) | `auth-private, loan` |
| `GET /api/v5/trade/order` | `private` | [`trade_get_order`](../v5/trade_get_order_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | [ex](../examples/trade_get_order) | `auth-private, orders` |
| `POST /api/v5/trade/order` | `private` | [`trade_place_order`](../v5/trade_place_order_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | — | `auth-private, orders` |
| `GET /api/v5/trade/order-algo` | `private` | [`trade_get_algo_order`](../v5/trade_get_algo_order_service.go) | ✅ [test](../v5/trade_get_algo_order_service_test.go) | [ex](../examples/trade_get_algo_order) | `auth-private, orders, algo` |
| `POST /api/v5/trade/order-algo` | `private` | [`trade_place_algo_order`](../v5/trade_place_algo_order_service.go) | ✅ [test](../v5/trade_place_algo_order_service_test.go) | [ex](../examples/trade_place_algo_order) | `auth-private, orders, algo` |
| `POST /api/v5/trade/order-precheck` | `private` | [`trade_order_precheck`](../v5/trade_order_precheck_service.go) | ✅ [test](../v5/trade_order_precheck_service_test.go) | [ex](../examples/trade_order_precheck) | `auth-private, orders` |
| `GET /api/v5/trade/orders-algo-history` | `private` | [`trade_orders_algo_history`](../v5/trade_orders_algo_history_service.go) | ✅ [test](../v5/trade_orders_algo_history_service_test.go) | [ex](../examples/trade_orders_algo_history) | `auth-private, orders, algo` |
| `GET /api/v5/trade/orders-algo-pending` | `private` | [`trade_orders_algo_pending`](../v5/trade_orders_algo_pending_service.go) | ✅ [test](../v5/trade_orders_algo_pending_service_test.go) | [ex](../examples/trade_orders_algo_pending) | `auth-private, orders, algo` |
| `GET /api/v5/trade/orders-history` | `private` | [`trade_orders_history`](../v5/trade_orders_history_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | [ex](../examples/trade_orders_history) | `auth-private, orders` |
| `GET /api/v5/trade/orders-history-archive` | `private` | [`trade_orders_history_archive`](../v5/trade_orders_history_archive_service.go) | ✅ [test](../v5/trade_orders_history_archive_service_test.go) | [ex](../examples/trade_orders_history_archive) | `auth-private, orders` |
| `GET /api/v5/trade/orders-pending` | `private` | [`trade_orders_pending`](../v5/trade_orders_pending_service.go) | ✅聚合 [test](../v5/trade_services_test.go) | [ex](../examples/trade_orders_pending) | `auth-private, orders` |

### Account（账户）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `POST /api/v5/account/account-level-switch-preset` | `private` | [`account_level_switch_preset`](../v5/account_level_switch_preset_service.go) | ✅ [test](../v5/account_level_switch_preset_service_test.go) | [ex](../examples/account_level_switch_preset) | `auth-private` |
| `GET /api/v5/account/account-position-risk` | `private` | [`account_position_risk`](../v5/account_position_risk_service.go) | ✅ [test](../v5/account_position_risk_service_test.go) | [ex](../examples/account_position_risk) | `auth-private, risk` |
| `POST /api/v5/account/activate-option` | `private` | [`account_activate_option`](../v5/account_activate_option_service.go) | ✅ [test](../v5/account_activate_option_service_test.go) | [ex](../examples/account_activate_option) | `auth-private, option` |
| `GET /api/v5/account/adjust-leverage-info` | `private` | [`account_adjust_leverage_info`](../v5/account_adjust_leverage_info_service.go) | ✅ [test](../v5/account_adjust_leverage_info_service_test.go) | [ex](../examples/account_adjust_leverage_info) | `auth-private, leverage` |
| `GET /api/v5/account/balance` | `private` | [`account_balance`](../v5/account_balance_service.go) | ✅ [test](../v5/account_balance_service_test.go) | [ex](../examples/account_balance) | `auth-private, balance` |
| `GET /api/v5/account/bills` | `private` | [`account_bills`](../v5/account_bills_service.go) | ✅ [test](../v5/account_bills_service_test.go) | — | `auth-private, bills` |
| `GET /api/v5/account/bills-archive` | `private` | [`account_bills_archive`](../v5/account_bills_archive_service.go) | ✅ [test](../v5/account_bills_archive_service_test.go) | — | `auth-private, bills` |
| `GET /api/v5/account/bills-history-archive` | `private` | [`account_bills_history_archive`](../v5/account_bills_history_archive_service.go) | ✅ [test](../v5/account_bills_history_archive_service_test.go) | [ex](../examples/account_bills_history_archive) | `auth-private, bills` |
| `POST /api/v5/account/bills-history-archive` | `private` | [`account_bills_history_archive_apply`](../v5/account_bills_history_archive_apply_service.go) | ✅ [test](../v5/account_bills_history_archive_apply_service_test.go) | [ex](../examples/account_bills_history_archive_apply) | `auth-private, bills` |
| `GET /api/v5/account/collateral-assets` | `private` | [`account_collateral_assets`](../v5/account_collateral_assets_service.go) | ✅ [test](../v5/account_collateral_assets_service_test.go) | [ex](../examples/account_collateral_assets) | `auth-private` |
| `GET /api/v5/account/config` | `private` | [`account_config`](../v5/account_config_service.go) | ✅ [test](../v5/account_config_service_test.go) | [ex](../examples/account_config) | `auth-private` |
| `GET /api/v5/account/greeks` | `private` | [`account_greeks`](../v5/account_greeks_service.go) | ✅ [test](../v5/account_greeks_service_test.go) | [ex](../examples/account_greeks) | `auth-private, option` |
| `GET /api/v5/account/instruments` | `private` | [`account_instruments`](../v5/account_instruments_service.go) | ✅ [test](../v5/account_instruments_service_test.go) | [ex](../examples/account_instruments) | `auth-private` |
| `GET /api/v5/account/interest-accrued` | `private` | [`account_interest_accrued`](../v5/account_interest_accrued_service.go) | ✅ [test](../v5/account_interest_accrued_service_test.go) | [ex](../examples/account_interest_accrued) | `auth-private, loan` |
| `GET /api/v5/account/interest-limits` | `private` | [`account_interest_limits`](../v5/account_interest_limits_service.go) | ✅ [test](../v5/account_interest_limits_service_test.go) | [ex](../examples/account_interest_limits) | `auth-private, loan` |
| `GET /api/v5/account/interest-rate` | `private` | [`account_interest_rate`](../v5/account_interest_rate_service.go) | ✅ [test](../v5/account_interest_rate_service_test.go) | [ex](../examples/account_interest_rate) | `auth-private, loan` |
| `GET /api/v5/account/leverage-info` | `private` | [`account_leverage_info`](../v5/account_leverage_info_service.go) | ✅ [test](../v5/account_leverage_info_service_test.go) | [ex](../examples/account_leverage_info) | `auth-private, leverage` |
| `GET /api/v5/account/max-avail-size` | `private` | [`account_max_avail_size`](../v5/account_max_avail_size_service.go) | ✅ [test](../v5/account_max_avail_size_service_test.go) | [ex](../examples/account_max_avail_size) | `auth-private` |
| `GET /api/v5/account/max-loan` | `private` | [`account_max_loan`](../v5/account_max_loan_service.go) | ✅ [test](../v5/account_max_loan_service_test.go) | [ex](../examples/account_max_loan) | `auth-private, loan` |
| `GET /api/v5/account/max-size` | `private` | [`account_max_size`](../v5/account_max_size_service.go) | ✅ [test](../v5/account_max_size_service_test.go) | [ex](../examples/account_max_size) | `auth-private` |
| `GET /api/v5/account/max-withdrawal` | `private` | [`account_max_withdrawal`](../v5/account_max_withdrawal_service.go) | ✅ [test](../v5/account_max_withdrawal_service_test.go) | [ex](../examples/account_max_withdrawal) | `auth-private` |
| `GET /api/v5/account/mmp-config` | `private` | [`account_mmp_config`](../v5/account_mmp_config_service.go) | ✅ [test](../v5/account_mmp_config_service_test.go) | [ex](../examples/account_mmp_config) | `auth-private, mmp` |
| `POST /api/v5/account/mmp-config` | `private` | [`account_set_mmp_config`](../v5/account_set_mmp_config_service.go) | ✅ [test](../v5/account_set_mmp_config_service_test.go) | [ex](../examples/account_set_mmp_config) | `auth-private, mmp` |
| `POST /api/v5/account/mmp-reset` | `private` | [`account_mmp_reset`](../v5/account_mmp_reset_service.go) | ✅ [test](../v5/account_mmp_reset_service_test.go) | [ex](../examples/account_mmp_reset) | `auth-private, mmp` |
| `POST /api/v5/account/move-positions` | `private` | [`account_move_positions`](../v5/account_move_positions_service.go) | ✅ [test](../v5/account_move_positions_service_test.go) | [ex](../examples/account_move_positions) | `auth-private, positions` |
| `GET /api/v5/account/move-positions-history` | `private` | [`account_move_positions_history`](../v5/account_move_positions_history_service.go) | ✅ [test](../v5/account_move_positions_history_service_test.go) | [ex](../examples/account_move_positions_history) | `auth-private, positions` |
| `POST /api/v5/account/position-builder` | `private` | [`account_position_builder`](../v5/account_position_builder_service.go) | ✅ [test](../v5/account_position_builder_service_test.go) | [ex](../examples/account_position_builder) | `auth-private` |
| `POST /api/v5/account/position-builder-graph` | `private` | [`account_position_builder_graph`](../v5/account_position_builder_graph_service.go) | ✅ [test](../v5/account_position_builder_graph_service_test.go) | [ex](../examples/account_position_builder_graph) | `auth-private` |
| `GET /api/v5/account/position-tiers` | `private` | [`account_position_tiers`](../v5/account_position_tiers_service.go) | ✅ [test](../v5/account_position_tiers_service_test.go) | [ex](../examples/account_position_tiers) | `auth-private` |
| `POST /api/v5/account/position/margin-balance` | `private` | [`account_position_margin_balance`](../v5/account_position_margin_balance_service.go) | ✅ [test](../v5/account_position_margin_balance_service_test.go) | [ex](../examples/account_position_margin_balance) | `auth-private, balance, margin` |

| `GET /api/v5/account/positions` | `private` | [`account_positions`](../v5/account_positions_service.go) | ✅ [test](../v5/account_positions_service_test.go) | [ex](../examples/account_positions) | `auth-private, positions` |
| `GET /api/v5/account/positions-history` | `private` | [`account_positions_history`](../v5/account_positions_history_service.go) | ✅ [test](../v5/account_positions_history_service_test.go) | [ex](../examples/account_positions_history) | `auth-private, positions` |
| `GET /api/v5/account/precheck-set-delta-neutral` | `private` | [`account_precheck_set_delta_neutral`](../v5/account_precheck_set_delta_neutral_service.go) | ✅ [test](../v5/account_precheck_set_delta_neutral_service_test.go) | [ex](../examples/account_precheck_set_delta_neutral) | `auth-private` |
| `GET /api/v5/account/risk-state` | `private` | [`account_risk_state`](../v5/account_risk_state_service.go) | ✅ [test](../v5/account_risk_state_service_test.go) | [ex](../examples/account_risk_state) | `auth-private, risk` |
| `POST /api/v5/account/set-account-level` | `private` | [`account_set_account_level`](../v5/account_set_account_level_service.go) | ✅ [test](../v5/account_set_account_level_service_test.go) | [ex](../examples/account_set_account_level) | `auth-private` |
| `GET /api/v5/account/set-account-switch-precheck` | `private` | [`account_switch_precheck`](../v5/account_switch_precheck_service.go) | ✅ [test](../v5/account_switch_precheck_service_test.go) | [ex](../examples/account_switch_precheck) | `auth-private` |
| `POST /api/v5/account/set-auto-earn` | `private` | [`account_set_auto_earn`](../v5/account_set_auto_earn_service.go) | ✅ [test](../v5/account_set_auto_earn_service_test.go) | [ex](../examples/account_set_auto_earn) | `auth-private` |
| `POST /api/v5/account/set-auto-loan` | `private` | [`account_set_auto_loan`](../v5/account_set_auto_loan_service.go) | ✅ [test](../v5/account_set_auto_loan_service_test.go) | [ex](../examples/account_set_auto_loan) | `auth-private, loan` |
| `POST /api/v5/account/set-auto-repay` | `private` | [`account_set_auto_repay`](../v5/account_set_auto_repay_service.go) | ✅ [test](../v5/account_set_auto_repay_service_test.go) | [ex](../examples/account_set_auto_repay) | `auth-private, loan` |
| `POST /api/v5/account/set-collateral-assets` | `private` | [`account_set_collateral_assets`](../v5/account_set_collateral_assets_service.go) | ✅ [test](../v5/account_set_collateral_assets_service_test.go) | [ex](../examples/account_set_collateral_assets) | `auth-private` |
| `POST /api/v5/account/set-fee-type` | `private` | [`account_set_fee_type`](../v5/account_set_fee_type_service.go) | ✅ [test](../v5/account_set_fee_type_service_test.go) | [ex](../examples/account_set_fee_type) | `auth-private` |
| `POST /api/v5/account/set-greeks` | `private` | [`account_set_greeks`](../v5/account_set_greeks_service.go) | ✅ [test](../v5/account_set_greeks_service_test.go) | [ex](../examples/account_set_greeks) | `auth-private, option` |
| `POST /api/v5/account/set-isolated-mode` | `private` | [`account_set_isolated_mode`](../v5/account_set_isolated_mode_service.go) | ✅ [test](../v5/account_set_isolated_mode_service_test.go) | [ex](../examples/account_set_isolated_mode) | `auth-private` |
| `POST /api/v5/account/set-leverage` | `private` | [`account_set_leverage`](../v5/account_set_leverage_service.go) | ✅ [test](../v5/account_set_leverage_service_test.go) | [ex](../examples/account_set_leverage) | `auth-private, leverage` |
| `POST /api/v5/account/set-position-mode` | `private` | [`account_set_position_mode`](../v5/account_set_position_mode_service.go) | ✅ [test](../v5/account_set_position_mode_service_test.go) | [ex](../examples/account_set_position_mode) | `auth-private` |
| `POST /api/v5/account/set-riskOffset-amt` | `private` | [`account_set_risk_offset_amt`](../v5/account_set_risk_offset_amt_service.go) | ✅ [test](../v5/account_set_risk_offset_amt_service_test.go) | [ex](../examples/account_set_risk_offset_amt) | `auth-private, risk` |
| `POST /api/v5/account/set-settle-currency` | `private` | [`account_set_settle_currency`](../v5/account_set_settle_currency_service.go) | ✅ [test](../v5/account_set_settle_currency_service_test.go) | [ex](../examples/account_set_settle_currency) | `auth-private` |
| `POST /api/v5/account/set-trading-config` | `private` | [`account_set_trading_config`](../v5/account_set_trading_config_service.go) | ✅ [test](../v5/account_set_trading_config_service_test.go) | [ex](../examples/account_set_trading_config) | `auth-private` |
| `GET /api/v5/account/spot-borrow-repay-history` | `private` | [`account_spot_borrow_repay_history`](../v5/account_spot_borrow_repay_history_service.go) | ✅ [test](../v5/account_spot_borrow_repay_history_service_test.go) | [ex](../examples/account_spot_borrow_repay_history) | `auth-private, loan, spot` |
| `POST /api/v5/account/spot-manual-borrow-repay` | `private` | [`account_spot_manual_borrow_repay`](../v5/account_spot_manual_borrow_repay_service.go) | ✅ [test](../v5/account_spot_manual_borrow_repay_service_test.go) | [ex](../examples/account_spot_manual_borrow_repay) | `auth-private, loan, spot` |
| `GET /api/v5/account/subaccount/balances` | `private` | [`account_subaccount_balances`](../v5/account_subaccount_balances_service.go) | ✅ [test](../v5/account_subaccount_balances_service_test.go) | [ex](../examples/account_subaccount_balances) | `auth-private, subaccount, balance` |
| `GET /api/v5/account/subaccount/max-withdrawal` | `private` | [`account_subaccount_max_withdrawal`](../v5/account_subaccount_max_withdrawal_service.go) | ✅ [test](../v5/account_subaccount_max_withdrawal_service_test.go) | [ex](../examples/account_subaccount_max_withdrawal) | `auth-private, subaccount` |
| `GET /api/v5/account/trade-fee` | `private` | [`account_trade_fee`](../v5/account_trade_fee_service.go) | ✅ [test](../v5/account_trade_fee_service_test.go) | [ex](../examples/account_trade_fee) | `auth-private` |

### Asset（资产）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/asset/asset-valuation` | `private` | [`asset_valuation`](../v5/asset_valuation_service.go) | ✅ [test](../v5/asset_valuation_service_test.go) | [ex](../examples/asset_valuation) | `auth-private` |
| `GET /api/v5/asset/balances` | `private` | [`asset_balances`](../v5/asset_balances_service.go) | ✅ [test](../v5/asset_balances_service_test.go) | — | `auth-private` |
| `GET /api/v5/asset/bills` | `private` | [`asset_bills`](../v5/asset_bills_service.go) | ✅ [test](../v5/asset_bills_service_test.go) | [ex](../examples/asset_bills) | `auth-private, bills` |
| `GET /api/v5/asset/bills-history` | `private` | [`asset_bills_history`](../v5/asset_bills_history_service.go) | ✅ [test](../v5/asset_bills_history_service_test.go) | [ex](../examples/asset_bills_history) | `auth-private, bills` |
| `POST /api/v5/asset/cancel-withdrawal` | `private` | [`asset_cancel_withdrawal`](../v5/asset_cancel_withdrawal_service.go) | ✅ [test](../v5/asset_cancel_withdrawal_service_test.go) | [ex](../examples/asset_cancel_withdrawal) | `auth-private, withdrawal` |
| `GET /api/v5/asset/convert/currencies` | `private` | [`asset_convert_currencies`](../v5/asset_convert_currencies_service.go) | ✅ [test](../v5/asset_convert_currencies_service_test.go) | [ex](../examples/asset_convert_currencies) | `auth-private, convert` |
| `GET /api/v5/asset/convert/currency-pair` | `private` | [`asset_convert_currency_pair`](../v5/asset_convert_currency_pair_service.go) | ✅ [test](../v5/asset_convert_currency_pair_service_test.go) | [ex](../examples/asset_convert_currency_pair) | `auth-private, convert` |
| `POST /api/v5/asset/convert/estimate-quote` | `private` | [`asset_convert_estimate_quote`](../v5/asset_convert_estimate_quote_service.go) | ✅ [test](../v5/asset_convert_estimate_quote_service_test.go) | [ex](../examples/asset_convert_estimate_quote) | `auth-private, convert` |
| `GET /api/v5/asset/convert/history` | `private` | [`asset_convert_history`](../v5/asset_convert_history_service.go) | ✅ [test](../v5/asset_convert_history_service_test.go) | [ex](../examples/asset_convert_history) | `auth-private, convert` |
| `POST /api/v5/asset/convert/trade` | `private` | [`asset_convert_trade`](../v5/asset_convert_trade_service.go) | ✅ [test](../v5/asset_convert_trade_service_test.go) | [ex](../examples/asset_convert_trade) | `auth-private, convert` |
| `GET /api/v5/asset/currencies` | `private` | [`asset_currencies`](../v5/asset_currencies_service.go) | ✅ [test](../v5/asset_currencies_service_test.go) | [ex](../examples/asset_currencies) | `auth-private` |
| `GET /api/v5/asset/deposit-address` | `private` | [`asset_deposit_address`](../v5/asset_deposit_address_service.go) | ✅ [test](../v5/asset_deposit_address_service_test.go) | [ex](../examples/asset_deposit_address) | `auth-private, deposit` |
| `GET /api/v5/asset/deposit-history` | `private` | [`asset_deposit_history`](../v5/asset_deposit_history_service.go) | ✅ [test](../v5/asset_deposit_history_service_test.go) | [ex](../examples/asset_deposit_history) | `auth-private, deposit` |
| `GET /api/v5/asset/deposit-withdraw-status` | `private` | [`asset_deposit_withdraw_status`](../v5/asset_deposit_withdraw_status_service.go) | ✅ [test](../v5/asset_deposit_withdraw_status_service_test.go) | [ex](../examples/asset_deposit_withdraw_status) | `auth-private, deposit` |
| `GET /api/v5/asset/exchange-list` | `public` | [`asset_exchange_list`](../v5/asset_exchange_list_service.go) | ✅ [test](../v5/asset_exchange_list_service_test.go) | [ex](../examples/asset_exchange_list) | `auth-public` |
| `GET /api/v5/asset/monthly-statement` | `private` | [`asset_monthly_statement`](../v5/asset_monthly_statement_service.go) | ✅ [test](../v5/asset_monthly_statement_service_test.go) | [ex](../examples/asset_monthly_statement) | `auth-private, statement` |
| `POST /api/v5/asset/monthly-statement` | `private` | [`asset_monthly_statement_apply`](../v5/asset_monthly_statement_apply_service.go) | ✅ [test](../v5/asset_monthly_statement_apply_service_test.go) | [ex](../examples/asset_monthly_statement_apply) | `auth-private, statement` |
| `GET /api/v5/asset/non-tradable-assets` | `private` | [`asset_non_tradable_assets`](../v5/asset_non_tradable_assets_service.go) | ✅ [test](../v5/asset_non_tradable_assets_service_test.go) | [ex](../examples/asset_non_tradable_assets) | `auth-private` |
| `GET /api/v5/asset/subaccount/balances` | `private` | [`asset_subaccount_balances`](../v5/asset_subaccount_balances_service.go) | ✅ [test](../v5/asset_subaccount_balances_service_test.go) | [ex](../examples/asset_subaccount_balances) | `auth-private, subaccount` |
| `GET /api/v5/asset/subaccount/bills` | `private` | [`asset_subaccount_bills`](../v5/asset_subaccount_bills_service.go) | ✅ [test](../v5/asset_subaccount_bills_service_test.go) | [ex](../examples/asset_subaccount_bills) | `auth-private, subaccount, bills` |
| `GET /api/v5/asset/subaccount/managed-subaccount-bills` | `private` | [`asset_subaccount_managed_subaccount_bills`](../v5/asset_subaccount_managed_subaccount_bills_service.go) | ✅ [test](../v5/asset_subaccount_managed_subaccount_bills_service_test.go) | [ex](../examples/asset_subaccount_managed_subaccount_bills) | `auth-private, subaccount, bills` |
| `POST /api/v5/asset/subaccount/transfer` | `private` | [`asset_subaccount_transfer`](../v5/asset_subaccount_transfer_service.go) | ✅ [test](../v5/asset_subaccount_transfer_service_test.go) | [ex](../examples/asset_subaccount_transfer) | `auth-private, subaccount, transfer` |
| `POST /api/v5/asset/transfer` | `private` | [`asset_transfer`](../v5/asset_transfer_service.go) | ✅ [test](../v5/asset_transfer_service_test.go) | — | `auth-private, transfer` |
| `GET /api/v5/asset/transfer-state` | `private` | [`asset_transfer_state`](../v5/asset_transfer_state_service.go) | ✅ [test](../v5/asset_transfer_state_service_test.go) | — | `auth-private, transfer` |
| `POST /api/v5/asset/withdrawal` | `private` | [`asset_withdrawal`](../v5/asset_withdrawal_service.go) | ✅ [test](../v5/asset_withdrawal_service_test.go) | [ex](../examples/asset_withdrawal) | `auth-private, withdrawal` |
| `GET /api/v5/asset/withdrawal-history` | `private` | [`asset_withdrawal_history`](../v5/asset_withdrawal_history_service.go) | ✅ [test](../v5/asset_withdrawal_history_service_test.go) | [ex](../examples/asset_withdrawal_history) | `auth-private, withdrawal` |

### Users（用户/子账户）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/users/entrust-subaccount-list` | `private` | [`users_entrust_subaccount_list`](../v5/users_entrust_subaccount_list_service.go) | ✅ [test](../v5/users_entrust_subaccount_list_service_test.go) | [ex](../examples/users_entrust_subaccount_list) | `auth-private, subaccount` |
| `GET /api/v5/users/subaccount/apikey` | `private` | [`users_subaccount_apikeys`](../v5/users_subaccount_apikeys_service.go) | ✅ [test](../v5/users_subaccount_apikeys_service_test.go) | [ex](../examples/users_subaccount_apikeys) | `auth-private, subaccount, apikey` |
| `POST /api/v5/users/subaccount/apikey` | `private` | [`users_subaccount_create_apikey`](../v5/users_subaccount_create_apikey_service.go) | ✅ [test](../v5/users_subaccount_create_apikey_service_test.go) | [ex](../examples/users_subaccount_create_apikey) | `auth-private, subaccount, apikey` |
| `POST /api/v5/users/subaccount/create-subaccount` | `private` | [`users_subaccount_create_subaccount`](../v5/users_subaccount_create_subaccount_service.go) | ✅ [test](../v5/users_subaccount_create_subaccount_service_test.go) | [ex](../examples/users_subaccount_create_subaccount) | `auth-private, subaccount` |
| `POST /api/v5/users/subaccount/delete-apikey` | `private` | [`users_subaccount_delete_apikey`](../v5/users_subaccount_delete_apikey_service.go) | ✅ [test](../v5/users_subaccount_delete_apikey_service_test.go) | [ex](../examples/users_subaccount_delete_apikey) | `auth-private, subaccount, apikey` |
| `GET /api/v5/users/subaccount/list` | `private` | [`users_subaccount_list`](../v5/users_subaccount_list_service.go) | ✅ [test](../v5/users_subaccount_list_service_test.go) | [ex](../examples/users_subaccount_list) | `auth-private, subaccount` |
| `POST /api/v5/users/subaccount/modify-apikey` | `private` | [`users_subaccount_modify_apikey`](../v5/users_subaccount_modify_apikey_service.go) | ✅ [test](../v5/users_subaccount_modify_apikey_service_test.go) | [ex](../examples/users_subaccount_modify_apikey) | `auth-private, subaccount, apikey` |
| `POST /api/v5/users/subaccount/set-transfer-out` | `private` | [`users_subaccount_set_transfer_out`](../v5/users_subaccount_set_transfer_out_service.go) | ✅ [test](../v5/users_subaccount_set_transfer_out_service_test.go) | [ex](../examples/users_subaccount_set_transfer_out) | `auth-private, subaccount, transfer` |

### RFQ（大宗交易/询价）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/rfq/public-trades` | `public` | [`rfq_public_trades`](../v5/rfq_public_trades_service.go) | ✅ [test](../v5/rfq_public_trades_service_test.go) | [ex](../examples/rfq_public_trades) | `rest, auth-public, rfq` |

### SPRD（价差交易）

| Endpoint | 鉴权 | Service | 测试 | 示例 | 标签 |
|---|---|---|---|---|---|
| `GET /api/v5/sprd/spreads` | `public` | [`sprd_spreads`](../v5/sprd_spreads_service.go) | ✅ [test](../v5/sprd_spreads_service_test.go) | [ex](../examples/sprd_spreads) | `rest, auth-public, sprd` |
| `GET /api/v5/sprd/books` | `public` | [`sprd_books`](../v5/sprd_books_service.go) | ✅ [test](../v5/sprd_books_service_test.go) | [ex](../examples/sprd_books) | `rest, auth-public, sprd, orderbook` |
| `GET /api/v5/sprd/public-trades` | `public` | [`sprd_public_trades`](../v5/sprd_public_trades_service.go) | ✅ [test](../v5/sprd_public_trades_service_test.go) | [ex](../examples/sprd_public_trades) | `rest, auth-public, sprd` |

## WebSocket

> 说明：WS 的订阅以 `WSArg{Channel: ..., InstId/InstType/...}` 为核心；连接管理（自动重连/重订阅/心跳/notice）在 `WSClient` 内部处理。

### Public（无需登录）

- 频道（已解析）：`tickers`、`trades`、深度 `books/books-elp/books5/bbo-tbt/books-l2-tbt/books50-l2-tbt`、`open-interest`、`funding-rate`、`mark-price`、`index-tickers`、`price-limit`、`opt-summary`、`liquidation-orders`、`option-trades`、`call-auction-details`
- Examples：[examples/ws_public_tickers](../examples/ws_public_tickers)、[examples/ws_public_trades](../examples/ws_public_trades)、[examples/ws_public_books](../examples/ws_public_books)、[examples/ws_public_books5](../examples/ws_public_books5)、[examples/ws_public_open_interest](../examples/ws_public_open_interest)、[examples/ws_public_funding_rate](../examples/ws_public_funding_rate)、[examples/ws_public_mark_price](../examples/ws_public_mark_price)、[examples/ws_public_index_tickers](../examples/ws_public_index_tickers)、[examples/ws_public_opt_summary](../examples/ws_public_opt_summary)

### Private（需要登录）

- 频道（已解析 + typed handler）：`orders`、`fills`、`account`、`positions`、`balance_and_position`
- 业务 op（交易闭环）：`order`、`cancel-order`、`amend-order`（含 batch 等待/错误归一）
- Examples：[examples/ws_private_orders](../examples/ws_private_orders)、[examples/ws_private_fills](../examples/ws_private_fills)、[examples/ws_private_account](../examples/ws_private_account)、[examples/ws_private_positions](../examples/ws_private_positions)、[examples/ws_private_balance_and_position](../examples/ws_private_balance_and_position)、[examples/ws_private_trade_order](../examples/ws_private_trade_order)、[examples/ws_private_trade_cancel](../examples/ws_private_trade_cancel)、[examples/ws_private_trade_amend](../examples/ws_private_trade_amend)、[examples/ws_private_trade_batch_ops](../examples/ws_private_trade_batch_ops)

### Business（按频道决定是否需要登录）

- 频道（已解析）：K 线 `candle*`、标记价格 K 线 `mark-price-candle*`、指数 K 线 `index-candle*`、`trades-all`、`sprd-public-trades`
- 资金推送（需要登录）：`deposit-info`、`withdrawal-info`
- Examples：[examples/ws_business_candles](../examples/ws_business_candles)、[examples/ws_business_mark_price_candles](../examples/ws_business_mark_price_candles)、[examples/ws_business_index_candles](../examples/ws_business_index_candles)、[examples/ws_business_trades_all](../examples/ws_business_trades_all)、[examples/ws_business_sprd_public_trades](../examples/ws_business_sprd_public_trades)、[examples/ws_business_deposit_info](../examples/ws_business_deposit_info)、[examples/ws_business_withdrawal_info](../examples/ws_business_withdrawal_info)

## 场景索引（标签）

> 使用方式：在本文件里搜索标签（如 `orderbook`/`orders`/`deposit`），或在代码里用 `rg \"<tag>\" docs/coverage.md` 快速定位相关端点。

- 行情：`tickers`、`trades`、`orderbook`、`candles`
- 下单链路：`orders`、`fills`、`algo`
- 鉴权：`auth-public`、`auth-private`
- 资产：`transfer`、`bills`、`deposit`、`withdrawal`、`convert`、`statement`
- 账户/仓位：`balance`、`positions`、`leverage`、`loan`、`margin`、`spot`
- 子账户：`subaccount`、`apikey`
- 期权相关：`option`、`mmp`
- 风控/限速：`risk`
- 大宗交易：`rfq`
- 价差交易：`sprd`

## 维护说明

- REST 覆盖以 `v5/*_service.go` 中调用 `c.do/c.doWithHeaders` 的 `\"/api/v5/...\"` 为准。
- 端点新增建议同步更新：
  - 单测：优先 `*_service_test.go`（或聚合到 `*_services_test.go`）
  - 示例：优先 `examples/<service>`（按需）
  - 本文件：补齐对应行（或重生成表格）
