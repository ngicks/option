package either

import "log/slog"

var (
	_ slog.LogValuer = Either[any, any]{}
	_ slog.LogValuer = Result[any]{}
)

// LogValue implements slog.LogValuer.
func (e Either[L, R]) LogValue() slog.Value {
	if e.IsRight() {
		return slog.AnyValue(e.r)
	}
	return slog.AnyValue(e.l)
}

// LogValue implements slog.LogValuer, delegating to the internal either value.
func (r Result[T]) LogValue() slog.Value {
	return r.e.LogValue()
}
