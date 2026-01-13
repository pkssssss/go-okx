package okx

// FiatBuySellCcy 表示买卖交易币种项。
type FiatBuySellCcy struct {
	Ccy string `json:"ccy"`
}

// FiatBuySellCurrencies 表示买卖交易可用的法币/加密货币列表。
type FiatBuySellCurrencies struct {
	FiatCcyList   []FiatBuySellCcy `json:"fiatCcyList"`
	CryptoCcyList []FiatBuySellCcy `json:"cryptoCcyList"`
}

// FiatBuySellCurrencyPair 表示买卖交易币对能力信息。
// 数值字段保持为 string（无损）。
type FiatBuySellCurrencyPair struct {
	Side    string `json:"side"`
	FromCcy string `json:"fromCcy"`
	ToCcy   string `json:"toCcy"`

	SingleTradeMax string `json:"singleTradeMax"`
	SingleTradeMin string `json:"singleTradeMin"`

	FixedPxDailyLimit          string `json:"fixedPxDailyLimit"`
	FixedPxRemainingDailyQuota string `json:"fixedPxRemainingDailyQuota"`

	PaymentMethods []string `json:"paymentMethods"`
}

// FiatBuySellQuote 表示买卖交易报价。
// 数值字段保持为 string（无损）。
type FiatBuySellQuote struct {
	QuoteId string `json:"quoteId"`
	Side    string `json:"side,omitempty"`

	FromCcy string `json:"fromCcy"`
	ToCcy   string `json:"toCcy"`

	RfqAmt string `json:"rfqAmt"`
	RfqCcy string `json:"rfqCcy"`

	QuotePx      string    `json:"quotePx"`
	QuoteCcy     string    `json:"quoteCcy"`
	QuoteFromAmt string    `json:"quoteFromAmt"`
	QuoteToAmt   string    `json:"quoteToAmt"`
	QuoteTime    UnixMilli `json:"quoteTime"`
	TtlMs        string    `json:"ttlMs"`
}

// FiatBuySellOrder 表示买卖交易订单/历史条目。
// 数值字段保持为 string（无损）。
type FiatBuySellOrder struct {
	OrdId   string `json:"ordId"`
	ClOrdId string `json:"clOrdId"`
	QuoteId string `json:"quoteId"`

	State string `json:"state,omitempty"`
	Side  string `json:"side,omitempty"`

	FromCcy string `json:"fromCcy"`
	ToCcy   string `json:"toCcy"`

	RfqAmt string `json:"rfqAmt"`
	RfqCcy string `json:"rfqCcy"`

	FillPx       string `json:"fillPx"`
	FillQuoteCcy string `json:"fillQuoteCcy"`
	FillFromAmt  string `json:"fillFromAmt"`
	FillToAmt    string `json:"fillToAmt"`

	CTime UnixMilli `json:"cTime"`
	UTime UnixMilli `json:"uTime,omitempty"`
}
