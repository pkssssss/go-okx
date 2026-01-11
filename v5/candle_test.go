package okx

import (
	"encoding/json"
	"testing"
)

func TestCandle_UnmarshalJSON(t *testing.T) {
	t.Run("sprd_candles_7_fields", func(t *testing.T) {
		var c Candle
		if err := json.Unmarshal([]byte(`["1","1","2","0.5","1.5","100","0"]`), &c); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		if c.TS != 1 || c.Open != "1" || c.Vol != "100" || c.Confirm != "0" {
			t.Fatalf("candle = %#v", c)
		}
		if c.VolCcy != "" || c.VolCcyQuote != "" {
			t.Fatalf("unexpected optional fields: VolCcy=%q VolCcyQuote=%q", c.VolCcy, c.VolCcyQuote)
		}
	})

	t.Run("8_fields_volCcy_and_confirm", func(t *testing.T) {
		var c Candle
		if err := json.Unmarshal([]byte(`["1","1","2","0.5","1.5","100","200","1"]`), &c); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		if c.VolCcy != "200" || c.Confirm != "1" {
			t.Fatalf("candle = %#v", c)
		}
		if c.VolCcyQuote != "" {
			t.Fatalf("unexpected VolCcyQuote = %q", c.VolCcyQuote)
		}
	})

	t.Run("9_fields_full", func(t *testing.T) {
		var c Candle
		if err := json.Unmarshal([]byte(`["1","1","2","0.5","1.5","100","200","300","1"]`), &c); err != nil {
			t.Fatalf("json.Unmarshal() error = %v", err)
		}
		if c.VolCcy != "200" || c.VolCcyQuote != "300" || c.Confirm != "1" {
			t.Fatalf("candle = %#v", c)
		}
	})
}
