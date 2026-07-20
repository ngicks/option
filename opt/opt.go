// Package opt defines 2-state value of T | null.
package opt

// Option represents an optional value.
type Option[T any] struct {
	some bool
	v    T
}

func Some[T any](v T) Option[T] {
	return Option[T]{
		some: true,
		v:    v,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsZero() bool {
	return o.IsNone()
}

// Value returns its internal as T.
// T would be zero value if o is None.
func (o Option[T]) Value() T {
	return o.v
}

// Pointer transforms o to *T, the plain conventional Go representation of an optional value.
// The pointer is reference to internal value.
// Use lock mechanism if multiple goroutines need to acccess it.
func (o *Option[T]) Pointer() *T {
	if o.IsNone() {
		return nil
	}
	return &o.v
}

// CloneFunc clones o using the cloneT function.
func (o Option[T]) CloneFunc(cloneT func(T) T) Option[T] {
	return o.Map(func(t T) T {
		return cloneT(t)
	})
}

// EqualFunc tests o and other if both are Some or None.
// If their state does not match, it returns false immediately.
// If both have value, it tests equality of their values by cmp.
//
// If T is just a comparable type, use [Equal].
// If T is an implementor of interface { Equal(t T) bool }, e.g time.Time, use [EqualEqualer].
func (o Option[T]) EqualFunc(other Option[T], cmp func(i, j T) bool) bool {
	if !o.some || !other.some {
		return o.some == other.some
	}

	return cmp(o.v, other.v)
}
