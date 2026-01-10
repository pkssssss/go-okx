package okx

// EasyConvertFromData 表示小币一键兑换可兑换币种信息。
// 数值字段按 OKX 返回保持为 string（无损）。
type EasyConvertFromData struct {
	FromAmt string `json:"fromAmt"`
	FromCcy string `json:"fromCcy"`
}

// EasyConvertCurrencyList 表示小币一键兑换主流币币种列表返回项。
type EasyConvertCurrencyList struct {
	FromData []EasyConvertFromData `json:"fromData"`
	ToCcy    []string              `json:"toCcy"`
}

// EasyConvertAck 表示小币一键兑换主流币交易返回项。
// 数值字段按 OKX 返回保持为 string（无损）。
type EasyConvertAck struct {
	FillFromSz string `json:"fillFromSz"`
	FillToSz   string `json:"fillToSz"`
	FromCcy    string `json:"fromCcy"`
	Status     string `json:"status"`
	ToCcy      string `json:"toCcy"`
	UTime      int64  `json:"uTime,string"`
}

// EasyConvertHistory 表示小币一键兑换主流币历史记录返回项。
// 数值字段按 OKX 返回保持为 string（无损）。
type EasyConvertHistory struct {
	FillFromSz string `json:"fillFromSz"`
	FillToSz   string `json:"fillToSz"`
	FromCcy    string `json:"fromCcy"`
	Status     string `json:"status"`
	ToCcy      string `json:"toCcy"`
	Acct       string `json:"acct"`
	UTime      int64  `json:"uTime,string"`
}
