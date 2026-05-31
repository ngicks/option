package und

import (
	"database/sql"

	"github.com/ngicks/option/opt"
)

// Und[T] is a type that can express T (a value of type T), *null* (exists but empty), or *undefined* (absent, unspecified).
//
// Und[T] implements json.Unmarshaler so that it can be unmarshaled from all of those type.
//
// Und[T] implements IsZero.
// *undefined* Und[T] struct fields are omitted by [json.Marshal] (or similar functions)
// if `json:",omitzero"` option is attached to those fields.
type Und[T any] struct {
	opt opt.Option[opt.Option[T]]
}

// Defined returns a defined Und[T] whose internal value is t.
func Defined[T any](t T) Und[T] {
	return Und[T]{
		opt: opt.Some(opt.Some(t)),
	}
}

// Null returns a null Und[T].
func Null[T any]() Und[T] {
	return Und[T]{
		opt: opt.Some(opt.None[T]()),
	}
}

// Undefined returns an undefined Und[T].
func Undefined[T any]() Und[T] {
	return Und[T]{}
}

// FromPointer converts *T into Und[T].
// If v is nil, it returns an undefined Und.
// Otherwise, it returns Defined[T] whose value is the dereferenced v.
//
// If you need to keep t as pointer, use [WrapPointer] instead.
func FromPointer[T any](v *T) Und[T] {
	if v == nil {
		return Undefined[T]()
	}
	return Defined(*v)
}

// WrapPointer converts *T into Und[*T].
// The und value is defined if t is non nil, undefined otherwise.
//
// If you want t to be dereferenced, use [FromPointer] instead.
func WrapPointer[T any](t *T) Und[*T] {
	if t == nil {
		return Undefined[*T]()
	}
	return Defined(t)
}

// FromOption converts o into an Und[T].
// o is retained by the returned value.
func WrapOption[T any](o opt.Option[opt.Option[T]]) Und[T] {
	return Und[T]{opt: o}
}

// FromSqlNull converts a valid sql.Null[T] to a defined Und[T]
// and invalid one into a null Und[T].
func FromSqlNull[T any](v sql.Null[T]) Und[T] {
	if !v.Valid {
		return Null[T]()
	}
	return Defined(v.V)
}

// IsZero is an alias for IsUndefined.
func (u Und[T]) IsZero() bool {
	return u.IsUndefined()
}

// IsDefined returns true if u is a defined value, otherwise false.
func (u Und[T]) IsDefined() bool {
	return u.opt.IsSome() && u.opt.Value().IsSome()
}

// IsNull returns true if u is a null value, otherwise false.
func (u Und[T]) IsNull() bool {
	return u.opt.IsSome() && u.opt.Value().IsNone()
}

// IsUndefined returns true if u is an undefined value, otherwise false.
func (u Und[T]) IsUndefined() bool {
	return u.opt.IsNone()
}

// EqualFunc reports whether two Und values are equal.
// EqualFunc checks state of both. If both state does not match, it returns false.
// If both are *defined* state, then it checks equality of their value by cmp,
// then returns true if they are equal.
//
// If T is just a comparable type, use [Equal].
// If T is an implementor of interface { Equal(t T) bool }, e.g time.Time, use [EqualEqualer].
func (u Und[T]) EqualFunc(t Und[T], cmp func(i, j T) bool) bool {
	return u.opt.EqualFunc(
		t.opt,
		func(i, j opt.Option[T]) bool {
			return i.EqualFunc(j, cmp)
		},
	)
}

// CloneFunc clones u using the cloneT functions.
func (u Und[T]) CloneFunc(cloneT func(T) T) Und[T] {
	return u.InnerMap(func(o opt.Option[opt.Option[T]]) opt.Option[opt.Option[T]] {
		return o.CloneFunc(func(o opt.Option[T]) opt.Option[T] {
			return o.CloneFunc(cloneT)
		})
	})
}

// Value returns an internal value.
func (u Und[T]) Value() T {
	if u.IsDefined() {
		return u.opt.Value().Value()
	}
	var zero T
	return zero
}

func (u Und[T]) Pointer() *opt.Option[T] {
	if u.IsUndefined() {
		return nil
	}
	return u.opt.Pointer()
}

// Unwrap returns u's internal value.
func (u Und[T]) Inner() opt.Option[opt.Option[T]] {
	return u.opt
}

// SqlNull converts u into sql.Null[T].
func (u Und[T]) SqlNull() sql.Null[T] {
	return u.opt.Value().SqlNull()
}

