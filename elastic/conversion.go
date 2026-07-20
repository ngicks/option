package elastic

import (
	"github.com/ngicks/option/opt"
	"github.com/ngicks/option/und"
)

// Strip returns the internal und value.
func (e Elastic[T]) Strip() und.Und[opt.Options[T]] {
	return e.v
}

// Wrap wraps given u with Elastic wrapper.
func Wrap[T any](u und.Und[opt.Options[T]]) Elastic[T] {
	return Elastic[T]{v: u}
}

// FromPointer converts *T into Elastic[T].
// If t is nil, it returns an undefined Elastic.
// Otherwise, it returns a defined Elastic whose single option is some of the dereferenced t.
func FromPointer[T any](t *T) Elastic[T] {
	if t == nil {
		return Undefined[T]()
	}
	return FromValue(*t)
}

// FromPointers converts each of ps into an option and returns a defined Elastic[T] holding them.
// Nil elements convert to none options, non-nil ones to some options of the dereferenced value.
func FromPointers[T any](ps ...*T) Elastic[T] {
	options := make(opt.Options[T], len(ps))
	for i, p := range ps {
		options[i] = opt.FromPointer(p)
	}
	return Elastic[T]{v: und.Defined(options)}
}

// FromUnd converts und.Und[T] into Elastic[T], preserving its state.
// A defined u becomes a defined Elastic whose single option is some of u's value.
func FromUnd[T any](u und.Und[T]) Elastic[T] {
	switch {
	case u.IsDefined():
		return FromValue(u.Value())
	case u.IsNull():
		return Null[T]()
	default:
		return Undefined[T]()
	}
}

// Und converts e into und.Und[T], preserving its state.
// A defined e becomes an und.Und of its first option:
// defined if the option is some, null if it is none or e holds no options.
func (e Elastic[T]) Und() und.Und[T] {
	switch {
	case e.IsUndefined():
		return und.Undefined[T]()
	case e.IsNull():
		return und.Null[T]()
	default:
		if options := e.v.Value(); len(options) > 0 {
			return und.FromOption(options[0], true)
		}
		return und.Null[T]()
	}
}
