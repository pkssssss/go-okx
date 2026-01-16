# 使用指南

> 模块：`github.com/pkssssss/go-okx/v5`  
> 最低 Go 版本：`go1.25`

## 1. 安装与运行方式

本仓库使用 Go workspace（`go.work`），建议从仓库根目录运行示例：

```bash
go run ./examples/public_time
```

SDK 主模块代码位于 `v5/`；示例位于 `examples/`。

## 2. Client 初始化

### 2.1 公共接口（无需鉴权）

```go
c := okx.NewClient()
```

### 2.2 私有接口（需要鉴权）

```go
c := okx.NewClient(
	okx.WithCredentials(okx.Credentials{
		APIKey:     os.Getenv("OKX_API_KEY"),
		SecretKey:  os.Getenv("OKX_API_SECRET"),
		Passphrase: os.Getenv("OKX_API_PASSPHRASE"),
	}),
	okx.WithDemoTrading(os.Getenv("OKX_DEMO") == "1"),
)
```

如果要使用 WS private/business（需要登录），建议先做一次校时以降低登录失败概率：

```go
_, _ = c.SyncTime(ctx)
```

## 3. 常用入口（你大概率只需要这些）

### 3.1 构造与配置

- `okx.NewClient(...)`
- `okx.WithCredentials(...)`：私有 REST/WS
- `okx.WithDemoTrading(true)`：模拟盘
- `(*Client).SyncTime(ctx)`：建议 WS 登录前调用

### 3.2 REST

- 统一风格：`c.NewXXXService().<Setters...>().Do(ctx)`
- 常用定位方式：直接在 [`coverage.md`](coverage.md) 搜 endpoint（每行都链接到 Service/Test/Example）

### 3.3 WebSocket

- 端点：`c.NewWSPublic()` / `c.NewWSPrivate()` / `c.NewWSBusiness()` / `c.NewWSBusinessPrivate()`
- 订阅：`SubscribeAndWait`（推荐）/ `Subscribe`
- typed handler：`ws.OnTickers/OnTrades/OnOrderBook/OnOrders/...`
- handler 较重：`okx.WithWSTypedHandlerAsync(1024)`
- 深度合并：`okx.NewWSOrderBookStore(channel, instId)`（配合 `OnOrderBook`）

## 4. REST 调用风格

端点以 `Service + Do(ctx)` 形式暴露（对标 `go-binance`）：

```go
ticker, err := c.NewMarketTickerService().InstId("BTC-USDT").Do(ctx)
```

完整端点清单与对应 Service/Test/Example 以 [`coverage.md`](coverage.md) 为准。

## 5. WebSocket 使用建议

### 4.1 选择 WS 端点

- `c.NewWSPublic()`：公共数据，无需登录
- `c.NewWSPrivate()`：私有数据，需要登录（订单/成交/账户/仓位等）
- `c.NewWSBusiness()`：business（是否需要登录取决于频道）
- `c.NewWSBusinessPrivate()`：business + 强制登录

详细说明见 [`ws.md`](ws.md)。

### 4.2 typed handler（推荐）

默认 typed handler 在 WS read goroutine 中执行；若 handler 较重，建议启用异步队列：

```go
ws := c.NewWSPrivate(okx.WithWSTypedHandlerAsync(1024))
```

深度（books 系列）建议配合 `WSOrderBookStore` 做 snapshot/update 合并与 seq/checksum 校验，见示例 `examples/ws_public_books_store_typed`。

## 6. 类型/精度约定（字段策略）

- 价格/数量/费率等小数：SDK 层优先用 `string`（无损），避免 `float64` 精度问题。
- 时间戳：常见为 Unix 毫秒（string/number），部分字段使用 `UnixMilli` 兼容解析。
- 枚举：多为 `string`（建议上层自行做常量约束/校验）。

## 7. 错误处理

REST 失败会返回 `*okx.APIError`（HTTP 错误或业务 `code != "0"`），可用 `errors.As` 获取结构化信息：

```go
var apiErr *okx.APIError
if errors.As(err, &apiErr) {
	// apiErr.Code / apiErr.Message / apiErr.HTTPStatus / apiErr.RequestPath ...
}
```

常用判断函数：

- `okx.IsAuthError(err)`
- `okx.IsRateLimitError(err)`
- `okx.IsTimeSkewError(err)`

## 8. 如何快速定位“某个接口怎么用”

优先使用覆盖矩阵：[`coverage.md`](coverage.md)（每一行都链接到 Service/Test/Example）。

常用操作：

1) 在 [`coverage.md`](coverage.md) 搜索 endpoint（如 `GET /api/v5/market/ticker`）  
2) 打开对应 `v5/*_service.go` 看参数/返回类型  
3) 直接运行对应 `examples/*`（通常可用默认参数跑通；私有接口再补齐环境变量）

## 9. 示例运行的安全约束

部分示例会触发真实交易/撤单/改单（尤其 trade/WS op 链路），通常要求显式设置：

```bash
export OKX_CONFIRM=YES
```

强烈建议先使用模拟盘：

```bash
export OKX_DEMO=1
```
