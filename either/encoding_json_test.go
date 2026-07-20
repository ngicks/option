package either_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ngicks/option/either"
)

func TestEither_JSONv1_RoundTrip(t *testing.T) {
	type wrapper struct {
		A either.Either[int, string] `json:"a"`
		B either.Either[int, string] `json:"b"`
	}
	in := wrapper{A: either.Left[int, string](5), B: either.Right[int]("foo")}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `{"a":{"left":5},"b":{"right":"foo"}}`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}
	var out wrapper
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !out.A.IsLeft() || out.A.UnwrapLeft() != 5 {
		t.Errorf("A: got %v, want left(5)", out.A)
	}
	if !out.B.IsRight() || out.B.UnwrapRight() != "foo" {
		t.Errorf("B: got %v, want right(foo)", out.B)
	}
}

func TestResult_JSONv1_RoundTrip(t *testing.T) {
	for _, tt := range []struct {
		in   either.Result[int]
		want string
	}{
		{either.Ok(5), `{"ok":5}`},
		{either.Err[int](errors.New("boom")), `{"err":"boom"}`},
	} {
		b, err := json.Marshal(tt.in)
		if err != nil {
			t.Fatalf("Marshal: %v", err)
		}
		if string(b) != tt.want {
			t.Errorf("Marshal = %s, want %s", b, tt.want)
		}
		var out either.Result[int]
		if err := json.Unmarshal(b, &out); err != nil {
			t.Fatalf("Unmarshal(%s): %v", b, err)
		}
		if out.IsOk() != tt.in.IsOk() || out.Value() != tt.in.Value() {
			t.Errorf("round trip: got %v, want %v", out, tt.in)
		}
		if tt.in.IsErr() && out.Err().Error() != tt.in.Err().Error() {
			t.Errorf("round trip: error message %q, want %q", out.Err(), tt.in.Err())
		}
	}
}

func TestUntagged_JSONv1(t *testing.T) {
	b, err := json.Marshal(either.UntaggedEither[int, string]{either.Right[int]("foo")})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `"foo"`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}
	var ue either.UntaggedEither[int, string]
	if err := json.Unmarshal([]byte(`5`), &ue); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !ue.IsLeft() || ue.UnwrapLeft() != 5 {
		t.Errorf("want left(5), got %v", ue.Either)
	}

	b, err = json.Marshal(either.UntaggedResult[int]{either.Err[int](errors.New("boom"))})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if want := `"boom"`; string(b) != want {
		t.Errorf("Marshal = %s, want %s", b, want)
	}
	var ur either.UntaggedResult[int]
	if err := json.Unmarshal(b, &ur); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !ur.IsErr() || ur.Err().Error() != "boom" {
		t.Errorf("want err(boom), got %v", ur.Result)
	}
}
