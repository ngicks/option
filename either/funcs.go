package either

// CloneCloner clones e by calling the Clone method on its held value.
//
// For custom clone functions, use [Either.CloneFunc].
// An Either[L, R] whose L and R only need a shallow copy can be copied by plain assignment.
func CloneCloner[L interface{ Clone() L }, R interface{ Clone() R }](e Either[L, R]) Either[L, R] {
	return e.CloneFunc(
		func(l L) L { return l.Clone() },
		func(r R) R { return r.Clone() },
	)
}

// FlattenLeft converts a nested Either[Either[L, R], R] into an Either[L, R].
// The result is the inner either if e is left, Right[L] holding e's right value otherwise.
func FlattenLeft[L, R any](e Either[Either[L, R], R]) Either[L, R] {
	if e.IsRight() {
		return Right[L](e.r)
	}
	return e.l
}

// FlattenRight converts a nested Either[L, Either[L, R]] into an Either[L, R].
// The result is the inner either if e is right, Left[L, R] holding e's left value otherwise.
func FlattenRight[L, R any](e Either[L, Either[L, R]]) Either[L, R] {
	if e.IsLeft() {
		return Left[L, R](e.l)
	}
	return e.r
}

// FlattenEither converts an Either[Either[L, R], Either[L, R]] into an Either[L, R],
// returning the inner either of whichever side e holds.
func FlattenEither[L, R any](e Either[Either[L, R], Either[L, R]]) Either[L, R] {
	if e.IsRight() {
		return e.r
	}
	return e.l
}

// Equal tests equality of l and r then returns true if they are equal, false otherwise.
func Equal[L, R comparable](l, r Either[L, R]) bool {
	return l.EqualFunc(
		r,
		func(i, j L) bool { return i == j },
		func(i, j R) bool { return i == j },
	)
}

// EqualEqualer tests equality of l and r by calling the Equal method implemented on their held
// values.
func EqualEqualer[
	L interface{ Equal(l L) bool },
	R interface{ Equal(r R) bool },
](l, r Either[L, R]) bool {
	return l.EqualFunc(
		r,
		func(i, j L) bool { return i.Equal(j) },
		func(i, j R) bool { return i.Equal(j) },
	)
}
