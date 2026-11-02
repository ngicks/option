package und_test

import (
	"slices"
	"testing"

	"github.com/ngicks/option/und"
)

func TestUnd_States(t *testing.T) {
	d := und.Defined(5)
	if !d.IsDefined() || d.IsNull() || d.IsUndefined() || d.IsZero() {
		t.Errorf(
			"Defined(5): defined=%t null=%t undefined=%t zero=%t",
			d.IsDefined(),
			d.IsNull(),
			d.IsUndefined(),
			d.IsZero(),
		)
	}
	if d.Value() != 5 {
		t.Errorf("Defined(5).Value() = %d, want 5", d.Value())
	}

	n := und.Null[int]()
	if n.IsDefined() || !n.IsNull() || n.IsUndefined() {
		t.Errorf(
			"Null: defined=%t null=%t undefined=%t",
			n.IsDefined(),
			n.IsNull(),
			n.IsUndefined(),
		)
	}

	u := und.Undefined[int]()
	if u.IsDefined() || u.IsNull() || !u.IsUndefined() || !u.IsZero() {
		t.Errorf(
			"Undefined: defined=%t null=%t undefined=%t zero=%t",
			u.IsDefined(),
			u.IsNull(),
			u.IsUndefined(),
			u.IsZero(),
		)
	}
	// Zero value of Und[T] is undefined.
	var zero und.Und[int]
	if !zero.IsUndefined() {
		t.Error("zero Und is not undefined")
	}
}

func TestUnd_FromPointer(t *testing.T) {
	v := 5
	if u := und.FromPointer(&v); !u.IsDefined() || u.Value() != 5 {
		t.Errorf("FromPointer(&5) = %v", u)
	}
	if u := und.FromPointer[int](nil); !u.IsUndefined() {
		t.Error("FromPointer(nil) is not undefined")
	}
}

func TestUnd_Pointer(t *testing.T) {
	// Pointer returns a pointer to the inner opt.Option[T]; nil only when undefined.
	if und.Undefined[int]().Pointer() != nil {
		t.Error("Undefined.Pointer() != nil")
	}
	if p := und.Null[int]().Pointer(); p == nil || p.IsSome() {
		t.Errorf("Null.Pointer() = %v, want pointer to a none option", p)
	}
	if p := und.Defined(5).Pointer(); p == nil || !p.IsSome() || p.Value() != 5 {
		t.Errorf("Defined(5).Pointer() = %v, want pointer to some(5)", p)
	}
}

func TestUnd_Strip(t *testing.T) {
	// Strip exposes the backing option so Option's combinators can be used.
	if o := und.Defined(5).Strip(); !o.IsSome() || !o.Value().IsSome() || o.Value().Value() != 5 {
		t.Errorf("Defined(5).Strip() = %v, want some(some(5))", o)
	}
	if o := und.Null[int]().Strip(); !o.IsSome() || o.Value().IsSome() {
		t.Errorf("Null.Strip() = %v, want some(none)", o)
	}
	if o := und.Undefined[int]().Strip(); o.IsSome() {
		t.Errorf("Undefined.Strip() = %v, want none", o)
	}
}

func TestUnd_Equal(t *testing.T) {
	cases := []struct {
		name string
		l, r und.Und[int]
		want bool
	}{
		{"defined eq", und.Defined(5), und.Defined(5), true},
		{"defined ne", und.Defined(5), und.Defined(6), false},
		{"null eq", und.Null[int](), und.Null[int](), true},
		{"undefined eq", und.Undefined[int](), und.Undefined[int](), true},
		{"defined vs null", und.Defined(5), und.Null[int](), false},
		{"null vs undefined", und.Null[int](), und.Undefined[int](), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := und.Equal(c.l, c.r); got != c.want {
				t.Errorf("Equal(%v, %v) = %t, want %t", c.l, c.r, got, c.want)
			}
		})
	}
}

func TestUnd_CloneFunc(t *testing.T) {
	// CloneFunc deep-clones the inner value via the supplied function.
	orig := und.Defined([]int{1, 2, 3})
	cloned := orig.CloneFunc(slices.Clone)
	orig.Value()[0] = 99
	if got := cloned.Value(); got[0] != 1 {
		t.Errorf("CloneFunc did not deep copy: cloned[0] = %d, want 1", got[0])
	}
}
