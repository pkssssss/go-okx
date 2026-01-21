package okx

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSensitiveString_RedactsOnFmtAndJSON(t *testing.T) {
	s := SensitiveString{v: "secret"}

	if got, want := s.Value(), "secret"; got != want {
		t.Fatalf("Value() = %q, want %q", got, want)
	}
	if got, want := fmt.Sprintf("%v", s), "REDACTED"; got != want {
		t.Fatalf("fmt %%v = %q, want %q", got, want)
	}
	if got, want := fmt.Sprintf("%s", s), "REDACTED"; got != want {
		t.Fatalf("fmt %%s = %q, want %q", got, want)
	}

	b, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got, want := string(b), `"REDACTED"`; got != want {
		t.Fatalf("json = %q, want %q", got, want)
	}
}
