package either

import "github.com/ngicks/option/opt"

// Left converts e into opt.Option[L].
// The option is some wrapping the left value if e is left, none otherwise.
func (e Either[L, R]) Left() opt.Option[L] {
	if e.IsRight() {
		return opt.None[L]()
	}
	return opt.Some(e.l)
}

// Right converts e into opt.Option[R].
// The option is some wrapping the right value if e is right, none otherwise.
func (e Either[L, R]) Right() opt.Option[R] {
	if e.IsLeft() {
		return opt.None[R]()
	}
	return opt.Some(e.r)
}

// TransposeLeft converts an Either[opt.Option[L], R] into an opt.Option[Either[L, R]].
// A left none maps to none.
// A left some maps to some wrapping left, and a right maps to some wrapping right.
func TransposeLeft[L, R any](e Either[opt.Option[L], R]) opt.Option[Either[L, R]] {
	if e.IsRight() {
		return opt.Some(Right[L](e.r))
	}
	if e.l.IsNone() {
		return opt.None[Either[L, R]]()
	}
	return opt.Some(Left[L, R](e.l.Value()))
}

// TransposeRight converts an Either[L, opt.Option[R]] into an opt.Option[Either[L, R]].
// A right none maps to none.
// A right some maps to some wrapping right, and a left maps to some wrapping left.
func TransposeRight[L, R any](e Either[L, opt.Option[R]]) opt.Option[Either[L, R]] {
	if e.IsLeft() {
		return opt.Some(Left[L, R](e.l))
	}
	if e.r.IsNone() {
		return opt.None[Either[L, R]]()
	}
	return opt.Some(Right[L](e.r.Value()))
}
