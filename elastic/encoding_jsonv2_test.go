package elastic_test

import (
	jsonv2 "encoding/json/v2"
	"slices"
	"testing"

	"github.com/ngicks/option/elastic"
	"github.com/ngicks/option/opt"
)

func TestElastic_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name string
		in   elastic.Elastic[int]
		want string
	}{
		{"multiple values", elastic.FromValues(1, 2, 3), `[1,2,3]`},
		{"single value", elastic.FromValue(5), `[5]`},
		{"none options", elastic.FromOptions(opt.Some(1), opt.None[int](), opt.Some(3)), `[1,null,3]`},
		{"empty", elastic.FromValues[int](), `[]`},
		{"null", elastic.Null[int](), `null`},
		{"undefined", elastic.Undefined[int](), `null`},
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

func TestElastic_UnmarshalJSONFrom(t *testing.T) {
	t.Run("null decodes to null state", func(t *testing.T) {
		e := elastic.FromValue(5)
		if err := jsonv2.Unmarshal([]byte(`null`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsNull() {
			t.Errorf("want null, got defined=%t undefined=%t", e.IsDefined(), e.IsUndefined())
		}
	})

	t.Run("array decodes to values", func(t *testing.T) {
		var e elastic.Elastic[int]
		if err := jsonv2.Unmarshal([]byte(`[1,2,3]`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsDefined() || !slices.Equal(e.Values(), []int{1, 2, 3}) {
			t.Errorf("want defined [1 2 3], got %v", e.Values())
		}
	})

	t.Run("null array elements decode to none options", func(t *testing.T) {
		var e elastic.Elastic[int]
		if err := jsonv2.Unmarshal([]byte(`[1,null,3]`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !opt.EqualOptions(e.Options(), opt.Options[int]{opt.Some(1), opt.None[int](), opt.Some(3)}) {
			t.Errorf("want [some(1) none some(3)], got %v", e.Options())
		}
	})

	t.Run("single value decodes to single-element values", func(t *testing.T) {
		var e elastic.Elastic[string]
		if err := jsonv2.Unmarshal([]byte(`"foo"`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsDefined() || !slices.Equal(e.Values(), []string{"foo"}) {
			t.Errorf("want defined [foo], got %v", e.Values())
		}
	})

	t.Run("decode error keeps elastic unchanged", func(t *testing.T) {
		e := elastic.FromValue(7)
		if err := jsonv2.Unmarshal([]byte(`"not a number"`), &e); err == nil {
			t.Fatal("want error, got nil")
		}
		if !e.IsDefined() || e.Value() != 7 {
			t.Errorf("state changed on error: defined=%t value=%d", e.IsDefined(), e.Value())
		}
	})
}

func TestElastic_JSONv2_omitzero(t *testing.T) {
	type wrapper struct {
		A elastic.Elastic[int]    `json:"a"`
		B elastic.Elastic[string] `json:"b,omitzero"`
		C elastic.Elastic[int]    `json:"c"`
	}

	t.Run("undefined b is omitted via omitzero", func(t *testing.T) {
		got, err := jsonv2.Marshal(wrapper{A: elastic.FromValue(1), B: elastic.Undefined[string](), C: elastic.Null[int]()})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"a":[1],"c":null}`; string(got) != want {
			t.Errorf("Marshal = %s, want %s", got, want)
		}
	})

	t.Run("null b is kept", func(t *testing.T) {
		got, err := jsonv2.Marshal(wrapper{A: elastic.FromValue(1), B: elastic.Null[string](), C: elastic.FromValues(2, 3)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"a":[1],"b":null,"c":[2,3]}`; string(got) != want {
			t.Errorf("Marshal = %s, want %s", got, want)
		}
	})
}

func TestElastic_JSONv2_RoundTrip(t *testing.T) {
	type wrapper struct {
		A elastic.Elastic[int]    `json:"a"`
		B elastic.Elastic[string] `json:"b"`
	}
	in := wrapper{A: elastic.FromValues(1, 2), B: elastic.Null[string]()}
	b, err := jsonv2.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var out wrapper
	if err := jsonv2.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !elastic.Equal(out.A, in.A) {
		t.Errorf("A: got %v, want %v", out.A.Values(), in.A.Values())
	}
	if !out.B.IsNull() {
		t.Errorf("B: not null (defined=%t undefined=%t)", out.B.IsDefined(), out.B.IsUndefined())
	}
}
