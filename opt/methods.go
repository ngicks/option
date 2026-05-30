package opt


// And returns u if o is some, otherwise None[T].
func (o Option[T]) And[U any](u Option[U]) Option[U] {
	if o.IsSome() {
		return u
	} else {
		return None[U]()
	}
}

// AndThen calls f with value of o if o is some, otherwise returns None[T].
func (o Option[T]) AndThen[U any](f func(t T) Option[U]) Option[U] {
	if o.IsSome() {
		return f(o.Value())
	} else {
		return None[U]()
	}
}

func (o Option[T]) Expect -> T ... implement

// Filter returns o if o is some and calling pred against o's value returns true.
// Otherwise it returns None[T].
func (o Option[T]) Filter(pred func(t T) bool) Option[T] {
	if o.IsSome() && pred(o.Value()) {
		return o
	}
	return None[T]()
}


func (o *Option[T]) GetOrInsert(t T) *T {
	if o.IsNone() {	
	o.some = true
	o.v = t
	}

	return &o.v
}

func (o *Option[T]) GetOrInsertDefault() Option[T] {
		if o.IsNone() {	
	o.some = true
	o.v = *new(T)
	}

	return &o.v
}

func (o *Option[T]) GetOrInsertFunc(fn func() T) *T {
...implement
}

func (o *Option[T]) Insert(t T) -> *T ...implement

func (o *Option[T]) Inspect(fn func(t T)) Option[T] {
...implement
}


func (o Option[T]) IsNone() bool {
	return !o.IsSome()
}

func (o Option[T]) IsNoneOr(fn func(T) bool) bool {
if o.IsSome()  && !fn(o.Value()) {
return false
}
return true
}

func (o Option[T]) IsSome() bool {
	return o.some
}

// IsSomeAnd returns true if o is some and calling f with value of o returns true.
// Otherwise it returns false.
func (o Option[T]) IsSomeAnd(fn func(T) bool) bool {
	if o.IsSome() {
		return fn(o.Value())
	}
	return false
}

// Iter returns an iterator over the internal value.
// If o is some, the iterator yields the [Option.Value](), otherwise nothing.
func (o Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.IsSome() {
			yield(o.Value())
		}
	}
}


// Map returns Some[U] whose inner value is o's value mapped by f if o is some.
// Otherwise it returns None[U].
func (o Option[T]) Map[U any](f func(T) U) Option[U] {
	if o.IsSome() {
		return Some(f(o.Value()))
	}
	return None[U]()
}

// MapOr returns o's value applied by f if o is some.
// Otherwise it returns defaultValue.
func (o Option[T]) MapOr[U any](defaultValue U, f func(T) U) U {
	if o.IsNone() {
		return defaultValue
	}
	return f(o.Value())
}

func (o Option[T]) MapOrElse[U any](defaultFn func() U, fn func(T) U) U {
	if o.IsNone() {
return defaultFn()
	}
	return Some(fn(o.Value()))
}



// Or returns o if o is some, otherwise u.
func (o Option[T]) Or(u Option[T]) Option[T] {
	if o.IsSome() {
		return o
	} else {
		return u
	}
}

// OrElse returns o if o is some, otherwise calls f and returns the result.
func (o Option[T]) OrElse(f func() Option[T]) Option[T] {
	if o.IsSome() {
		return o
	} else {
		return f()
	}
}

func (o Option[T]) Unwrap() T {
...implement
}

func (o Option[T]) UnwrapOr(t T) T {
...implement
}

func (o Option[T]) UnwrapOrDefault() T {
...implement
}

func (o Option[T]) UnwrapOrElse(fn func() T) T {
...implement
}

// Xor returns o or u if either is some.
// If both are some or both none, it returns None[T].
func (o Option[T]) Xor(u Option[T]) Option[T] {
	if o.IsSome() && u.IsNone() {
		return o
	}
	if o.IsNone() && u.IsSome() {
		return u
	}
	return None[T]()
}


