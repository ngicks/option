package elastic

import (
	"bytes"
	"encoding/json"

	"github.com/ngicks/option/opt"
)

var (
	_ json.Marshaler   = Elastic[any]{}
	_ json.Unmarshaler = (*Elastic[any])(nil)
)

// MarshalJSON implements json.Marshaler.
//
// It encodes a defined Elastic as the JSON array of its options, none options as the JSON null,
// even if it was unmarshaled from a single value,
// and both null and undefined Elastic as the JSON null.
func (e Elastic[T]) MarshalJSON() ([]byte, error) {
	if !e.IsDefined() {
		return []byte(`null`), nil
	}
	options := e.v.Value()
	if options == nil {
		// keep defined state encoded as an array, not null.
		return []byte(`[]`), nil
	}
	return json.Marshal(options)
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It decodes the JSON null into a null Elastic,
// a JSON array into a defined Elastic holding its elements as options, null elements as none,
// and any other JSON value into a defined Elastic holding it as the single some option.
func (e *Elastic[T]) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimLeft(data, " \t\r\n")
	if string(trimmed) == "null" {
		*e = Null[T]()
		return nil
	}

	if len(trimmed) > 0 && trimmed[0] == '[' {
		var options opt.Options[T]
		err := json.Unmarshal(data, &options)
		if err != nil {
			return err
		}
		*e = FromOptions(options...)
		return nil
	}

	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return err
	}
	*e = FromValue(t)
	return nil
}
