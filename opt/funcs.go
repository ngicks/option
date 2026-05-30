package opt

// CloneCloner clones o by calling the Clone method on its inner value if o is some.
//
// For a custom clone function, use [Option.CloneFunc].
// An Option[T] whose T only needs a shallow copy can be copied by plain assignment.
func CloneCloner[T interface{ Clone() T }](o Option[T]) Option[T] {
	return o.CloneFunc(func(t T) T {
		return t.Clone()
	})
}

// Equal tests equality of l and r then returns true if they are equal, false otherwise
func Equal[T comparable](l, r Option[T]) bool {
	return l.EqualFunc(r, func(i, j T) bool { return i == j })
}

// EqualEqualer tests equality of l and r by calling Equal method implemented on l.
func EqualEqualer[T interface{ Equal(t T) bool }](l, r Option[T]) bool {
	return l.EqualFunc(r, func(i, j T) bool {
		return i.Equal(j)
	})
}

// Flatten converts Option[Option[T]] into Option[T].
func Flatten[T any](o Option[Option[T]]) Option[T] {
	if o.IsNone() {
		return None[T]()
	}
	v := o.Value()
	if v.IsNone() {
		return None[T]()
	}
	return v
}
