package und

// Coalesce squashes 3-state u into 2-state u,
// (defined | null) if toNull is true, (defined | undefined) otherwise.
func (u Und[T]) Coalesce(toNull bool) Und[T] {
	if u.IsDefined() {
		return u
	}
	if toNull {
		return Null[T]()
	}
	return Undefined[T]()
}
