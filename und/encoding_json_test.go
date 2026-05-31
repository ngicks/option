package und_test

import (
	"encoding/json"
	"testing"

	"github.com/ngicks/option/und"
)

func TestUnd_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   und.Und[int]
		want string
	}{
		{"defined", und.Defined(123), `123`},
		{"defined zero", und.Defined(0), `0`},
		{"null", und.Null[int](), `null`},
		{"undefined", und.Undefined[int](), `null`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

func TestUnd_UnmarshalJSON(t *testing.T) {
	t.Run("null decodes to null state", func(t *testing.T) {
		u := und.Defined(999)
		if err := json.Unmarshal([]byte(`null`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsNull() {
			t.Errorf("want null, got state %d", u.State())
		}
	})

	t.Run("value decodes to defined", func(t *testing.T) {
		var u und.Und[int]
		if err := json.Unmarshal([]byte(`123`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsDefined() || u.Value() != 123 {
			t.Errorf("want defined(123), got defined=%t value=%d", u.IsDefined(), u.Value())
		}
	})

	t.Run("decode error keeps und unchanged", func(t *testing.T) {
		u := und.Defined(7)
		if err := json.Unmarshal([]byte(`"not a number"`), &u); err == nil {
			t.Fatal("want error, got nil")
		}
		if !u.IsDefined() || u.Value() != 7 {
			t.Errorf("state changed on error: defined=%t value=%d", u.IsDefined(), u.Value())
		}
	})
}

func TestUnd_JSON_omitzero(t *testing.T) {
	type wrapper struct {
		A und.Und[int]    `json:"a,omitzero"`
		B und.Und[string] `json:"b,omitzero"`
		C und.Und[int]    `json:"c,omitzero"`
	}

	// undefined (C) is omitted via omitzero/IsZero; null (B) is kept as null.
	got, err := json.Marshal(wrapper{A: und.Defined(1), B: und.Null[string](), C: und.Undefined[int]()})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"a":1,"b":null}`; string(got) != want {
		t.Errorf("Marshal = %s, want %s", got, want)
	}
}

func TestUnd_JSON_RoundTrip(t *testing.T) {
	type wrapper struct {
		A und.Und[int]    `json:"a"`
		B und.Und[string] `json:"b"`
	}
	in := wrapper{A: und.Defined(42), B: und.Null[string]()}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var out wrapper
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !out.A.IsDefined() || out.A.Value() != 42 {
		t.Errorf("A: got %v, want defined(42)", out.A)
	}
	if !out.B.IsNull() {
		t.Errorf("B: got state %d, want null", out.B.State())
	}
}
