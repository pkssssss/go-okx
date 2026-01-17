package redact

// Mask 用于在示例日志中脱敏敏感标识（如 APIKey）。
// 规则：保留前后各 4 位，中间用 "****" 替换；过短字符串直接返回 "****"。
func Mask(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}
