package okx

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// SensitiveString 用于承载“必须返回给调用方，但不应默认外显”的敏感字符串（如 secretKey/passphrase）。
//
// 设计要点：
// - Value() 返回原始值（用于安全存储/二次处理）。
// - fmt / JSON 默认输出脱敏后的值，降低误打日志风险。
type SensitiveString struct {
	v string
}

// Value 返回未脱敏的原始值。
func (s SensitiveString) Value() string {
	return s.v
}

func (s SensitiveString) String() string {
	if s.v == "" {
		return ""
	}
	return "REDACTED"
}

func (s SensitiveString) Format(w fmt.State, verb rune) {
	redacted := s.String()
	switch verb {
	case 'q':
		_, _ = io.WriteString(w, strconv.Quote(redacted))
	default:
		_, _ = io.WriteString(w, redacted)
	}
}

func (s SensitiveString) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *SensitiveString) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	s.v = v
	return nil
}
