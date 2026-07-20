package either_test

import (
	"strconv"
	"testing"

	"github.com/ngicks/option/either"
	"github.com/ngicks/option/opt"
)

func TestEither_States(t *testing.T) {
	l := either.Left[int, string](5)
	if !l.IsLeft() || l.IsRight() {
		t.Errorf("Left(5): left=%t right=%t", l.IsLeft(), l.IsRight())
	}
	if l.UnwrapLeft() != 5 {
		t.Errorf("Left(5).UnwrapLeft() = %d, want 5", l.UnwrapLeft())
	}

	r := either.Right[int]("foo")
	if r.IsLeft() || !r.IsRight() {
		t.Errorf("Right(foo): left=%t right=%t", r.IsLeft(), r.IsRight())
	}
	if r.UnwrapRight() != "foo" {
		t.Errorf("Right(foo).UnwrapRight() = %s, want foo", r.UnwrapRight())
	}

	// Zero value of Either is left holding the zero value of L.
	var zero either.Either[int, string]
	if !zero.IsLeft() || zero.UnwrapLeft() != 0 {
		t.Errorf("zero Either: left=%t value=%d", zero.IsLeft(), zero.UnwrapLeft())
	}
}

func TestEither_LeftRightOptions(t *testing.T) {
	l := either.Left[int, string](5)
	if o := l.Left(); !o.IsSome() || o.Value() != 5 {
		t.Errorf("Left(5).Left() = %v, want some(5)", o)
	}
	if o := l.Right(); o.IsSome() {
		t.Errorf("Left(5).Right() = %v, want none", o)
	}

	r := either.Right[int]("foo")
	if o := r.Left(); o.IsSome() {
		t.Errorf("Right(foo).Left() = %v, want none", o)
	}
	if o := r.Right(); !o.IsSome() || o.Value() != "foo" {
		t.Errorf("Right(foo).Right() = %v, want some(foo)", o)
	}
}

func TestEither_Fold_Map_Flip(t *testing.T) {
	l := either.Left[int, string](5)
	r := either.Right[int]("foo")

	if got := l.Fold(strconv.Itoa, func(s string) string { return s }); got != "5" {
		t.Errorf("Left(5).Fold = %s, want 5", got)
	}
	if got := r.Fold(strconv.Itoa, func(s string) string { return s }); got != "foo" {
		t.Errorf("Right(foo).Fold = %s, want foo", got)
	}

	if got := l.MapLeft(strconv.Itoa); !got.IsLeft() || got.UnwrapLeft() != "5" {
		t.Errorf("Left(5).MapLeft(itoa) = %v, want left(5)", got)
	}
	if got := r.MapLeft(strconv.Itoa); !got.IsRight() || got.UnwrapRight() != "foo" {
		t.Errorf("Right(foo).MapLeft(itoa) = %v, want right(foo)", got)
	}
	if got := r.MapRight(func(s string) int { return len(s) }); !got.IsRight() || got.UnwrapRight() != 3 {
		t.Errorf("Right(foo).MapRight(len) = %v, want right(3)", got)
	}

	if got := l.Flip(); !got.IsRight() || got.UnwrapRight() != 5 {
		t.Errorf("Left(5).Flip() = %v, want right(5)", got)
	}
	if got := r.Flip(); !got.IsLeft() || got.UnwrapLeft() != "foo" {
		t.Errorf("Right(foo).Flip() = %v, want left(foo)", got)
	}
}

func TestEither_Or(t *testing.T) {
	l := either.Left[int, string](5)
	r := either.Right[int]("foo")

	if got := l.LeftOr(99); got != 5 {
		t.Errorf("Left(5).LeftOr(99) = %d, want 5", got)
	}
	if got := r.LeftOr(99); got != 99 {
		t.Errorf("Right(foo).LeftOr(99) = %d, want 99", got)
	}
	if got := l.RightOr("bar"); got != "bar" {
		t.Errorf("Left(5).RightOr(bar) = %s, want bar", got)
	}
	if got := r.RightOrElse(strconv.Itoa); got != "foo" {
		t.Errorf("Right(foo).RightOrElse = %s, want foo", got)
	}
	if got := l.LeftOrElse(func(s string) int { return len(s) }); got != 5 {
		t.Errorf("Left(5).LeftOrElse = %d, want 5", got)
	}
}

