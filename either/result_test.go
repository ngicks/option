package either_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/ngicks/option/either"
	"github.com/ngicks/option/opt"
)

func TestResult(t *testing.T) {
	ok := either.Ok(5)
	if !ok.IsOk() || ok.IsErr() || ok.Value() != 5 || ok.Err() != nil {
		t.Errorf("Ok(5): ok=%t err=%v value=%d", ok.IsOk(), ok.Err(), ok.Value())
	}

	sentinel := errors.New("boom")
	e := either.Err[int](sentinel)
	if e.IsOk() || !e.IsErr() || !errors.Is(e.Err(), sentinel) {
		t.Errorf("Err(boom): ok=%t err=%v", e.IsOk(), e.Err())
	}

	if r := either.FromResult(strconv.Atoi("5")); r.Err() != nil || r.Value() != 5 {
		t.Errorf("FromResult(Atoi(5)) = %d, %v", r.Value(), r.Err())
	}
	if r := either.FromResult(strconv.Atoi("x")); r.Err() == nil {
		t.Error("FromResult(Atoi(x)) has no error")
	}

	if got := e.UnwrapOr(42); got != 42 {
		t.Errorf("Err.UnwrapOr(42) = %d, want 42", got)
	}
	if got := ok.Unwrap(); got != 5 {
		t.Errorf("Ok(5).Unwrap() = %d, want 5", got)
	}
}

func TestResult_Combinators(t *testing.T) {
	sentinel := errors.New("boom")
	ok := either.Ok(5)
	e := either.Err[int](sentinel)

	if got := ok.And(either.Ok("foo")); !got.IsOk() || got.Value() != "foo" {
		t.Errorf("Ok(5).And(Ok(foo)) = %v", got)
	}
	if got := e.And(either.Ok("foo")); !got.IsErr() || !errors.Is(got.Err(), sentinel) {
		t.Errorf("Err.And(Ok(foo)) = %v, want err(boom)", got)
	}
	if got := ok.AndThen(func(i int) either.Result[string] { return either.Ok(strconv.Itoa(i)) }); got.Value() != "5" {
		t.Errorf("Ok(5).AndThen(itoa) = %v, want ok(5)", got)
	}

	if got := ok.Map(strconv.Itoa); !got.IsOk() || got.Value() != "5" {
		t.Errorf("Ok(5).Map(itoa) = %v, want ok(5)", got)
	}
	if got := e.Map(strconv.Itoa); !got.IsErr() || !errors.Is(got.Err(), sentinel) {
		t.Errorf("Err.Map(itoa) = %v, want err(boom)", got)
	}
	wrapped := errors.New("wrapped")
	if got := e.MapErr(func(error) error { return wrapped }); !errors.Is(got.Err(), wrapped) {
		t.Errorf("Err.MapErr = %v, want err(wrapped)", got)
	}
	if got := ok.MapOr("x", strconv.Itoa); got != "5" {
		t.Errorf("Ok(5).MapOr = %s, want 5", got)
	}
	if got := e.MapOrElse(func(err error) string { return err.Error() }, strconv.Itoa); got != "boom" {
		t.Errorf("Err.MapOrElse = %s, want boom", got)
	}

	if got := e.Or(ok); !got.IsOk() || got.Value() != 5 {
		t.Errorf("Err.Or(Ok(5)) = %v, want ok(5)", got)
	}
	if got := e.OrElse(func(err error) either.Result[int] { return either.Ok(1) }); got.Value() != 1 {
		t.Errorf("Err.OrElse = %v, want ok(1)", got)
	}

	if !ok.IsOkAnd(func(i int) bool { return i == 5 }) || ok.IsOkAnd(func(i int) bool { return i == 6 }) {
		t.Error("IsOkAnd misbehaves on Ok(5)")
	}
	if !e.IsErrAnd(func(err error) bool { return errors.Is(err, sentinel) }) || ok.IsErrAnd(func(error) bool { return true }) {
		t.Error("IsErrAnd misbehaves")
	}

	inspected := 0
	ok.Inspect(func(i int) { inspected = i }).InspectErr(func(error) { inspected = -1 })
	if inspected != 5 {
		t.Errorf("Inspect/InspectErr: inspected = %d, want 5", inspected)
	}

	var collected []int
	for v := range ok.Iter() {
		collected = append(collected, v)
	}
	for range e.Iter() {
		t.Error("Err.Iter() yielded a value")
	}
	if len(collected) != 1 || collected[0] != 5 {
		t.Errorf("Ok(5).Iter() collected %v, want [5]", collected)
	}

	if got := ok.Expect("nope"); got != 5 {
		t.Errorf("Ok(5).Expect = %d, want 5", got)
	}
	if err := e.ExpectErr("nope"); !errors.Is(err, sentinel) {
		t.Errorf("Err.ExpectErr = %v, want boom", err)
	}
	if err := e.UnwrapErr(); !errors.Is(err, sentinel) {
		t.Errorf("Err.UnwrapErr = %v, want boom", err)
	}
	if got := e.UnwrapOrDefault(); got != 0 {
		t.Errorf("Err.UnwrapOrDefault = %d, want 0", got)
	}
	if got := e.UnwrapOrElse(func(err error) int { return len(err.Error()) }); got != 4 {
		t.Errorf("Err.UnwrapOrElse = %d, want 4", got)
	}

	t.Run("Unwrap on err panics with the error", func(t *testing.T) {
		defer func() {
			if got := recover(); got == nil || !errors.Is(got.(error), sentinel) {
				t.Errorf("panic value = %v, want boom", got)
			}
		}()
		e.Unwrap()
	})
	t.Run("UnwrapErr on ok panics", func(t *testing.T) {
		defer func() {
			if recover() == nil {
				t.Error("UnwrapErr on ok did not panic")
			}
		}()
		ok.UnwrapErr()
	})
}

