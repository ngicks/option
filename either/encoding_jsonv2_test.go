package either_test

import (
	jsonv2 "encoding/json/v2"
	"errors"
	"testing"

	"github.com/ngicks/option/either"
)

func TestEither_MarshalJSONTo(t *testing.T) {
	tests := []struct {
		name string
		in   either.Either[int, string]
		want string
	}{
		{"left", either.Left[int, string](5), `{"left":5}`},
		{"right", either.Right[int]("foo"), `{"right":"foo"}`},
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

func TestEither_UnmarshalJSONFrom(t *testing.T) {
	t.Run("left key decodes to left", func(t *testing.T) {
		var e either.Either[int, string]
		if err := jsonv2.Unmarshal([]byte(`{"left":5}`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsLeft() || e.UnwrapLeft() != 5 {
			t.Errorf("want left(5), got %v", e)
		}
	})

	t.Run("right key decodes to right", func(t *testing.T) {
		var e either.Either[int, string]
		if err := jsonv2.Unmarshal([]byte(`{"right":"foo"}`), &e); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !e.IsRight() || e.UnwrapRight() != "foo" {
			t.Errorf("want right(foo), got %v", e)
		}
	})

	t.Run("invalid inputs error and keep e unchanged", func(t *testing.T) {
		for _, in := range []string{`5`, `{}`, `{"up":5}`, `{"left":5,"right":"foo"}`, `{"left":"not a number"}`} {
			e := either.Left[int, string](5)
			if err := jsonv2.Unmarshal([]byte(in), &e); err == nil {
				t.Fatalf("Unmarshal(%s): want error, got nil", in)
			}
			if !e.IsLeft() || e.UnwrapLeft() != 5 {
				t.Errorf("Unmarshal(%s): state changed on error: %v", in, e)
			}
		}
	})
}

func TestResult_JSONv2(t *testing.T) {
	t.Run("ok marshals with ok tag", func(t *testing.T) {
		got, err := jsonv2.Marshal(either.Ok(5))
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"ok":5}`; string(got) != want {
			t.Errorf("Marshal(Ok(5)) = %s, want %s", got, want)
		}
	})

	t.Run("err marshals with err tag as message string", func(t *testing.T) {
		got, err := jsonv2.Marshal(either.Err[int](errors.New("boom")))
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `{"err":"boom"}`; string(got) != want {
			t.Errorf("Marshal(Err(boom)) = %s, want %s", got, want)
		}
	})

	t.Run("round trip", func(t *testing.T) {
		for _, in := range []either.Result[int]{either.Ok(5), either.Err[int](errors.New("boom"))} {
			b, err := jsonv2.Marshal(in)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			var out either.Result[int]
			if err := jsonv2.Unmarshal(b, &out); err != nil {
				t.Fatalf("Unmarshal(%s): %v", b, err)
			}
			if out.IsOk() != in.IsOk() || out.Value() != in.Value() {
				t.Errorf("round trip: got %v, want %v", out, in)
			}
			if in.IsErr() && out.Err().Error() != in.Err().Error() {
				t.Errorf("round trip: error message %q, want %q", out.Err(), in.Err())
			}
		}
	})

	t.Run("unknown key errors", func(t *testing.T) {
		var r either.Result[int]
		if err := jsonv2.Unmarshal([]byte(`{"left":5}`), &r); err == nil {
			t.Fatal("want error, got nil")
		}
	})
}

func TestUntagged_JSONv2(t *testing.T) {
	t.Run("either marshals held value directly", func(t *testing.T) {
		got, err := jsonv2.Marshal(either.UntaggedEither[int, string]{either.Left[int, string](5)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `5`; string(got) != want {
			t.Errorf("Marshal = %s, want %s", got, want)
		}
	})

	t.Run("either unmarshals left first then right", func(t *testing.T) {
		var u either.UntaggedEither[int, string]
		if err := jsonv2.Unmarshal([]byte(`5`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsLeft() || u.UnwrapLeft() != 5 {
			t.Errorf("want left(5), got %v", u.Either)
		}
		if err := jsonv2.Unmarshal([]byte(`"foo"`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsRight() || u.UnwrapRight() != "foo" {
			t.Errorf("want right(foo), got %v", u.Either)
		}
		if err := jsonv2.Unmarshal([]byte(`{"a":1}`), &u); err == nil {
			t.Fatal("want error for value decoding as neither side, got nil")
		}
	})

	t.Run("result marshals held value directly", func(t *testing.T) {
		got, err := jsonv2.Marshal(either.UntaggedResult[int]{either.Ok(5)})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `5`; string(got) != want {
			t.Errorf("Marshal(Ok(5)) = %s, want %s", got, want)
		}
		got, err = jsonv2.Marshal(either.UntaggedResult[int]{either.Err[int](errors.New("boom"))})
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if want := `"boom"`; string(got) != want {
			t.Errorf("Marshal(Err(boom)) = %s, want %s", got, want)
		}
	})

	t.Run("result unmarshals ok first then error message", func(t *testing.T) {
		var u either.UntaggedResult[int]
		if err := jsonv2.Unmarshal([]byte(`5`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsOk() || u.Value() != 5 {
			t.Errorf("want ok(5), got %v", u.Result)
		}
		if err := jsonv2.Unmarshal([]byte(`"boom"`), &u); err != nil {
			t.Fatalf("Unmarshal: %v", err)
		}
		if !u.IsErr() || u.Err().Error() != "boom" {
			t.Errorf("want err(boom), got %v", u.Result)
		}
	})
}

func TestEither_JSONv2_RoundTrip(t *testing.T) {
	type wrapper struct {
		A either.Either[int, string]   `json:"a"`
		B either.Either[bool, []int]   `json:"b"`
		C either.Either[int, []string] `json:"c"`
	}
	in := wrapper{
		A: either.Right[int]("foo"),
		B: either.Left[bool, []int](true),
		C: either.Right[int]([]string{"a", "b"}),
	}
	b, err := jsonv2.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"a":{"right":"foo"},"b":{"left":true},"c":{"right":["a","b"]}}`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}
	var out wrapper
	if err := jsonv2.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !out.A.IsRight() || out.A.UnwrapRight() != "foo" {
		t.Errorf("A: got %v, want right(foo)", out.A)
	}
	if !out.B.IsLeft() || out.B.UnwrapLeft() != true {
		t.Errorf("B: got %v, want left(true)", out.B)
	}
	if !out.C.IsRight() || len(out.C.UnwrapRight()) != 2 {
		t.Errorf("C: got %v, want right([a b])", out.C)
	}
}
