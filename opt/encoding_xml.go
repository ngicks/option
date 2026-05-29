package opt

var (
	_ xml.Marshaler    = Option[any]{}
	_ xml.Unmarshaler  = (*Option[any])(nil)
)
