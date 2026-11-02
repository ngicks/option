package elastic_test

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/ngicks/option/elastic"
	"github.com/ngicks/option/opt"
)

func TestElastic_MarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   elastic.Elastic[int]
		want string
	}{
		{"multiple values", elastic.FromValues(1, 2, 3), `[1,2,3]`},
		{"single value", elastic.FromValue(5), `[5]`},
		{
			"none options",
			elastic.FromOptions(opt.Some(1), opt.None[int](), opt.Some(3)),
			`[1,null,3]`,
		},
		{"empty", elastic.FromValues[int](), `[]`},
		{"null", elastic.Null[int](), `null`},
		{"undefined", elastic.Undefined[int](), `null`},
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

func TestElastic_UnmarshalJSON(t *testing.T) {
	t.Run("null decodes to null state", func(t *testing.T) {
		e := elastic.FromValue(5)
		if err := json.Unmarshal([]byte(`null`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsNull() {
			t.Errorf("want null, got defined=%t undefined=%t", e.IsDefined(), e.IsUndefined())
		}
	})

	t.Run("array decodes to values", func(t *testing.T) {
		var e elastic.Elastic[int]
		if err := json.Unmarshal([]byte(`[1,2,3]`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsDefined() || !slices.Equal(e.Values(), []int{1, 2, 3}) {
			t.Errorf("want defined [1 2 3], got %v", e.Values())
		}
	})

	t.Run("null array elements decode to none options", func(t *testing.T) {
		var e elastic.Elastic[int]
		if err := json.Unmarshal([]byte(`[1,null,3]`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !opt.EqualOptions(
			e.Options(),
			opt.Options[int]{opt.Some(1), opt.None[int](), opt.Some(3)},
		) {
			t.Errorf("want [some(1) none some(3)], got %v", e.Options())
		}
	})

	t.Run("single value decodes to single-element values", func(t *testing.T) {
		var e elastic.Elastic[int]
		if err := json.Unmarshal([]byte(`5`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsDefined() || !slices.Equal(e.Values(), []int{5}) {
			t.Errorf("want defined [5], got %v", e.Values())
		}
	})

	t.Run("decode error keeps elastic unchanged", func(t *testing.T) {
		e := elastic.FromValue(7)
		if err := json.Unmarshal([]byte(`"not a number"`), &e); err == nil {
			t.Fatal("want error, got nil")
		}
		if !e.IsDefined() || e.Value() != 7 {
			t.Errorf("state changed on error: defined=%t value=%d", e.IsDefined(), e.Value())
		}
	})
}
