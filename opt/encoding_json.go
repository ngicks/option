package opt

var (
	_ json.Marshaler   = Option[any]{}
	_ json.Unmarshaler = (*Option[any])(nil)
)
