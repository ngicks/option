package elastic

import (
	"slices"

	"github.com/ngicks/option/opt"
	"github.com/ngicks/option/und"
)

// Coalesce squashes 3-state e into 2-state e,
// (defined | null) if toNull is true, (defined | undefined) otherwise.
func (e Elastic[T]) Coalesce(toNull bool) Elastic[T] {
	if e.IsDefined() {
		return e
	}
	if toNull {
		return Null[T]()
	}
	return Undefined[T]()
}

// MapInner returns a new Elastic[T] whose internal value is e's mapped by f.
//
// The internal slice of e is capped to its length before passed to f,
// so appending to it in f does not overwrite the backing array shared with e.
func (e Elastic[T]) MapInner(f func(und.Und[opt.Options[T]]) und.Und[opt.Options[T]]) Elastic[T] {
	v := e.v
	if v.IsDefined() {
		v = und.Defined(slices.Clip(v.Value()))
	}
	return Elastic[T]{v: f(v)}
}
