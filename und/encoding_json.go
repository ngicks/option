package und

import "encoding/json"

var (
	_ json.Marshaler   = Und[any]{}
	_ json.Unmarshaler = (*Und[any])(nil)
)

// MarshalJSON implements json.Marshaler.
//
// It encodes a defined Und as its inner value and both null and undefined Und as the JSON null.
func (u Und[T]) MarshalJSON() ([]byte, error) {
	if !u.IsDefined() {
		return []byte(`null`), nil
	}
	return json.Marshal(u.opt.Value().Value())
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It decodes the JSON null into a null Und and any other JSON value into a defined Und.
func (u *Und[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*u = Null[T]()
		return nil
	}

	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}

	*u = Defined(t)
	return nil
}
