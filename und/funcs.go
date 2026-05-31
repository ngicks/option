package und

// Equal tests equality of l and r then returns true if they are equal, false otherwise.
// For those types that are comparable but need special tests, e.g. time.Time, you should use [Und.EqualFunc] instead.
//
// If T is an implementor of interface { Equal(t T) bool }, e.g time.Time, use [EqualEqualer].
func Equal[T comparable](l, r Und[T]) bool {
	return l.EqualFunc(r, func(i, j T) bool { return i == j })
}

// EqualEqualer tests equality of l and r by calling Equal method implemented on l.
func EqualEqualer[T interface{ Equal(t T) bool }](l, r Und[T]) bool {
	return l.EqualFunc(r, func(i, j T) bool {
		return i.Equal(j)
	})
}

