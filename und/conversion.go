package und

import (
	"database/sql"

	"github.com/ngicks/option/opt"
)

// Strip returns the internal option value..
func (u Und[T]) Strip() opt.Option[opt.Option[T]] {
	return u.opt
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

// WrapSlice converts s into Und[S].
// The und value is defined if s is non nil, undefined otherwise.
// An empty non-nil s is defined.
func WrapSlice[S ~[]E, E any](s S) Und[S] {
	if s == nil {
		return Undefined[S]()
	}
	return Defined(s)
}

// WrapMap converts m into Und[M].
// The und value is defined if m is non nil, undefined otherwise.
// An empty non-nil m is defined.
func WrapMap[M ~map[K]V, K comparable, V any](m M) Und[M] {
	if m == nil {
		return Undefined[M]()
	}
	return Defined(m)
}

// Wrap wraps given o with Und wrapper.
func Wrap[T any](o opt.Option[opt.Option[T]]) Und[T] {
	return Und[T]{opt: o}
}

func FromSqlNull[T any](nullable sql.Null[T]) Und[T] {
	if nullable.Valid {
		return Defined(nullable.V)
	}
	return Null[T]()
}

// SqlNull converts u into sql.Null[T].
func (u Und[T]) SqlNull() sql.Null[T] {
	return u.opt.Value().SqlNull()
}

func FromOption[T any](o opt.Option[T], noneIsNull bool) Und[T] {
	if o.IsSome() {
		return Defined(o.Value())
	}
	if noneIsNull {
		return Null[T]()
	}
	return Undefined[T]()
}

func (u Und[T]) Flatten() opt.Option[T] {
	return opt.Flatten(u.opt)
}