func TestResult_Ok(t *testing.T) {
	if got := either.Ok(5).Ok(); !got.IsSome() || got.Value() != 5 {
		t.Errorf("Ok(5).Ok() = %v, want some(5)", got)
	}
	if got := either.Err[int](errors.New("boom")).Ok(); got.IsSome() {
		t.Errorf("Err(boom).Ok() = %v, want none", got)
	}
}

func TestResult_ErrNilPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Err(nil) did not panic")
		}
	}()
	either.Err[int](nil)
}

func TestResult_Flatten_Transpose(t *testing.T) {
	sentinel := errors.New("boom")

	if got := either.Flatten(either.Ok(either.Ok(5))); !got.IsOk() || got.Value() != 5 {
		t.Errorf("Flatten(Ok(Ok(5))) = %v, want ok(5)", got)
	}
	if got := either.Flatten(either.Ok(either.Err[int](sentinel))); !got.IsErr() || !errors.Is(got.Err(), sentinel) {
		t.Errorf("Flatten(Ok(Err(boom))) = %v, want err(boom)", got)
	}
	if got := either.Flatten(either.Err[either.Result[int]](sentinel)); !got.IsErr() || !errors.Is(got.Err(), sentinel) {
		t.Errorf("Flatten(Err(boom)) = %v, want err(boom)", got)
	}

	if got := either.Transpose(either.Ok(opt.Some(5))); !got.IsSome() || got.Value().Value() != 5 {
		t.Errorf("Transpose(Ok(Some(5))) = %v, want some(ok(5))", got)
	}
	if got := either.Transpose(either.Ok(opt.None[int]())); got.IsSome() {
		t.Errorf("Transpose(Ok(None)) = %v, want none", got)
	}
	got := either.Transpose(either.Err[opt.Option[int]](sentinel))
	if !got.IsSome() || !errors.Is(got.Value().Err(), sentinel) {
		t.Errorf("Transpose(Err(boom)) = %v, want some(err(boom))", got)
	}
}
