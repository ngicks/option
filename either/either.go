// Package either represents a value ether L or R.
package either

type Either[L, R any] struct {
	l     L
	r     R
	right bool // true when value is R
}

func Left[L, R any](l L) Either[L, R] {
	return Either[L, R]{l: l}
}

func Right[L, R any](r R) Either[L, R] {
	return Either[L, R]{r: r, right: true}
}

// IsLeft returns true if e holds an L value.
// Note that the zero value of Either[L, R] is left holding the zero value of L.
func (e Either[L, R]) IsLeft() bool {
	return !e.right
}

// IsRight returns true if e holds an R value.
func (e Either[L, R]) IsRight() bool {
	return e.right
}

// Unpack returns both internal values of e.
// The value for the side e does not hold is the zero value of its type.
func (e Either[L, R]) Unpack() (l L, r R) {
	return e.l, e.r
}

// EqualFunc tests e and other if both are left or right.
// If their state does not match, it returns false immediately.
// Otherwise it tests equality of the held values by cmpL or cmpR respectively.
//
// If L and R are just comparable types, use [Equal].
// If L and R are implementors of interface { Equal(t T) bool }, e.g time.Time, use [EqualEqualer].
func (e Either[L, R]) EqualFunc(
	other Either[L, R],
	cmpL func(i, j L) bool,
	cmpR func(i, j R) bool,
) bool {
	if e.right != other.right {
		return false
	}
	if e.right {
		return cmpR(e.r, other.r)
	}
	return cmpL(e.l, other.l)
}

// CloneFunc clones e using the cloneL or cloneR function for the held value.
func (e Either[L, R]) CloneFunc(cloneL func(L) L, cloneR func(R) R) Either[L, R] {
	if e.right {
		return Right[L](cloneR(e.r))
	}
	return Left[L, R](cloneL(e.l))
}
