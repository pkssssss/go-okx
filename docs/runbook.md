# 运行手册（Runbook）/ 生产故障处置

> 目标：把“能用”升级为“可运行（可观测/可处置/可演练）”。  
> 适用：OKX V5（REST + WebSocket）接入层与交易系统。

---

## 0. 通用处置原则（先保资金，再保服务）

1. **进入安全模式**（先止血）：暂停新开仓/新下单；必要时撤销所有未成交订单；将风险敞口收敛到可控范围。
2. **固定证据**（可追溯）：记录版本（commit）、关键错误、请求 `requestId`、WS 状态快照、当时的限频/网络指标。
3. **定位根因**（可判定）：区分“请求未发送/已发送未响应/已响应错误”；区分“WS 断线/handler 堵塞/失订阅”。
4. **恢复与验证**（不猜测）：恢复后必须做一次“订单/持仓/余额”对账（REST 为准），再逐步放开交易频率。

---

## 1. 证据采集清单（建议固化为结构化日志）

### 1.1 REST

- 错误分类：
  - `*okx.APIError`：已收到 HTTP 响应（可直接用 `RequestID` 对齐交易所侧日志/工单）。
  - `*okx.RequestStateError`：未形成响应的失败（可区分阶段）：
    - `stage=gate dispatched=false`：**请求未发送**（在 request gate 排队/获取名额阶段失败）
    - `stage=http dispatched=true`：**已尝试发送**（进入 HTTP Do 阶段失败/超时）
- 关键字段：`method`、`path`、`requestId`、`code/msg/http`、`stage/dispatched`、`latency`、`retryIndex`

### 1.2 WebSocket

建议周期性采集 `ws.Stats()` 并上报指标/日志：

- 活性：`LastRecv`、`LastPing`、`Connected`
- 稳定性：`Reconnects`、`LastError`
- 订阅质量：`SubscribeOK/SubscribeError`
- 背压：`TypedQueueLen/Cap`、`RawQueueLen/Cap`

---

## 2. 限频（429 / 50011 / 50061）

### 2.1 典型症状

- REST 返回 `*okx.APIError` 且 `okx.IsRateLimitError(err)==true`
- 429/限频码密集出现，P99 延迟飙升，重试次数上升

### 2.2 立即处置（止血）

1. 暂停非关键轮询（行情/状态查询降频，避免“对抗限频”）。
2. 降低并发与发送速率（优先降低交易写入类接口的频率）。
3. 确认已启用 request gate（默认已启用并发闸门）并合理配置：
   - `okx.WithRequestGate(okx.RequestGateConfig{MaxConcurrent: ...})`
4. 主动拉取账户限频并让 SDK 联动调度（SDK 会在 `TradeAccountRateLimitService.Do()` 成功后更新 gate）：
   - `c.NewTradeAccountRateLimitService().Do(ctx)`

### 2.3 恢复与验证

- 观察 429/限频码下降后，**逐步**提高并发/频率；
- 恢复交易前执行一次对账（见第 6 节）。

---

## 3. 鉴权失败（401 / 501xx）

### 3.1 典型症状

- `okx.IsAuthError(err)==true`
- WS private 登录失败（login event/error）

### 3.2 排查步骤（从高概率到低概率）

1. **时间偏差**：先跑一次 `c.SyncTime(ctx)`；再重试登录/签名请求。
2. **凭证三元组**：`APIKey/SecretKey/Passphrase` 是否来自同一个 key；passphrase 是否正确。
3. **权限/白名单**：APIKey 权限（只读/交易/提现）是否满足；IP 白名单是否覆盖当前出口 IP。
4. **环境与模式**：是否误用模拟盘/实盘（`WithDemoTrading`）或网关 URL。

### 3.3 恢复与验证

- 恢复后做一次只读查询（余额/持仓/未成交订单），确认签名链路完全恢复，再放开交易。

---

## 4. WS 重连震荡 / 断流 / 失速

### 4.1 典型症状

- `ws.Stats().Reconnects` 快速增长
- `ws.Stats().LastRecv` 长时间不更新（断流）
- `ws.Stats().TypedQueueLen` / `RawQueueLen` 接近 `Cap`（背压）
- `ws.Stats().LastError` 持续更新（心跳超时/notice 触发重连/订阅失败）

### 4.2 处置步骤

1. **先判定是不是 handler 堵塞导致**：
   - 若 `TypedQueueLen/RawQueueLen` 长期逼近 `Cap`，说明 handler 消费不过来。
2. **降低 handler 负载**（优先级从高到低）：
   - handler 只做轻量分发，把重逻辑扔给你自己的 worker（避免阻塞）。
   - 调大 buffer：`WithWSTypedHandlerAsync(n)` / `WithWSRawHandlerAsync(n)`
   - 如 handler 极轻且追求最小延迟，可改为 inline：`WithWSTypedHandlerInline()` / `WithWSRawHandlerInline()`
3. **订阅可验证**：关键订阅使用 `SubscribeAndWait`，确保失败可见。
4. **必要时重启 WS 客户端**：`Close()` 后重新 `Start()`，并重新对账（见第 6 节）。

---

## 5. 深度断档/校验失败（Order Book）

### 5.1 典型症状

- `store.Apply(...)` 返回：
  - `*okx.WSOrderBookSequenceError`（prevSeqId 断档）
  - `*okx.WSOrderBookChecksumError`（checksum 不一致）
  - `*okx.WSOrderBookNotReadyError`（未收到 snapshot 却收到 update）

### 5.2 处置步骤（确定性恢复）

1. 立即停止基于该深度的交易决策（避免用错误盘口下单）。
2. `store.Reset()` 清空本地状态。
3. 重新订阅并等待 snapshot（推荐 `SubscribeAndWait`）。
4. 恢复后对比一次关键价位/盘口（如 best bid/ask）与 REST 查询（若你使用 REST 深度做旁路校验）。

---

## 6. 对账失败（订单/成交/持仓/余额分叉）

### 6.1 典型症状

- 本地订单状态与交易所不一致（例如本地认为撤单成功但交易所仍挂单）
- 本地持仓/余额与实际偏离
- WS 丢推送/断流后未做补偿查询导致状态漂移

### 6.2 处置步骤（以交易所为准）

1. 暂停交易（防止错误状态继续扩散）。
2. 拉取权威状态（REST）并重建本地状态机：
   - 余额/持仓/未成交订单/成交明细（按你系统的关键对象选择对应 Service）。
3. 对“不确定窗口”（请求超时/断线期间）的订单，禁止盲重试：
   - 通过 `clOrdId` 幂等键 + 查询确认订单真实状态，再决定下一步。
4. 必要时执行一次“只撤不加”的风险收敛（撤单/降杠杆/平仓），再逐步恢复策略。

---

## 7. 演练建议（不演练=不生产）

- 限频演练：压测到出现 429，验证 gate 生效、429 不形成风暴、可降级与恢复。
- WS 断网演练：验证重连+重订阅可恢复；恢复后对账通过。
- handler 堆积演练：刻意让 handler 变慢，观察队列堆积与断连行为是否符合预期，并调整 buffer/负载。
- 深度断档演练：模拟 prevSeqId 断档，验证 `Reset + 重订阅` 能确定性恢复。

