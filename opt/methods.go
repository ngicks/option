package opt

import "iter"

// And returns u if o is some, otherwise None[T].
func (o Option[T]) And[U any](u Option[U]) Option[U] {
	if o.IsSome() {
		return u
	} else {
		return None[U]()
	}
}

// AndThen calls f with value of o if o is some, otherwise returns None[T].
func (o Option[T]) AndThen[U any](f func(t T) Option[U]) Option[U] {
	if o.IsSome() {
		return f(o.Value())
	} else {
		return None[U]()
	}
}

// Expect returns o's value if o is some.
// Otherwise it panics with msg.
func (o Option[T]) Expect(msg string) T {
	if o.IsNone() {
		panic(msg)
	}
	return o.v
}

// Filter returns o if o is some and calling pred against o's value returns true.
// Otherwise it returns None[T].
func (o Option[T]) Filter(pred func(t T) bool) Option[T] {
	if o.IsSome() && pred(o.Value()) {
		return o
	}
	return None[T]()
}

// GetOrInsert inserts t into o if o is none, then returns a pointer to o's value.
// If o is already some, the existing value is kept and t is discarded.
//
// o must be addressable; the returned pointer aliases o's internal value.
func (o *Option[T]) GetOrInsert(t T) *T {
	if o.IsNone() {
		o.some = true
		o.v = t
	}
	return &o.v
}

// GetOrInsertDefault inserts the zero value of T into o if o is none,
// then returns a pointer to o's value.
// If o is already some, the existing value is kept.
//
// o must be addressable; the returned pointer aliases o's internal value.
func (o *Option[T]) GetOrInsertDefault() *T {
	if o.IsNone() {
		o.some = true
		var zero T
		o.v = zero
	}
	return &o.v
}

// GetOrInsertFunc inserts the result of fn into o if o is none,
// then returns a pointer to o's value.
// If o is already some, the existing value is kept and fn is not called.
//
// o must be addressable; the returned pointer aliases o's internal value.
func (o *Option[T]) GetOrInsertFunc(fn func() T) *T {
	if o.IsNone() {
		o.some = true
		o.v = fn()
	}
	return &o.v
}

// Insert sets o's value to t, overwriting any existing value, and returns a pointer to it.
//
// See [Option.GetOrInsert] which only inserts when o is none.
//
// o must be addressable; the returned pointer aliases o's internal value.
func (o *Option[T]) Insert(t T) *T {
	o.some = true
	o.v = t
	return &o.v
}

// Inspect calls fn with o's value if o is some, then returns o unchanged.
func (o Option[T]) Inspect(fn func(t T)) Option[T] {
	if o.IsSome() {
		fn(o.Value())
	}
	return o
}

func (o Option[T]) IsNone() bool {
	return !o.IsSome()
}

// IsNoneOr returns true if o is none, or if o is some and calling fn with o's value returns true.
func (o Option[T]) IsNoneOr(fn func(T) bool) bool {
	if o.IsSome() && !fn(o.Value()) {
		return false
	}
	return true
}

func (o Option[T]) IsSome() bool {
	return o.some
}

// IsSomeAnd returns true if o is some and calling f with value of o returns true.
// Otherwise it returns false.
func (o Option[T]) IsSomeAnd(fn func(T) bool) bool {
	if o.IsSome() {
		return fn(o.Value())
	}
	return false
}

// Iter returns an iterator over the internal value.
// If o is some, the iterator yields the [Option.Value](), otherwise nothing.
func (o Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.IsSome() {
			yield(o.Value())
		}
	}
}

// Map returns Some[U] whose inner value is o's value mapped by f if o is some.
// Otherwise it returns None[U].
func (o Option[T]) Map[U any](f func(T) U) Option[U] {
	if o.IsSome() {
		return Some(f(o.Value()))
	}
	return None[U]()
}

// MapOr returns o's value applied by f if o is some.
// Otherwise it returns defaultValue.
func (o Option[T]) MapOr[U any](defaultValue U, f func(T) U) U {
	if o.IsNone() {
		return defaultValue
	}
	return f(o.Value())
}

// MapOrElse returns o's value applied by fn if o is some.
// Otherwise it returns the result of defaultFn.
func (o Option[T]) MapOrElse[U any](defaultFn func() U, fn func(T) U) U {
	if o.IsNone() {
		return defaultFn()
	}
	return fn(o.Value())
}

// Or returns o if o is some, otherwise u.
func (o Option[T]) Or(u Option[T]) Option[T] {
	if o.IsSome() {
		return o
	} else {
		return u
	}
}

// OrElse returns o if o is some, otherwise calls f and returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	} else {
		return f()
	}
}

// Unwrap returns o's value if o is some.
// Otherwise it panics.
func (o Option[T]) Unwrap() T {
	if o.IsNone() {
		panic("opt: Unwrap called on a none Option")
	}
	return o.v
}

// UnwrapOr returns o's value if o is some, otherwise t.
func (o Option[T]) UnwrapOr(t T) T {
	if o.IsSome() {
		return o.v
	}
	return t
}

// UnwrapOrDefault returns o's value if o is some, otherwise the zero value of T.
func (o Option[T]) UnwrapOrDefault() T {
	if o.IsSome() {
		return o.v
	}
	return *new(T)
}

// UnwrapOrElse returns o's value if o is some, otherwise the result of fn.
func (o Option[T]) UnwrapOrElse(fn func() T) T {
	if o.IsSome() {
		return o.v
	}
	return fn()
}

// Xor returns o or u if either is some.
// If both are some or both none, it returns None[T].
func (o Option[T]) Xor(u Option[T]) Option[T] {
	if o.IsSome() && u.IsNone() {
		return o
	}
	if o.IsNone() && u.IsSome() {
		return u
	}
	return None[T]()
}
