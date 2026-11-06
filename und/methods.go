package und

import "github.com/ngicks/option/opt"

// CollapseToNull returns u as is if u is defined, Null otherwise.
func (u Und[T]) CollapseToNull() Und[T] {
	if u.IsDefined() {
		return u
	}
	return Null[T]()
}

// CollapseToUndefined returns u as is if u is defined, Undefined otherwise.
func (u Und[T]) CollapseToUndefined() Und[T] {
	if u.IsDefined() {
		return u
	}
	return Undefined[T]()
}

// MapInner returns a new Und[T] whose internal value is u's mapped by f.
func (u Und[T]) MapInner(f func(opt.Option[opt.Option[T]]) opt.Option[opt.Option[T]]) Und[T] {
	return Und[T]{opt: f(u.opt)}
}
