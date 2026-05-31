package und_test

import (
	"encoding/xml"
	"testing"

	"github.com/ngicks/option/und"
)

func TestUnd_XML(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name     `xml:"root"`
		V       und.Und[int] `xml:"v"`
	}

	t.Run("defined marshals the element", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: und.Defined(5)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `<root><v>5</v></root>`; string(b) != want {
			t.Errorf("Marshal = %s, want %s", b, want)
		}
	})

	t.Run("non-defined omits the element", func(t *testing.T) {
		for _, u := range []und.Und[int]{und.Null[int](), und.Undefined[int]()} {
			b, err := xml.Marshal(wrapper{V: u})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if want := `<root></root>`; string(b) != want {
				t.Errorf("Marshal(state %d) = %s, want %s", u.State(), b, want)
			}
		}
	})

	t.Run("round trip defined", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: und.Defined(5)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var out wrapper
		if err := xml.Unmarshal(b, &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsDefined() || out.V.Value() != 5 {
			t.Errorf("round trip: got %v, want defined(5)", out.V)
		}
	})
}
