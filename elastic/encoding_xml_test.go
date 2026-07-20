package elastic_test

import (
	"encoding/xml"
	"slices"
	"testing"

	"github.com/ngicks/option/elastic"
	"github.com/ngicks/option/opt"
)

func TestElastic_XML(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name             `xml:"root"`
		V       elastic.Elastic[int] `xml:"v"`
	}

	t.Run("defined marshals one element per value", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: elastic.FromValues(1, 2, 3)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `<root><v>1</v><v>2</v><v>3</v></root>`; string(b) != want {
			t.Errorf("Marshal = %s, want %s", b, want)
		}
	})

	t.Run("none options encode nothing", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: elastic.FromOptions(opt.Some(1), opt.None[int](), opt.Some(2))})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `<root><v>1</v><v>2</v></root>`; string(b) != want {
			t.Errorf("Marshal = %s, want %s", b, want)
		}
	})

	t.Run("non-defined omits the element", func(t *testing.T) {
		for _, e := range []elastic.Elastic[int]{elastic.Null[int](), elastic.Undefined[int]()} {
			b, err := xml.Marshal(wrapper{V: e})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if want := `<root></root>`; string(b) != want {
				t.Errorf("Marshal(defined=%t) = %s, want %s", e.IsDefined(), b, want)
			}
		}
	})

	t.Run("repeated elements accumulate", func(t *testing.T) {
		var out wrapper
		if err := xml.Unmarshal([]byte(`<root><v>1</v><v>2</v><v>3</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsDefined() || !slices.Equal(out.V.Values(), []int{1, 2, 3}) {
			t.Errorf("Unmarshal: got %v, want [1 2 3]", out.V.Values())
		}
	})

	t.Run("round trip", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: elastic.FromValues(4, 5)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var out wrapper
		if err := xml.Unmarshal(b, &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !slices.Equal(out.V.Values(), []int{4, 5}) {
			t.Errorf("round trip: got %v, want [4 5]", out.V.Values())
		}
	})
}
