package okx

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// UnixMilli 是兼容 string/number 两种 JSON 表达的 Unix 毫秒时间戳。
type UnixMilli int64

func (t *UnixMilli) UnmarshalJSON(b []byte) error {
	b = bytes.TrimSpace(b)
	if len(b) == 0 || bytes.Equal(b, []byte("null")) {
		*t = 0
		return nil
	}

	if b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		if s == "" {
			*t = 0
			return nil
		}
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*t = UnixMilli(v)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	v, err := n.Int64()
	if err != nil {
		return err
	}
	*t = UnixMilli(v)
	return nil
}
