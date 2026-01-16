# go-okx

OKX V5 API 的 Go SDK（REST + WebSocket），强调正确性与稳定性。

## 安装

```bash
go get github.com/pkssssss/go-okx/v5
```

> 最低 Go 版本：`go1.25`

## 快速开始

仓库使用 Go workspace（`go.work`），可直接从仓库根目录运行 `examples/`；SDK 主模块位于 `v5/`。

```bash
# 1) 公共接口：获取系统时间
go run ./examples/public_time

# 2) REST：获取单个产品行情（默认 BTC-USDT）
go run ./examples/market_ticker

# 3) WS：订阅公共 tickers（收到首条消息后退出）
go run ./examples/ws_public_tickers
```

私有接口（REST/WS）需要设置凭证，建议先使用模拟盘：

```bash
export OKX_API_KEY="..."
export OKX_API_SECRET="..."
export OKX_API_PASSPHRASE="..."
export OKX_DEMO=1
```

## 文档

- 文档导航：[`docs/README.md`](docs/README.md)
- 使用指南：[`docs/guide.md`](docs/guide.md)
- 覆盖矩阵（端点 -> Service/Test/Example）：[`docs/coverage.md`](docs/coverage.md)
- WebSocket 指南：[`docs/ws.md`](docs/ws.md)
- 设计记录：[`docs/design.md`](docs/design.md)
- Roadmap：[`docs/roadmap.md`](docs/roadmap.md)
