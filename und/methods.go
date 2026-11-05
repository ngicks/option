package und

import "github.com/ngicks/option/opt"

// Coalesce squashes 3-state u into 2-state u,
// (defined | null) if toNull is true, (defined | undefined) otherwise.
func (u Und[T]) Coalesce(toNull bool) Und[T] {
	if u.IsDefined() {
		return u
	}
	if toNull {
		return Null[T]()
	}
	return Undefined[T]()
}

// MapInner returns a new Und[T] whose internal value is u's mapped by f.
func (u Und[T]) MapInner(f func(opt.Option[opt.Option[T]]) opt.Option[opt.Option[T]]) Und[T] {
	return Und[T]{opt: f(u.opt)}
}
