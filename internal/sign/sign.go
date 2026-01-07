package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

const wsLoginRequestPath = "/users/self/verify"

// TimestampISO8601Millis 返回 OKX REST 使用的 UTC ISO8601（毫秒精度）时间戳。
func TimestampISO8601Millis(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

// TimestampUnixSeconds 返回 Unix Epoch 秒（UTC）。
func TimestampUnixSeconds(t time.Time) string {
	return strconv.FormatInt(t.UTC().Unix(), 10)
}

// PrehashREST 构造 OKX REST 的签名前字符串：timestamp + method + requestPath + body。
func PrehashREST(timestamp, method, requestPath, body string) string {
	return timestamp + method + requestPath + body
}

// PrehashWSLogin 构造 OKX WS 登录的签名前字符串：timestamp + "GET" + "/users/self/verify"。
func PrehashWSLogin(timestamp string) string {
	return timestamp + "GET" + wsLoginRequestPath
}

// SignHMACSHA256Base64 对 prehash 进行 HMAC-SHA256，再做 Base64 编码（OKX 使用该形式）。
func SignHMACSHA256Base64(secret, prehash string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(prehash))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
