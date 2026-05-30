package opt

import "encoding/json"

var (
	_ json.Marshaler   = Option[any]{}
	_ json.Unmarshaler = (*Option[any])(nil)
)

func (o Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsNone() {
		// same as bytes.Clone.
		return []byte(`null`), nil
	}
	return json.Marshal(o.v)
}

func (o *Option[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.some = false
		var zero T
		o.v = zero
		return nil
	}

	// not gonna call like json.Unmarshal(data, &o.v)
	// since it could be half-baked result if unmarshaling fails at some point.
	// o's state must stay valid.
	var v T
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	o.some = true
	o.v = v
	return nil
}