func TestEither_SideCombinators(t *testing.T) {
	l := either.Left[int, string](5)
	r := either.Right[int]("foo")

	if got := l.LeftAnd(either.Left[bool, string](true)); !got.IsLeft() || !got.UnwrapLeft() {
		t.Errorf("Left(5).LeftAnd(Left(true)) = %v, want left(true)", got)
	}
	if got := r.LeftAnd(either.Left[bool, string](true)); !got.IsRight() || got.UnwrapRight() != "foo" {
		t.Errorf("Right(foo).LeftAnd(Left(true)) = %v, want right(foo)", got)
	}
	if got := l.LeftAndThen(func(i int) either.Either[bool, string] { return either.Left[bool, string](i == 5) }); !got.UnwrapLeft() {
		t.Errorf("Left(5).LeftAndThen = %v, want left(true)", got)
	}
	if got := l.RightAnd(either.Right[int](true)); !got.IsLeft() || got.UnwrapLeft() != 5 {
		t.Errorf("Left(5).RightAnd(Right(true)) = %v, want left(5)", got)
	}
	if got := r.RightAndThen(func(s string) either.Either[int, bool] { return either.Right[int](s == "foo") }); !got.UnwrapRight() {
		t.Errorf("Right(foo).RightAndThen = %v, want right(true)", got)
	}

	if !l.IsLeftAnd(func(i int) bool { return i == 5 }) || r.IsLeftAnd(func(int) bool { return true }) {
		t.Error("IsLeftAnd misbehaves")
	}
	if !r.IsRightAnd(func(s string) bool { return s == "foo" }) || l.IsRightAnd(func(string) bool { return true }) {
		t.Error("IsRightAnd misbehaves")
	}

	if got := l.MapLeftOr("x", strconv.Itoa); got != "5" {
		t.Errorf("Left(5).MapLeftOr = %s, want 5", got)
	}
	if got := l.MapRightOr(-1, func(s string) int { return len(s) }); got != -1 {
		t.Errorf("Left(5).MapRightOr = %d, want -1", got)
	}

	var lefts []int
	for v := range l.IterLeft() {
		lefts = append(lefts, v)
	}
	if len(lefts) != 1 || lefts[0] != 5 {
		t.Errorf("Left(5).IterLeft() collected %v, want [5]", lefts)
	}
	for range l.IterRight() {
		t.Error("Left(5).IterRight() yielded a value")
	}
	var rights []string
	for v := range r.IterRight() {
		rights = append(rights, v)
	}
	if len(rights) != 1 || rights[0] != "foo" {
		t.Errorf("Right(foo).IterRight() collected %v, want [foo]", rights)
	}
}

func TestEither_Equal(t *testing.T) {
	cases := []struct {
		name string
		l, r either.Either[int, string]
		want bool
	}{
		{"left eq", either.Left[int, string](5), either.Left[int, string](5), true},
		{"left ne", either.Left[int, string](5), either.Left[int, string](6), false},
		{"right eq", either.Right[int]("a"), either.Right[int]("a"), true},
		{"right ne", either.Right[int]("a"), either.Right[int]("b"), false},
		{"left vs right", either.Left[int, string](5), either.Right[int]("a"), false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := either.Equal(c.l, c.r); got != c.want {
				t.Errorf("Equal(%v, %v) = %t, want %t", c.l, c.r, got, c.want)
			}
		})
	}
}

func TestEither_UnwrapPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("UnwrapLeft on right did not panic")
		}
	}()
	either.Right[int]("foo").UnwrapLeft()
}

func TestEither_Flatten(t *testing.T) {
	innerL := either.Left[int, string](5)
	innerR := either.Right[int]("foo")

	if got := either.FlattenLeft(either.Left[either.Either[int, string], string](innerL)); !got.IsLeft() || got.UnwrapLeft() != 5 {
		t.Errorf("FlattenLeft(Left(Left(5))) = %v, want left(5)", got)
	}
	if got := either.FlattenLeft(either.Right[either.Either[int, string]]("bar")); !got.IsRight() || got.UnwrapRight() != "bar" {
		t.Errorf("FlattenLeft(Right(bar)) = %v, want right(bar)", got)
	}

	if got := either.FlattenRight(either.Right[int](innerR)); !got.IsRight() || got.UnwrapRight() != "foo" {
		t.Errorf("FlattenRight(Right(Right(foo))) = %v, want right(foo)", got)
	}
	if got := either.FlattenRight(either.Left[int, either.Either[int, string]](7)); !got.IsLeft() || got.UnwrapLeft() != 7 {
		t.Errorf("FlattenRight(Left(7)) = %v, want left(7)", got)
	}

	if got := either.FlattenEither(either.Left[either.Either[int, string], either.Either[int, string]](innerL)); !got.IsLeft() || got.UnwrapLeft() != 5 {
		t.Errorf("FlattenEither(Left(Left(5))) = %v, want left(5)", got)
	}
	if got := either.FlattenEither(either.Right[either.Either[int, string]](innerR)); !got.IsRight() || got.UnwrapRight() != "foo" {
		t.Errorf("FlattenEither(Right(Right(foo))) = %v, want right(foo)", got)
	}
}

func TestEither_Transpose(t *testing.T) {
	if got := either.TransposeLeft(either.Left[opt.Option[int], string](opt.Some(5))); !got.IsSome() || !got.Value().IsLeft() || got.Value().UnwrapLeft() != 5 {
		t.Errorf("TransposeLeft(Left(Some(5))) = %v, want some(left(5))", got)
	}
	if got := either.TransposeLeft(either.Left[opt.Option[int], string](opt.None[int]())); got.IsSome() {
		t.Errorf("TransposeLeft(Left(None)) = %v, want none", got)
	}
	if got := either.TransposeLeft(either.Right[opt.Option[int]]("foo")); !got.IsSome() || !got.Value().IsRight() || got.Value().UnwrapRight() != "foo" {
		t.Errorf("TransposeLeft(Right(foo)) = %v, want some(right(foo))", got)
	}

	if got := either.TransposeRight(either.Right[int](opt.Some("foo"))); !got.IsSome() || !got.Value().IsRight() || got.Value().UnwrapRight() != "foo" {
		t.Errorf("TransposeRight(Right(Some(foo))) = %v, want some(right(foo))", got)
	}
	if got := either.TransposeRight(either.Right[int](opt.None[string]())); got.IsSome() {
		t.Errorf("TransposeRight(Right(None)) = %v, want none", got)
	}
	if got := either.TransposeRight(either.Left[int, opt.Option[string]](5)); !got.IsSome() || !got.Value().IsLeft() || got.Value().UnwrapLeft() != 5 {
		t.Errorf("TransposeRight(Left(5)) = %v, want some(left(5))", got)
	}
}
