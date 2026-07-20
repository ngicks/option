// Package elastic defines multi-state value of (T | null | undefined | [](null | T)). Useful for consuming Elasticsearch fields or config file.
package elastic

import (
	"github.com/ngicks/option/opt"
	"github.com/ngicks/option/und"
)

type Elastic[T any] struct {
	v und.Und[opt.Options[T]]
}

// FromValue returns a defined Elastic[T] whose internal value is a single some option of t.
func FromValue[T any](t T) Elastic[T] {
	return FromOptions(opt.Some(t))
}

// FromValues returns a defined Elastic[T] whose internal value is some options of each of ts.
func FromValues[T any](ts ...T) Elastic[T] {
	options := make(opt.Options[T], len(ts))
	for i, t := range ts {
		options[i] = opt.Some(t)
	}
	return Elastic[T]{v: und.Defined(options)}
}

// FromOptions returns a defined Elastic[T] whose internal value is options.
// options is retained by the returned value.
func FromOptions[T any](options ...opt.Option[T]) Elastic[T] {
	return Elastic[T]{v: und.Defined(opt.Options[T](options))}
}

// Null returns a null Elastic[T].
func Null[T any]() Elastic[T] {
	return Elastic[T]{v: und.Null[opt.Options[T]]()}
}

// Undefined returns an undefined Elastic[T].
func Undefined[T any]() Elastic[T] {
	return Elastic[T]{}
}

// IsZero is an alias for IsUndefined.
func (e Elastic[T]) IsZero() bool {
	return e.IsUndefined()
}

// IsDefined returns true if e holds at least zero options, otherwise false.
func (e Elastic[T]) IsDefined() bool {
	return e.v.IsDefined()
}

// IsNull returns true if e is a null value, otherwise false.
func (e Elastic[T]) IsNull() bool {
	return e.v.IsNull()
}

// IsUndefined returns true if e is an undefined value, otherwise false.
func (e Elastic[T]) IsUndefined() bool {
	return e.v.IsUndefined()
}

// Value returns the first value of e.
// T would be zero value if e is not defined, holds no options, or its first option is none.
func (e Elastic[T]) Value() T {
	if options := e.v.Value(); len(options) > 0 {
		return options[0].Value()
	}
	var zero T
	return zero
}

// Values returns all values e holds as []T.
// None options are represented as the zero value of T.
// It returns nil if e is not defined.
//
// For the option-preserving form, use [Elastic.Options].
func (e Elastic[T]) Values() []T {
	options := e.v.Value()
	if options == nil {
		return nil
	}
	vs := make([]T, len(options))
	for i, o := range options {
		vs[i] = o.Value()
	}
	return vs
}

// Options returns the internal options.
func (e Elastic[T]) Options() opt.Options[T] {
	return e.v.Value()
}

// Pointer returns a pointer to the first value of e,
// or nil if e is not defined, holds no options, or its first option is none.
func (e Elastic[T]) Pointer() *T {
	if options := e.v.Value(); len(options) > 0 {
		return options[0].Pointer()
	}
	return nil
}

// Pointers returns all values e holds as []*T.
// None options are represented as nil.
// It returns nil if e is not defined.
func (e Elastic[T]) Pointers() []*T {
	options := e.v.Value()
	if options == nil {
		return nil
	}
	ptrs := make([]*T, len(options))
	for i := range options {
		ptrs[i] = options[i].Pointer()
	}
	return ptrs
}

// EqualFunc reports whether two Elastic values are equal.
// EqualFunc checks state of both. If both state does not match, it returns false.
// If both are *defined* state, then it compares their options,
// testing values of two some options by cmp,
// then returns true if all of them are equal.
//
// If T is just a comparable type, use [Equal].
// If T is an implementor of interface { Equal(t T) bool }, e.g time.Time, use [EqualEqualer].
func (e Elastic[T]) EqualFunc(other Elastic[T], cmp func(i, j T) bool) bool {
	return e.v.EqualFunc(
		other.v,
		func(i, j opt.Options[T]) bool {
			return i.EqualFunc(j, cmp)
		},
	)
}

// CloneFunc clones e using the cloneT function for each value.
func (e Elastic[T]) CloneFunc(cloneT func(T) T) Elastic[T] {
	return Elastic[T]{
		v: e.v.CloneFunc(func(options opt.Options[T]) opt.Options[T] {
			return options.CloneFunc(cloneT)
		}),
	}
}
