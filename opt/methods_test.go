package opt_test

import (
	"strconv"
	"testing"

	"github.com/ngicks/option/opt"
)

func TestOption_Expect(t *testing.T) {
	if got := opt.Some(7).Expect("boom"); got != 7 {
		t.Errorf("Some(7).Expect = %d, want 7", got)
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("None.Expect did not panic")
		}
		if msg, _ := r.(string); msg != "boom" {
			t.Errorf("panic value = %v, want %q", r, "boom")
		}
	}()
	opt.None[int]().Expect("boom")
}

func TestOption_Unwrap(t *testing.T) {
	if got := opt.Some(7).Unwrap(); got != 7 {
		t.Errorf("Some(7).Unwrap = %d, want 7", got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("None.Unwrap did not panic")
		}
	}()
	opt.None[int]().Unwrap()
}

func TestOption_UnwrapOr(t *testing.T) {
	if got := opt.Some(7).UnwrapOr(3); got != 7 {
		t.Errorf("Some(7).UnwrapOr(3) = %d, want 7", got)
	}
	if got := opt.None[int]().UnwrapOr(3); got != 3 {
		t.Errorf("None.UnwrapOr(3) = %d, want 3", got)
	}
}

func TestOption_UnwrapOrDefault(t *testing.T) {
	if got := opt.Some(7).UnwrapOrDefault(); got != 7 {
		t.Errorf("Some(7).UnwrapOrDefault = %d, want 7", got)
	}
	if got := opt.None[int]().UnwrapOrDefault(); got != 0 {
		t.Errorf("None.UnwrapOrDefault = %d, want 0", got)
	}
}

func TestOption_UnwrapOrElse(t *testing.T) {
	if got := opt.Some(7).UnwrapOrElse(func() int { return 3 }); got != 7 {
		t.Errorf("Some(7).UnwrapOrElse = %d, want 7", got)
	}
	called := false
	got := opt.None[int]().UnwrapOrElse(func() int {
		called = true
		return 3
	})
	if got != 3 || !called {
		t.Errorf("None.UnwrapOrElse = %d (called=%t), want 3 (called=true)", got, called)
	}
}

func TestOption_Insert(t *testing.T) {
	o := opt.Some(1)
	p := o.Insert(2)
	if !o.IsSome() || o.Value() != 2 {
		t.Errorf("after Insert, o = %v, want Some(2)", o)
	}
	if *p != 2 {
		t.Errorf("Insert returned pointer to %d, want 2", *p)
	}
	// The returned pointer aliases the internal value.
	*p = 9
	if o.Value() != 9 {
		t.Errorf("after *p = 9, o.Value() = %d, want 9", o.Value())
	}
}

func TestOption_GetOrInsert(t *testing.T) {
	none := opt.None[int]()
	if p := none.GetOrInsert(5); *p != 5 || !none.IsSome() || none.Value() != 5 {
		t.Errorf("None.GetOrInsert(5): o = %v, *p = %d", none, *p)
	}

	some := opt.Some(1)
	if p := some.GetOrInsert(5); *p != 1 || some.Value() != 1 {
		t.Errorf("Some(1).GetOrInsert(5): o = %v, *p = %d, want kept 1", some, *p)
	}
}

func TestOption_GetOrInsertDefault(t *testing.T) {
	none := opt.None[int]()
	if p := none.GetOrInsertDefault(); *p != 0 || !none.IsSome() {
		t.Errorf("None.GetOrInsertDefault: o = %v, *p = %d", none, *p)
	}

	some := opt.Some(7)
	if p := some.GetOrInsertDefault(); *p != 7 || some.Value() != 7 {
		t.Errorf("Some(7).GetOrInsertDefault: o = %v, *p = %d, want kept 7", some, *p)
	}
}

func TestOption_GetOrInsertFunc(t *testing.T) {
	none := opt.None[int]()
	if p := none.GetOrInsertFunc(func() int { return 5 }); *p != 5 || none.Value() != 5 {
		t.Errorf("None.GetOrInsertFunc: o = %v, *p = %d", none, *p)
	}

	some := opt.Some(1)
	called := false
	p := some.GetOrInsertFunc(func() int {
		called = true
		return 5
	})
	if *p != 1 || called {
		t.Errorf("Some(1).GetOrInsertFunc: *p = %d, called = %t, want 1, false", *p, called)
	}
}

func TestOption_Inspect(t *testing.T) {
	seen := 0
	got := opt.Some(7).Inspect(func(v int) { seen = v })
	if seen != 7 {
		t.Errorf("Inspect did not call fn with 7, got %d", seen)
	}
	if !got.IsSome() || got.Value() != 7 {
		t.Errorf("Inspect returned %v, want Some(7)", got)
	}

	seen = 0
	opt.None[int]().Inspect(func(v int) { seen = v })
	if seen != 0 {
		t.Errorf("None.Inspect called fn, seen = %d", seen)
	}
}

func TestOption_MapOrElse(t *testing.T) {
	// Regression: MapOrElse used to wrap the result in Some, breaking the return type.
	if got := opt.Some(7).MapOrElse(func() string { return "default" }, strconv.Itoa); got != "7" {
		t.Errorf(`Some(7).MapOrElse = %q, want "7"`, got)
	}
	if got := opt.None[int]().MapOrElse(func() string { return "default" }, strconv.Itoa); got != "default" {
		t.Errorf(`None.MapOrElse = %q, want "default"`, got)
	}
}

func TestOption_Filter(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	if got := opt.Some(4).Filter(isEven); !got.IsSome() || got.Value() != 4 {
		t.Errorf("Some(4).Filter(isEven) = %v, want Some(4)", got)
	}
	if got := opt.Some(3).Filter(isEven); got.IsSome() {
		t.Errorf("Some(3).Filter(isEven) = %v, want none", got)
	}
	if got := opt.None[int]().Filter(isEven); got.IsSome() {
		t.Errorf("None.Filter(isEven) = %v, want none", got)
	}
}

func TestOption_Iter(t *testing.T) {
	var collected []int
	for v := range opt.Some(7).Iter() {
		collected = append(collected, v)
	}
	if len(collected) != 1 || collected[0] != 7 {
		t.Errorf("Some(7).Iter yielded %v, want [7]", collected)
	}

	collected = nil
	for v := range opt.None[int]().Iter() {
		collected = append(collected, v)
	}
	if len(collected) != 0 {
		t.Errorf("None.Iter yielded %v, want nothing", collected)
	}
}
