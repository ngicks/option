package und

import "log/slog"

var (
	_ slog.LogValuer = Und[any]{}
)

// LogValue implements slog.LogValuer.
func (u Und[T]) LogValue() slog.Value {
	return u.opt.Value().LogValue()
}
