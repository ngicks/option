package opt

import "log/slog"

var (
	_ slog.LogValuer = Option[any]{}
)

// LogValue implements slog.LogValuer
func (o Option[T]) LogValue() slog.Value {
	if o.IsNone() {
		return slog.AnyValue(nil)
	}
	return slog.AnyValue(o.Value())
}
