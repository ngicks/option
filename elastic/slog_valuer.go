package elastic

import "log/slog"

var (
	_ slog.LogValuer = Elastic[any]{}
)

// LogValue implements slog.LogValuer.
func (e Elastic[T]) LogValue() slog.Value {
	return e.v.LogValue()
}
