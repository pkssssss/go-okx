package okx

import (
	"encoding/json"
	"errors"
	"strconv"
)

// PriceCandle 表示指数/标记价格 K 线数据。
//
// OKX 返回为数组：["ts","o","h","l","c","confirm"]
type PriceCandle struct {
	TS int64

	Open  string
	High  string
	Low   string
	Close string

	Confirm string
}

func (c *PriceCandle) UnmarshalJSON(data []byte) error {
	*c = PriceCandle{}

	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) < 6 {
		return errors.New("okx: invalid price candle")
	}

	ts, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return errors.New("okx: invalid price candle ts")
	}

	c.TS = ts
	c.Open = arr[1]
	c.High = arr[2]
	c.Low = arr[3]
	c.Close = arr[4]
	c.Confirm = arr[5]
	return nil
}
