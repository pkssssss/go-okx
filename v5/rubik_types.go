package okx

import (
	"encoding/json"
	"errors"
	"strconv"
)

func rubikUnmarshalStringArray(data []byte, minLen int) ([]string, error) {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return nil, err
	}
	if len(arr) < minLen {
		return nil, errors.New("okx: invalid rubik data")
	}
	return arr, nil
}

func rubikParseUnixMilli(s string) (int64, error) {
	ts, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New("okx: invalid rubik ts")
	}
	return ts, nil
}

// RubikSupportCoin 表示交易大数据支持的币种。
type RubikSupportCoin struct {
	Contract []string `json:"contract"`
	Option   []string `json:"option"`
	Spot     []string `json:"spot"`
}

// RubikOpenInterestHistory 表示合约持仓量历史项（[ts, oi, oiCcy, oiUsd]）。
// 数值字段保持为 string（无损）。
type RubikOpenInterestHistory struct {
	TS int64

	OI    string
	OICcy string
	OIUsd string
}

func (r *RubikOpenInterestHistory) UnmarshalJSON(data []byte) error {
	*r = RubikOpenInterestHistory{}

	arr, err := rubikUnmarshalStringArray(data, 4)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.OI = arr[1]
	r.OICcy = arr[2]
	r.OIUsd = arr[3]
	return nil
}

// RubikTsRatio 表示二维时间序列（[ts, ratio]）。
// 数值字段保持为 string（无损）。
type RubikTsRatio struct {
	TS int64

	Ratio string
}

func (r *RubikTsRatio) UnmarshalJSON(data []byte) error {
	*r = RubikTsRatio{}

	arr, err := rubikUnmarshalStringArray(data, 2)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.Ratio = arr[1]
	return nil
}

// RubikTakerVolume 表示主动买入/卖出情况（[ts, sellVol, buyVol]）。
// 数值字段保持为 string（无损）。
type RubikTakerVolume struct {
	TS int64

	SellVol string
	BuyVol  string
}

func (r *RubikTakerVolume) UnmarshalJSON(data []byte) error {
	*r = RubikTakerVolume{}

	arr, err := rubikUnmarshalStringArray(data, 3)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.SellVol = arr[1]
	r.BuyVol = arr[2]
	return nil
}

// RubikOpenInterestVolume 表示持仓量及交易量（[ts, oi, vol]）。
// 数值字段保持为 string（无损）。
type RubikOpenInterestVolume struct {
	TS int64

	OI  string
	Vol string
}

func (r *RubikOpenInterestVolume) UnmarshalJSON(data []byte) error {
	*r = RubikOpenInterestVolume{}

	arr, err := rubikUnmarshalStringArray(data, 3)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.OI = arr[1]
	r.Vol = arr[2]
	return nil
}

// RubikOptionOpenInterestVolumeRatio 表示期权看涨/看跌持仓总量比/交易总量比（[ts, oiRatio, volRatio]）。
// 数值字段保持为 string（无损）。
type RubikOptionOpenInterestVolumeRatio struct {
	TS int64

	OIRatio  string
	VolRatio string
}

func (r *RubikOptionOpenInterestVolumeRatio) UnmarshalJSON(data []byte) error {
	*r = RubikOptionOpenInterestVolumeRatio{}

	arr, err := rubikUnmarshalStringArray(data, 3)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.OIRatio = arr[1]
	r.VolRatio = arr[2]
	return nil
}

// RubikOptionOpenInterestVolumeExpiry 表示按到期日分的看涨/看跌持仓量及交易量（[ts, expTime, callOI, putOI, callVol, putVol]）。
// 数值字段保持为 string（无损）。
type RubikOptionOpenInterestVolumeExpiry struct {
	TS int64

	ExpTime string
	CallOI  string
	PutOI   string
	CallVol string
	PutVol  string
}

func (r *RubikOptionOpenInterestVolumeExpiry) UnmarshalJSON(data []byte) error {
	*r = RubikOptionOpenInterestVolumeExpiry{}

	arr, err := rubikUnmarshalStringArray(data, 6)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.ExpTime = arr[1]
	r.CallOI = arr[2]
	r.PutOI = arr[3]
	r.CallVol = arr[4]
	r.PutVol = arr[5]
	return nil
}

// RubikOptionOpenInterestVolumeStrike 表示按执行价格分的看涨/看跌持仓量及交易量（[ts, strike, callOI, putOI, callVol, putVol]）。
// 数值字段保持为 string（无损）。
type RubikOptionOpenInterestVolumeStrike struct {
	TS int64

	Strike  string
	CallOI  string
	PutOI   string
	CallVol string
	PutVol  string
}

func (r *RubikOptionOpenInterestVolumeStrike) UnmarshalJSON(data []byte) error {
	*r = RubikOptionOpenInterestVolumeStrike{}

	arr, err := rubikUnmarshalStringArray(data, 6)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.Strike = arr[1]
	r.CallOI = arr[2]
	r.PutOI = arr[3]
	r.CallVol = arr[4]
	r.PutVol = arr[5]
	return nil
}

// RubikOptionTakerBlockVolume 表示看跌/看涨期权合约主动买入/卖出量（data 为单数组）。
// 数值字段保持为 string（无损）。
type RubikOptionTakerBlockVolume struct {
	TS int64

	CallBuyVol   string
	CallSellVol  string
	PutBuyVol    string
	PutSellVol   string
	CallBlockVol string
	PutBlockVol  string
}

func (r *RubikOptionTakerBlockVolume) UnmarshalJSON(data []byte) error {
	*r = RubikOptionTakerBlockVolume{}

	arr, err := rubikUnmarshalStringArray(data, 7)
	if err != nil {
		return err
	}
	ts, err := rubikParseUnixMilli(arr[0])
	if err != nil {
		return err
	}
	r.TS = ts
	r.CallBuyVol = arr[1]
	r.CallSellVol = arr[2]
	r.PutBuyVol = arr[3]
	r.PutSellVol = arr[4]
	r.CallBlockVol = arr[5]
	r.PutBlockVol = arr[6]
	return nil
}
