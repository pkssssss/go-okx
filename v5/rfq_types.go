package okx

// RFQ 表示询价单（REST/WS 通用）。
type RFQ = WSRFQ

// RFQLeg 表示询价单腿（REST/WS 通用）。
type RFQLeg = WSRFQLeg

// RFQAcctAlloc 表示组合询价单账户分配（REST/WS 通用）。
type RFQAcctAlloc = WSRFQAcctAlloc

// RFQAcctAllocLeg 表示组合询价单账户分配的腿（REST/WS 通用）。
type RFQAcctAllocLeg = WSRFQAcctAllocLeg

// Quote 表示报价单（REST/WS 通用）。
type Quote = WSQuote

// QuoteLeg 表示报价单腿（REST/WS 通用）。
type QuoteLeg = WSQuoteLeg

// StrucBlockTrade 表示大宗交易（REST/WS 通用）。
type StrucBlockTrade = WSStrucBlockTrade

// StrucBlockTradeLeg 表示大宗交易腿（REST/WS 通用）。
type StrucBlockTradeLeg = WSStrucBlockTradeLeg

// StrucBlockTradeAcctAlloc 表示组合询价单成交的账户分配（REST/WS 通用）。
type StrucBlockTradeAcctAlloc = WSStrucBlockTradeAcctAlloc

// StrucBlockTradeAcctAllocLeg 表示组合询价单成交的账户分配腿（REST/WS 通用）。
type StrucBlockTradeAcctAllocLeg = WSStrucBlockTradeAcctAllocLeg
