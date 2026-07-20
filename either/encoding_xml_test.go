package either_test

import (
	"encoding/xml"
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/ngicks/option/either"
)

func TestEither_XML(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name                   `xml:"root"`
		V       either.Either[int, string] `xml:"v"`
	}

	t.Run("marshals tagged elements", func(t *testing.T) {
		for _, tt := range []struct {
			in   either.Either[int, string]
			want string
		}{
			{either.Left[int, string](5), `<root><v><left>5</left></v></root>`},
			{either.Right[int]("foo"), `<root><v><right>foo</right></v></root>`},
		} {
			b, err := xml.Marshal(wrapper{V: tt.in})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(b) != tt.want {
				t.Errorf("Marshal(%v) = %s, want %s", tt.in, b, tt.want)
			}
		}
	})

	t.Run("round trip", func(t *testing.T) {
		for _, in := range []either.Either[int, string]{either.Left[int, string](5), either.Right[int]("foo")} {
			b, err := xml.Marshal(wrapper{V: in})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			var out wrapper
			if err := xml.Unmarshal(b, &out); err != nil {
				t.Fatalf("Unmarshal(%s): %v", b, err)
			}
			if !either.Equal(out.V, in) {
				t.Errorf("round trip: got %v, want %v", out.V, in)
			}
		}
	})

	t.Run("invalid inputs error and keep e unchanged", func(t *testing.T) {
		for _, in := range []string{
			`<root><v></v></root>`,
			`<root><v><up>5</up></v></root>`,
			`<root><v><left>5</left><right>foo</right></v></root>`,
			`<root><v><left>not a number</left></v></root>`,
		} {
			out := wrapper{V: either.Left[int, string](5)}
			if err := xml.Unmarshal([]byte(in), &out); err == nil {
				t.Fatalf("Unmarshal(%s): want error, got nil", in)
			}
			if !out.V.IsLeft() || out.V.UnwrapLeft() != 5 {
				t.Errorf("Unmarshal(%s): state changed on error: %v", in, out.V)
			}
		}
	})
}

// intOnly decodes only from plain integer character data,
// unlike int which encoding/xml leniently zeroes on empty content.
type intOnly int

func (i *intOnly) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return err
	}
	*i = intOnly(v)
	return nil
}

