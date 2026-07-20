package elastic_test

import (
	"slices"
	"testing"

	"github.com/ngicks/option/elastic"
	"github.com/ngicks/option/opt"
	"github.com/ngicks/option/und"
)

func TestElastic_States(t *testing.T) {
	d := elastic.FromValues(1, 2, 3)
	if !d.IsDefined() || d.IsNull() || d.IsUndefined() || d.IsZero() {
		t.Errorf("FromValues: defined=%t null=%t undefined=%t zero=%t", d.IsDefined(), d.IsNull(), d.IsUndefined(), d.IsZero())
	}
	if d.Value() != 1 {
		t.Errorf("FromValues(1,2,3).Value() = %d, want 1", d.Value())
	}
	if !slices.Equal(d.Values(), []int{1, 2, 3}) {
		t.Errorf("FromValues(1,2,3).Values() = %v", d.Values())
	}

	s := elastic.FromValue(5)
	if !s.IsDefined() || s.Value() != 5 || !slices.Equal(s.Values(), []int{5}) {
		t.Errorf("FromValue(5): defined=%t value=%d values=%v", s.IsDefined(), s.Value(), s.Values())
	}

	n := elastic.Null[int]()
	if n.IsDefined() || !n.IsNull() || n.IsUndefined() {
		t.Errorf("Null: defined=%t null=%t undefined=%t", n.IsDefined(), n.IsNull(), n.IsUndefined())
	}

	u := elastic.Undefined[int]()
	if u.IsDefined() || u.IsNull() || !u.IsUndefined() || !u.IsZero() {
		t.Errorf("Undefined: defined=%t null=%t undefined=%t zero=%t", u.IsDefined(), u.IsNull(), u.IsUndefined(), u.IsZero())
	}
	// Zero value of Elastic[T] is undefined.
	var zero elastic.Elastic[int]
	if !zero.IsUndefined() {
		t.Error("zero Elastic is not undefined")
	}
}

func TestElastic_Pointer(t *testing.T) {
	if elastic.Undefined[int]().Pointer() != nil {
		t.Error("Undefined.Pointer() != nil")
	}
	if elastic.FromValues[int]().Pointer() != nil {
		t.Error("empty FromValues().Pointer() != nil")
	}
	if p := elastic.FromValues(1, 2).Pointer(); p == nil || *p != 1 {
		t.Errorf("FromValues(1,2).Pointer() = %v, want pointer to 1", p)
	}
	if ps := elastic.FromValues(1, 2).Pointers(); len(ps) != 2 || *ps[0] != 1 || *ps[1] != 2 {
		t.Errorf("FromValues(1,2).Pointers() = %v", ps)
	}
	if elastic.Null[int]().Pointers() != nil {
		t.Error("Null.Pointers() != nil")
	}
}

func TestElastic_Conversion(t *testing.T) {
	if e := elastic.FromPointer[int](nil); !e.IsUndefined() {
		t.Error("FromPointer(nil) is not undefined")
	}
	v := 5
	if e := elastic.FromPointer(&v); !e.IsDefined() || e.Value() != 5 {
		t.Errorf("FromPointer(&5) = %v", e)
	}

	if e := elastic.FromUnd(und.Defined(5)); !e.IsDefined() || e.Value() != 5 {
		t.Errorf("FromUnd(Defined(5)) = %v", e)
	}
	if e := elastic.FromUnd(und.Null[int]()); !e.IsNull() {
		t.Error("FromUnd(Null) is not null")
	}
	if e := elastic.FromUnd(und.Undefined[int]()); !e.IsUndefined() {
		t.Error("FromUnd(Undefined) is not undefined")
	}

	if u := elastic.FromValues(1, 2).Und(); !u.IsDefined() || u.Value() != 1 {
		t.Errorf("FromValues(1,2).Und() = %v, want defined(1)", u)
	}
	if u := elastic.Null[int]().Und(); !u.IsNull() {
		t.Error("Null.Und() is not null")
	}
	if u := elastic.FromOptions(opt.None[int]()).Und(); !u.IsNull() {
		t.Error("FromOptions(none).Und() is not null")
	}
	if u := elastic.FromOptions[int]().Und(); !u.IsNull() {
		t.Error("FromOptions().Und() is not null")
	}

	// Wrap and Strip are inverses.
	u := und.Defined(opt.Options[int]{opt.Some(1), opt.None[int]()})
	if got := elastic.Wrap[int](u).Strip(); !got.EqualFunc(u, opt.EqualOptions) {
		t.Errorf("Wrap(u).Strip() = %v, want %v", got, u)
	}
}

