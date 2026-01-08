package sign

import (
	"testing"
	"time"
)

func TestTimestampISO8601Millis(t *testing.T) {
	tm := time.Date(2020, 3, 28, 12, 21, 41, 274_000_000, time.UTC)
	if got, want := TimestampISO8601Millis(tm), "2020-03-28T12:21:41.274Z"; got != want {
		t.Fatalf("TimestampISO8601Millis() = %q, want %q", got, want)
	}
}

func TestTimestampUnixSeconds(t *testing.T) {
	tm := time.Unix(1538054050, 0).UTC()
	if got, want := TimestampUnixSeconds(tm), "1538054050"; got != want {
		t.Fatalf("TimestampUnixSeconds() = %q, want %q", got, want)
	}
}

func TestPrehashREST(t *testing.T) {
	timestamp := "2020-12-08T09:08:57.715Z"
	method := "GET"
	requestPath := "/api/v5/account/balance?ccy=BTC"
	body := ""
	if got, want := PrehashREST(timestamp, method, requestPath, body), "2020-12-08T09:08:57.715ZGET/api/v5/account/balance?ccy=BTC"; got != want {
		t.Fatalf("PrehashREST() = %q, want %q", got, want)
	}
}

func TestPrehashWSLogin(t *testing.T) {
	timestamp := "1538054050"
	if got, want := PrehashWSLogin(timestamp), "1538054050GET/users/self/verify"; got != want {
		t.Fatalf("PrehashWSLogin() = %q, want %q", got, want)
	}
}

func TestSignHMACSHA256Base64(t *testing.T) {
	secret := "mysecret"

	t.Run("rest_get_with_query", func(t *testing.T) {
		prehash := "2020-12-08T09:08:57.715ZGET/api/v5/account/balance?ccy=BTC"
		if got, want := SignHMACSHA256Base64(secret, prehash), "OkAjRTXqfIRKxx7SRjIowU96vZPkf4n9X2G+8yduCf4="; got != want {
			t.Fatalf("SignHMACSHA256Base64() = %q, want %q", got, want)
		}
	})

	t.Run("rest_post_with_body", func(t *testing.T) {
		prehash := "2020-03-28T12:21:41.274ZPOST/api/v5/trade/order{\"instId\":\"BTC-USDT\",\"tdMode\":\"isolated\",\"side\":\"buy\",\"ordType\":\"limit\",\"px\":\"1\",\"sz\":\"1\"}"
		if got, want := SignHMACSHA256Base64(secret, prehash), "5JgIEfkRBluy4x31t6uitZqzoshK+kWWjq9f597WRqQ="; got != want {
			t.Fatalf("SignHMACSHA256Base64() = %q, want %q", got, want)
		}
	})

	t.Run("ws_login", func(t *testing.T) {
		prehash := "1538054050GET/users/self/verify"
		if got, want := SignHMACSHA256Base64(secret, prehash), "m+lzVL6siKIpimAa/6y8lHpWZe0SCpehAqymC8Nel0A="; got != want {
			t.Fatalf("SignHMACSHA256Base64() = %q, want %q", got, want)
		}
	})
}
