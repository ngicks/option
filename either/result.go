package either

import (
	"iter"

	"github.com/ngicks/option/opt"
)

type Result[T any] struct {
	e Either[T, error]
}

// Ok returns a Result[T] holding t.
func Ok[T any](t T) Result[T] {
	return Result[T]{e: Left[T, error](t)}
}

// Err returns a Result[T] holding err.
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("either: nil error is passed to Err")
	}
	return Result[T]{e: Right[T](err)}
}

// FromResult converts the conventional (t T, err error) pair into a Result[T].
// The result is err if err is non nil, ok wrapping t otherwise.
func FromResult[T any](t T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return Ok(t)
}

// Flatten converts a nested Result[Result[T]] into a Result[T].
// The result is the inner result if r is ok, Err[T] holding r's error otherwise.
func Flatten[T any](r Result[Result[T]]) Result[T] {
	if r.IsErr() {
		return Err[T](r.UnwrapErr())
	}
	return r.Unwrap()
}

// Transpose converts a Result[opt.Option[T]] into an opt.Option[Result[T]].
// An ok wrapping none maps to none.
// An ok wrapping some maps to some wrapping ok, and an err maps to some wrapping err.
func Transpose[T any](r Result[opt.Option[T]]) opt.Option[Result[T]] {
	if r.IsErr() {
		return opt.Some(Err[T](r.UnwrapErr()))
	}
	o := r.Unwrap()
	if o.IsNone() {
		return opt.None[Result[T]]()
	}
	return opt.Some(Ok(o.Value()))
}

// Either returns the internal either value.
func (r Result[T]) Either() Either[T, error] {
	return r.e
}

// And returns u if r is ok, otherwise Err[U] holding r's error.
func (r Result[T]) And[U any](u Result[U]) Result[U] {
	return Result[U]{e: r.e.LeftAnd(u.e)}
}

// AndThen calls fn with r's value if r is ok, otherwise returns Err[U] holding r's error.
func (r Result[T]) AndThen[U any](fn func(t T) Result[U]) Result[U] {
	return Result[U]{e: r.e.LeftAndThen(func(t T) Either[U, error] { return fn(t).e })}
}

// Err returns r's error if r is err, nil otherwise.
func (r Result[T]) Err() error {
	return r.e.RightOr(nil)
}

// Expect returns r's value if r is ok.
// Otherwise it panics with msg.
func (r Result[T]) Expect(msg string) T {
	return r.e.ExpectLeft(msg)
}

// ExpectErr returns r's error if r is err.
// Otherwise it panics with msg.
func (r Result[T]) ExpectErr(msg string) error {
	return r.e.ExpectRight(msg)
}

// Inspect calls fn with r's value if r is ok, then returns r unchanged.
func (r Result[T]) Inspect(fn func(t T)) Result[T] {
	return Result[T]{e: r.e.InspectLeft(fn)}
}

// InspectErr calls fn with r's error if r is err, then returns r unchanged.
func (r Result[T]) InspectErr(fn func(err error)) Result[T] {
	return Result[T]{e: r.e.InspectRight(fn)}
}

// IsErr returns true if r holds an error.
func (r Result[T]) IsErr() bool {
	return r.e.IsRight()
}

// IsErrAnd returns true if r is err and calling fn with r's error returns true.
// Otherwise it returns false.
func (r Result[T]) IsErrAnd(fn func(err error) bool) bool {
	return r.e.IsRightAnd(fn)
}

// IsOk returns true if r holds a T value.
func (r Result[T]) IsOk() bool {
	return r.e.IsLeft()
}

// IsOkAnd returns true if r is ok and calling fn with r's value returns true.
// Otherwise it returns false.
func (r Result[T]) IsOkAnd(fn func(t T) bool) bool {
	return r.e.IsLeftAnd(fn)
}

// Iter returns an iterator over the internal value.
// If r is ok, the iterator yields the value, otherwise nothing.
func (r Result[T]) Iter() iter.Seq[T] {
	return r.e.IterLeft()
}

// Map returns Ok[U] whose inner value is r's value mapped by fn if r is ok.
// Otherwise it returns Err[U] holding r's error.
func (r Result[T]) Map[U any](fn func(t T) U) Result[U] {
	return Result[U]{e: r.e.MapLeft(fn)}
}

// MapErr returns Err[T] whose error is r's error mapped by fn if r is err.
// Otherwise it returns r unchanged.
func (r Result[T]) MapErr(fn func(err error) error) Result[T] {
	return Result[T]{e: r.e.MapRight(fn)}
}

// MapOr returns r's value applied by fn if r is ok.
// Otherwise it returns defaultValue.
func (r Result[T]) MapOr[U any](defaultValue U, fn func(t T) U) U {
	return r.e.MapLeftOr(defaultValue, fn)
}

// MapOrElse returns r's value applied by fn if r is ok.
// Otherwise it returns the result of calling defaultFn with r's error.
func (r Result[T]) MapOrElse[U any](defaultFn func(err error) U, fn func(t T) U) U {
	return r.e.Fold(fn, defaultFn)
}

// Ok converts r into opt.Option[T].
// The option is some wrapping the value if r is ok, none otherwise.
func (r Result[T]) Ok() opt.Option[T] {
	return r.e.Left()
}

// Or returns r if r is ok, otherwise u.
func (r Result[T]) Or(u Result[T]) Result[T] {
	return Result[T]{e: r.e.RightAnd(u.e)}
}

// OrElse returns r if r is ok, otherwise calls fn with r's error and returns the result.
func (r Result[T]) OrElse(fn func(err error) Result[T]) Result[T] {
	return Result[T]{e: r.e.RightAndThen(func(err error) Either[T, error] { return fn(err).e })}
}

// Unwrap returns r's value if r is ok.
// Otherwise it panics with the held error.
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic(r.e.UnwrapRight())
	}
	return r.e.UnwrapLeft()
}

// UnwrapErr returns r's error if r is err.
// Otherwise it panics.
func (r Result[T]) UnwrapErr() error {
	if r.IsOk() {
		panic("either: UnwrapErr called on an ok Result")
	}
	return r.e.UnwrapRight()
}

// UnwrapOr returns r's value if r is ok, otherwise t.
func (r Result[T]) UnwrapOr(t T) T {
	return r.e.LeftOr(t)
}

// UnwrapOrDefault returns r's value if r is ok, otherwise the zero value of T.
func (r Result[T]) UnwrapOrDefault() T {
	var zero T
	return r.e.LeftOr(zero)
}

// UnwrapOrElse returns r's value if r is ok, otherwise the result of calling fn with r's error.
func (r Result[T]) UnwrapOrElse(fn func(err error) T) T {
	return r.e.LeftOrElse(fn)
}

// Value returns r's value as T.
// T would be zero value if r is err.
func (r Result[T]) Value() T {
	var zero T
	return r.e.LeftOr(zero)
}