func TestElastic_Options(t *testing.T) {
	e := elastic.FromOptions(opt.Some(1), opt.None[int](), opt.Some(3))
	if !e.IsDefined() {
		t.Error("FromOptions is not defined")
	}
	if !opt.EqualOptions(e.Options(), opt.Options[int]{opt.Some(1), opt.None[int](), opt.Some(3)}) {
		t.Errorf("Options() = %v", e.Options())
	}
	// None options surface as zero values in Values and nils in Pointers.
	if !slices.Equal(e.Values(), []int{1, 0, 3}) {
		t.Errorf("Values() = %v, want [1 0 3]", e.Values())
	}
	ps := e.Pointers()
	if len(ps) != 3 || *ps[0] != 1 || ps[1] != nil || *ps[2] != 3 {
		t.Errorf("Pointers() = %v, want [&1 nil &3]", ps)
	}

	// A none first option yields the zero value / nil pointer.
	n := elastic.FromOptions(opt.None[int](), opt.Some(2))
	if n.Value() != 0 || n.Pointer() != nil {
		t.Errorf("none-first: Value() = %d, Pointer() = %v", n.Value(), n.Pointer())
	}

	one := 1
	if got := elastic.FromPointers(&one, nil); !opt.EqualOptions(got.Options(), opt.Options[int]{opt.Some(1), opt.None[int]()}) {
		t.Errorf("FromPointers(&1, nil).Options() = %v", got.Options())
	}
}

func TestElastic_Equal(t *testing.T) {
	cases := []struct {
		name string
		l, r elastic.Elastic[int]
		want bool
	}{
		{"defined eq", elastic.FromValues(1, 2), elastic.FromValues(1, 2), true},
		{"defined ne value", elastic.FromValues(1, 2), elastic.FromValues(1, 3), false},
		{"defined ne length", elastic.FromValues(1, 2), elastic.FromValues(1), false},
		{"null eq", elastic.Null[int](), elastic.Null[int](), true},
		{"undefined eq", elastic.Undefined[int](), elastic.Undefined[int](), true},
		{"defined vs null", elastic.FromValues(1), elastic.Null[int](), false},
		{"null vs undefined", elastic.Null[int](), elastic.Undefined[int](), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := elastic.Equal(c.l, c.r); got != c.want {
				t.Errorf("Equal(%v, %v) = %t, want %t", c.l, c.r, got, c.want)
			}
		})
	}
}

func TestElastic_CloneFunc(t *testing.T) {
	orig := elastic.FromValues([]int{1}, []int{2})
	cloned := orig.CloneFunc(slices.Clone)
	orig.Values()[0][0] = 99
	if got := cloned.Values(); got[0][0] != 1 {
		t.Errorf("CloneFunc did not deep copy: cloned[0][0] = %d, want 1", got[0][0])
	}
}

func TestElastic_Coalesce(t *testing.T) {
	if e := elastic.Undefined[int]().Coalesce(true); !e.IsNull() {
		t.Error("Undefined.Coalesce(true) is not null")
	}
	if e := elastic.Null[int]().Coalesce(false); !e.IsUndefined() {
		t.Error("Null.Coalesce(false) is not undefined")
	}
	if e := elastic.FromValue(5).Coalesce(true); !e.IsDefined() || e.Value() != 5 {
		t.Errorf("FromValue(5).Coalesce(true) = %v, want defined(5)", e)
	}
}
