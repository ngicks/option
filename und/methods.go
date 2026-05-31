package und


// Deprecated: Renamed to [Und.InnerMap]. Und.Map has same name but behavior is inconsistent to [Map].
func (u Und[T]) Map(f func(opt.Option[opt.Option[T]]) opt.Option[opt.Option[T]]) Und[T] {
	return u.InnerMap(f)
}

// InnerMap returns a new Und[T] whose internal value is u's mapped by f.
// Unlike [Map], f is invoked even when u is not an undined value.
func (u Und[T]) InnerMap(f func(opt.Option[opt.Option[T]]) opt.Option[opt.Option[T]]) Und[T] {
	return Und[T]{opt: f(u.opt)}
}

