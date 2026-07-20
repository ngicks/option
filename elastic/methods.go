package elastic

// Coalesce squashes 3-state e into 2-state e,
// (defined | null) if toNull is true, (defined | undefined) otherwise.
func (e Elastic[T]) Coalesce(toNull bool) Elastic[T] {
	if e.IsDefined() {
		return e
	}
	if toNull {
		return Null[T]()
	}
	return Undefined[T]()
}
