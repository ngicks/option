package und_test

import (
	"slices"
	"strconv"
	"testing"

	"github.com/ngicks/option/und"
)

func TestUnd_States(t *testing.T) {
	d := und.Defined(5)
	if !d.IsDefined() || d.IsNull() || d.IsUndefined() || d.IsZero() {
		t.Errorf("Defined(5): defined=%t null=%t undefined=%t zero=%t", d.IsDefined(), d.IsNull(), d.IsUndefined(), d.IsZero())
	}
	if d.State() != und.StateDefined {
		t.Errorf("Defined(5).State() = %d, want StateDefined", d.State())
	}
	if d.Value() != 5 {
		t.Errorf("Defined(5).Value() = %d, want 5", d.Value())
	}

	n := und.Null[int]()
	if n.IsDefined() || !n.IsNull() || n.IsUndefined() {
		t.Errorf("Null: defined=%t null=%t undefined=%t", n.IsDefined(), n.IsNull(), n.IsUndefined())
	}
	if n.State() != und.StateNull {
		t.Errorf("Null.State() = %d, want StateNull", n.State())
	}

	u := und.Undefined[int]()
	if u.IsDefined() || u.IsNull() || !u.IsUndefined() || !u.IsZero() {
		t.Errorf("Undefined: defined=%t null=%t undefined=%t zero=%t", u.IsDefined(), u.IsNull(), u.IsUndefined(), u.IsZero())
	}
	if u.State() != und.StateUndefined {
		t.Errorf("Undefined.State() = %d, want StateUndefined", u.State())
	}
	// Zero value of Und[T] is undefined.
	var zero und.Und[int]
	if !zero.IsUndefined() {
		t.Errorf("zero Und is not undefined: %d", zero.State())
	}
}

func TestUnd_FromPointer(t *testing.T) {
	v := 5
	if u := und.FromPointer(&v); !u.IsDefined() || u.Value() != 5 {
		t.Errorf("FromPointer(&5) = %v", u)
	}
	if u := und.FromPointer[int](nil); !u.IsUndefined() {
		t.Errorf("FromPointer(nil) = %v, want undefined", u.State())
	}
}

func TestUnd_Pointer(t *testing.T) {
	d := und.Defined(5)
	p := d.Pointer()
	if p == nil || *p != 5 {
		t.Fatalf("Defined(5).Pointer() = %v", p)
	}
	// The returned pointer is a copy; mutating it must not affect d.
	*p = 9
	if d.Value() != 5 {
		t.Errorf("Pointer aliases internal value: d.Value() = %d, want 5", d.Value())
	}
	if und.Null[int]().Pointer() != nil {
		t.Error("Null.Pointer() != nil")
	}
	if und.Undefined[int]().Pointer() != nil {
		t.Error("Undefined.Pointer() != nil")
	}
}

func TestUnd_DoublePointer(t *testing.T) {
	if und.Undefined[int]().DoublePointer() != nil {
		t.Error("Undefined.DoublePointer() != nil")
	}
	if pp := und.Null[int]().DoublePointer(); pp == nil || *pp != nil {
		t.Errorf("Null.DoublePointer() = %v, want &(*int)(nil)", pp)
	}
	if pp := und.Defined(5).DoublePointer(); pp == nil || *pp == nil || **pp != 5 {
		t.Errorf("Defined(5).DoublePointer() = %v, want **5", pp)
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

func TestUnd_Map(t *testing.T) {
	if got := und.Map(und.Defined(5), strconv.Itoa); !got.IsDefined() || got.Value() != "5" {
		t.Errorf("Map(Defined(5), Itoa) = %v, want Defined(\"5\")", got)
	}
	if got := und.Map(und.Null[int](), strconv.Itoa); !got.IsNull() {
		t.Errorf("Map(Null, Itoa) = %v, want null", got.State())
	}
	if got := und.Map(und.Undefined[int](), strconv.Itoa); !got.IsUndefined() {
		t.Errorf("Map(Undefined, Itoa) = %v, want undefined", got.State())
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
