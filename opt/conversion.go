package opt

import "database/sql"

// FromPointer converts *T into Option[T].
// If v is nil, it returns a none Option.
// Otherwise, it returns some Option whose value is the dereferenced v.
//
// If you need to keep t as pointer, use [WrapPointer] instead.
func FromPointer[T any](t *T) Option[T] {
	if t != nil {
		return Some(*t)
	}
	return None[T]()
}

// WrapPointer converts *T into Option[*T].
// The option is some if t is non nil, none otherwise.
//
// If you want t to be dereferenced, use [FromPointer] instead.
func WrapPointer[T any](t *T) Option[*T] {
	if t != nil {
		return Some(t)
	}
	return None[*T]()
}

func FromSqlNull[T any](v sql.Null[T]) Option[T] {
	if !v.Valid {
		return None[T]()
	}
	return Some(v.V)
}

func (o Option[T]) SqlNull() sql.Null[T] {
	if o.IsNone() {
		return sql.Null[T]{}
	}
	return sql.Null[T]{
		Valid: true,
		V:     o.Value(),
	}
}

// FromOk converts conventional (t T, ok bool) into an option.
// The options is some if ok is true, none otherwise.
//
// For getting values from maps or slices, instead you may want to use [IndexMap], [IndexSlice]
// respectively.
func FromOk[T any](t T, ok bool) Option[T] {
	if ok {
		return Some(t)
	}
	return None[T]()
}

// Assert type-asserts v into T.
// If v's internal value is T then returns Some of that value,
// None otherwise.
func Assert[T any](v any) Option[T] {
	if t, ok := v.(T); ok {
		return Some(t)
	}
	return None[T]()
}

// IndexMap gets a value associated with key.
// If key has a value, the Option is some wrapping the value.
// Otherwise it returns none Option.
func IndexMap[M ~map[K]V, K comparable, V any](m M, key K) Option[V] {
	v, ok := m[key]
	if !ok {
		return None[V]()
	}
	return Some(v)
}

// IndexSlice gets a value associated with idx.
// If idx is within interval [0, len(s)), then the Option is some wrapping a value associated to the
// idx.
// Otherwise it returns none Option.
func IndexSlice[S ~[]T, T any](s S, idx int) Option[T] {
	if idx < 0 || len(s) <= idx {
		return None[T]()
	}
	return Some(s[idx])
}
