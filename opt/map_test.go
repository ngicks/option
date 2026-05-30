package opt_test

import (
	"strconv"
	"testing"

	"github.com/ngicks/option/opt"
)

func TestOption_Map(t *testing.T) {
	// Generic method mapping across types: Option[int] -> Option[string].
	got := opt.Some(7).Map(strconv.Itoa)
	if !got.IsSome() || got.Value() != "7" {
		t.Errorf(`Some(7).Map(Itoa) = %v, want Some("7")`, got)
	}

	none := opt.None[int]().Map(strconv.Itoa)
	if none.IsSome() {
		t.Errorf("None[int]().Map(Itoa) = %v, want none", none)
	}
}

func TestOption_MapOr(t *testing.T) {
	// Generic method: Option[int] -> string, with a string default.
	if got := opt.Some(7).MapOr("default", strconv.Itoa); got != "7" {
		t.Errorf(`Some(7).MapOr("default", Itoa) = %q, want "7"`, got)
	}

	if got := opt.None[int]().MapOr("default", strconv.Itoa); got != "default" {
		t.Errorf(`None[int]().MapOr("default", Itoa) = %q, want "default"`, got)
	}
}
