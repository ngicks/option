package opt


func Clone ... implement

func CloneCloner ... implement

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