func TestUntagged_XML(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name                           `xml:"root"`
		V       either.UntaggedEither[int, string] `xml:"v"`
	}
	type resultWrapper struct {
		XMLName xml.Name                   `xml:"root"`
		V       either.UntaggedResult[int] `xml:"v"`
	}

	t.Run("marshals held value directly", func(t *testing.T) {
		b, err := xml.Marshal(wrapper{V: either.UntaggedEither[int, string]{either.Left[int, string](5)}})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `<root><v>5</v></root>`; string(b) != want {
			t.Errorf("Marshal = %s, want %s", b, want)
		}

		b, err = xml.Marshal(resultWrapper{V: either.UntaggedResult[int]{either.Err[int](errors.New("boom"))}})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `<root><v>boom</v></root>`; string(b) != want {
			t.Errorf("Marshal = %s, want %s", b, want)
		}
	})

	t.Run("either unmarshals left first then right", func(t *testing.T) {
		var out wrapper
		if err := xml.Unmarshal([]byte(`<root><v>5</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsLeft() || out.V.UnwrapLeft() != 5 {
			t.Errorf("want left(5), got %v", out.V.Either)
		}
		if err := xml.Unmarshal([]byte(`<root><v>foo</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsRight() || out.V.UnwrapRight() != "foo" {
			t.Errorf("want right(foo), got %v", out.V.Either)
		}
	})

	t.Run("nested elements replay", func(t *testing.T) {
		// Plain int would win here: encoding/xml leniently decodes an element
		// without character data into an int as 0. Use the strict intOnly to
		// exercise the fallback replay of nested tokens.
		type point struct {
			X int `xml:"x"`
			Y int `xml:"y"`
		}
		type pointWrapper struct {
			XMLName xml.Name                              `xml:"root"`
			V       either.UntaggedEither[intOnly, point] `xml:"v"`
		}
		var out pointWrapper
		if err := xml.Unmarshal([]byte(`<root><v>5</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsLeft() || out.V.UnwrapLeft() != 5 {
			t.Errorf("want left(5), got %v", out.V.Either)
		}
		if err := xml.Unmarshal([]byte(`<root><v><x>1</x><y>2</y></v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsRight() || out.V.UnwrapRight() != (point{X: 1, Y: 2}) {
			t.Errorf("want right({1 2}), got %v", out.V.Either)
		}
	})

	t.Run("result unmarshals ok first then error message", func(t *testing.T) {
		var out resultWrapper
		if err := xml.Unmarshal([]byte(`<root><v>5</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsOk() || out.V.Value() != 5 {
			t.Errorf("want ok(5), got %v", out.V.Result)
		}
		if err := xml.Unmarshal([]byte(`<root><v>boom</v></root>`), &out); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !out.V.IsErr() || out.V.Err().Error() != "boom" {
			t.Errorf("want err(boom), got %v", out.V.Result)
		}
	})

	t.Run("neither side errors and keeps value unchanged", func(t *testing.T) {
		type strictWrapper struct {
			XMLName xml.Name                         `xml:"root"`
			V       either.UntaggedEither[int, bool] `xml:"v"`
		}
		out := strictWrapper{V: either.UntaggedEither[int, bool]{either.Left[int, bool](5)}}
		if err := xml.Unmarshal([]byte(`<root><v>foo</v></root>`), &out); err == nil {
			t.Fatal("want error, got nil")
		}
		if !out.V.IsLeft() || out.V.UnwrapLeft() != 5 {
			t.Errorf("state changed on error: %v", out.V.Either)
		}
	})

	t.Run("round trip", func(t *testing.T) {
		in := wrapper{V: either.UntaggedEither[int, string]{either.Right[int]("foo")}}
		b, err := xml.Marshal(in)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		var out wrapper
		if err := xml.Unmarshal(b, &out); err != nil {
			t.Fatalf("Unmarshal(%s): %v", b, err)
		}
		if !either.Equal(out.V.Either, in.V.Either) {
			t.Errorf("round trip: got %v, want %v", out.V.Either, in.V.Either)
		}
	})
}

func TestResult_XML(t *testing.T) {
	type wrapper struct {
		XMLName xml.Name           `xml:"root"`
		V       either.Result[int] `xml:"v"`
	}

	t.Run("marshals tagged elements", func(t *testing.T) {
		for _, tt := range []struct {
			in   either.Result[int]
			want string
		}{
			{either.Ok(5), `<root><v><ok>5</ok></v></root>`},
			{either.Err[int](errors.New("boom")), `<root><v><err>boom</err></v></root>`},
		} {
			b, err := xml.Marshal(wrapper{V: tt.in})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(b) != tt.want {
				t.Errorf("Marshal = %s, want %s", b, tt.want)
			}
		}
	})

	t.Run("round trip", func(t *testing.T) {
		for _, in := range []either.Result[int]{either.Ok(5), either.Err[int](errors.New("boom"))} {
			b, err := xml.Marshal(wrapper{V: in})
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			var out wrapper
			if err := xml.Unmarshal(b, &out); err != nil {
				t.Fatalf("Unmarshal(%s): %v", b, err)
			}
			if out.V.IsOk() != in.IsOk() || out.V.Value() != in.Value() {
				t.Errorf("round trip: got %v, want %v", out.V, in)
			}
			if in.IsErr() && out.V.Err().Error() != in.Err().Error() {
				t.Errorf("round trip: error message %q, want %q", out.V.Err(), in.Err())
			}
		}
	})
}
