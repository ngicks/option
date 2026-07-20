package either

import "iter"

// ExpectLeft returns e's left value if e is left.
// Otherwise it panics with msg.
func (e Either[L, R]) ExpectLeft(msg string) L {
	if e.IsRight() {
		panic(msg)
	}
	return e.l
}

// ExpectRight returns e's right value if e is right.
// Otherwise it panics with msg.
func (e Either[L, R]) ExpectRight(msg string) R {
	if e.IsLeft() {
		panic(msg)
	}
	return e.r
}

// Flip converts Either[L, R] into Either[R, L], a left into a right and vice versa.
func (e Either[L, R]) Flip() Either[R, L] {
	if e.IsRight() {
		return Left[R, L](e.r)
	}
	return Right[R](e.l)
}

// Fold applies onLeft to a left value or onRight to a right value,
// squashing both cases into a single U.
func (e Either[L, R]) Fold[U any](onLeft func(l L) U, onRight func(r R) U) U {
	if e.IsRight() {
		return onRight(e.r)
	}
	return onLeft(e.l)
}

// InspectLeft calls fn with e's left value if e is left, then returns e unchanged.
func (e Either[L, R]) InspectLeft(fn func(l L)) Either[L, R] {
	if e.IsLeft() {
		fn(e.l)
	}
	return e
}

// InspectRight calls fn with e's right value if e is right, then returns e unchanged.
func (e Either[L, R]) InspectRight(fn func(r R)) Either[L, R] {
	if e.IsRight() {
		fn(e.r)
	}
	return e
}

// IsLeftAnd returns true if e is left and calling fn with e's left value returns true.
// Otherwise it returns false.
func (e Either[L, R]) IsLeftAnd(fn func(l L) bool) bool {
	if e.IsLeft() {
		return fn(e.l)
	}
	return false
}

// IsRightAnd returns true if e is right and calling fn with e's right value returns true.
// Otherwise it returns false.
func (e Either[L, R]) IsRightAnd(fn func(r R) bool) bool {
	if e.IsRight() {
		return fn(e.r)
	}
	return false
}

// IterLeft returns an iterator over the left value.
// If e is left, the iterator yields the left value, otherwise nothing.
func (e Either[L, R]) IterLeft() iter.Seq[L] {
	return func(yield func(L) bool) {
		if e.IsLeft() {
			yield(e.l)
		}
	}
}

// IterRight returns an iterator over the right value.
// If e is right, the iterator yields the right value, otherwise nothing.
func (e Either[L, R]) IterRight() iter.Seq[R] {
	return func(yield func(R) bool) {
		if e.IsRight() {
			yield(e.r)
		}
	}
}

// LeftAnd returns u if e is left, otherwise Right[U] holding e's right value.
func (e Either[L, R]) LeftAnd[U any](u Either[U, R]) Either[U, R] {
	if e.IsRight() {
		return Right[U](e.r)
	}
	return u
}

// LeftAndThen calls fn with e's left value if e is left,
// otherwise returns Right[U] holding e's right value.
func (e Either[L, R]) LeftAndThen[U any](fn func(l L) Either[U, R]) Either[U, R] {
	if e.IsRight() {
		return Right[U](e.r)
	}
	return fn(e.l)
}

// LeftOr returns e's left value if e is left, otherwise l.
func (e Either[L, R]) LeftOr(l L) L {
	if e.IsLeft() {
		return e.l
	}
	return l
}

// LeftOrElse returns e's left value if e is left,
// otherwise the result of calling fn with e's right value.
func (e Either[L, R]) LeftOrElse(fn func(r R) L) L {
	if e.IsLeft() {
		return e.l
	}
	return fn(e.r)
}

// Map maps a left value by onLeft or a right value by onRight.
func (e Either[L, R]) Map[L2, R2 any](onLeft func(l L) L2, onRight func(r R) R2) Either[L2, R2] {
	if e.IsRight() {
		return Right[L2](onRight(e.r))
	}
	return Left[L2, R2](onLeft(e.l))
}

// MapLeft maps a left value by fn, leaving a right value untouched.
func (e Either[L, R]) MapLeft[U any](fn func(l L) U) Either[U, R] {
	if e.IsRight() {
		return Right[U](e.r)
	}
	return Left[U, R](fn(e.l))
}

// MapLeftOr returns e's left value applied by fn if e is left.
// Otherwise it returns defaultValue.
func (e Either[L, R]) MapLeftOr[U any](defaultValue U, fn func(l L) U) U {
	if e.IsRight() {
		return defaultValue
	}
	return fn(e.l)
}

// MapRight maps a right value by fn, leaving a left value untouched.
func (e Either[L, R]) MapRight[U any](fn func(r R) U) Either[L, U] {
	if e.IsRight() {
		return Right[L](fn(e.r))
	}
	return Left[L, U](e.l)
}

// MapRightOr returns e's right value applied by fn if e is right.
// Otherwise it returns defaultValue.
func (e Either[L, R]) MapRightOr[U any](defaultValue U, fn func(r R) U) U {
	if e.IsLeft() {
		return defaultValue
	}
	return fn(e.r)
}

// RightAnd returns u if e is right, otherwise Left[L] holding e's left value.
func (e Either[L, R]) RightAnd[U any](u Either[L, U]) Either[L, U] {
	if e.IsLeft() {
		return Left[L, U](e.l)
	}
	return u
}

// RightAndThen calls fn with e's right value if e is right,
// otherwise returns Left[L] holding e's left value.
func (e Either[L, R]) RightAndThen[U any](fn func(r R) Either[L, U]) Either[L, U] {
	if e.IsLeft() {
		return Left[L, U](e.l)
	}
	return fn(e.r)
}

// RightOr returns e's right value if e is right, otherwise r.
func (e Either[L, R]) RightOr(r R) R {
	if e.IsRight() {
		return e.r
	}
	return r
}

// RightOrElse returns e's right value if e is right,
// otherwise the result of calling fn with e's left value.
func (e Either[L, R]) RightOrElse(fn func(l L) R) R {
	if e.IsRight() {
		return e.r
	}
	return fn(e.l)
}

// UnwrapLeft returns e's left value if e is left.
// Otherwise it panics.
func (e Either[L, R]) UnwrapLeft() L {
	if e.IsRight() {
		panic("either: UnwrapLeft called on a right Either")
	}
	return e.l
}

// UnwrapRight returns e's right value if e is right.
// Otherwise it panics.
func (e Either[L, R]) UnwrapRight() R {
	if e.IsLeft() {
		panic("either: UnwrapRight called on a left Either")
	}
	return e.r
}
