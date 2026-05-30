//go:build goexperiment.jsonv2

package opt_test

import (
	jsonv2 "encoding/json/v2"
	"slices"
	"testing"

	"github.com/ngicks/option/opt"
)

func TestOption_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name string
		in   opt.Option[int]
		want string
	}{
		{"some", opt.Some(123), `123`},
		{"some zero value", opt.Some(0), `0`},
		{"none", opt.None[int](), `null`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := jsonv2.Marshal(tt.in)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(got) != tt.want {
				t.Errorf("Marshal(%v) = %s, want %s", tt.in, got, tt.want)
			}
		})
	}
}

func TestOption_UnmarshalJSONFrom(t *testing.T) {
	t.Run("null decodes to none", func(t *testing.T) {
		// Start from a some value to ensure it is reset to none.
		o := opt.Some(999)
		if err := jsonv2.Unmarshal([]byte(`null`), &o); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if o.IsSome() {
			t.Errorf("want none, got some(%v)", o.Value())
		}
		if o.Value() != 0 {
			t.Errorf("want zero value after none, got %v", o.Value())
		}
	})

	t.Run("value decodes to some", func(t *testing.T) {
		var o opt.Option[int]
		if err := jsonv2.Unmarshal([]byte(`123`), &o); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !o.IsSome() || o.Value() != 123 {
			t.Errorf("want some(123), got some=%t value=%v", o.IsSome(), o.Value())
		}
	})

	t.Run("decode error keeps option unchanged", func(t *testing.T) {
		// A failed decode must not leave o in a half-baked state.
		o := opt.Some(7)
		if err := jsonv2.Unmarshal([]byte(`"not a number"`), &o); err == nil {
			t.Fatal("want error, got nil")
		}
		if !o.IsSome() || o.Value() != 7 {
			t.Errorf("state changed on error: some=%t value=%v", o.IsSome(), o.Value())
		}
	})
}

func TestOption_JSONv2_RoundTrip(t *testing.T) {
	type wrapper struct {
		A opt.Option[int]    `json:"a"`
		B opt.Option[string] `json:"b,omitzero"`
		C opt.Option[[]int]  `json:"c"`
	}

	t.Run("none b is omitted via omitzero", func(t *testing.T) {
		got, err := jsonv2.Marshal(wrapper{A: opt.Some(1), B: opt.None[string](), C: opt.Some([]int{1, 2})})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"a":1,"c":[1,2]}`; string(got) != want {
			t.Errorf("Marshal = %s, want %s", got, want)
		}
	})

	t.Run("some b is kept", func(t *testing.T) {
		got, err := jsonv2.Marshal(wrapper{A: opt.None[int](), B: opt.Some("x"), C: opt.None[[]int]()})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"a":null,"b":"x","c":null}`; string(got) != want {
			t.Errorf("Marshal = %s, want %s", got, want)
		}
	})

	t.Run("round trip", func(t *testing.T) {
		in := wrapper{A: opt.Some(42), B: opt.Some("hello"), C: opt.Some([]int{3, 4, 5})}
		b, err := jsonv2.Marshal(in)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var out wrapper
		if err := jsonv2.Unmarshal(b, &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if out.A != in.A {
			t.Errorf("A: got %v, want %v", out.A, in.A)
		}
		if out.B != in.B {
			t.Errorf("B: got %v, want %v", out.B, in.B)
		}
		if !out.C.IsSome() || !slices.Equal(out.C.Value(), in.C.Value()) {
			t.Errorf("C: got %v, want %v", out.C, in.C)
		}
	})
}
