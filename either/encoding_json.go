package either

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ json.Marshaler   = Either[any, any]{}
	_ json.Unmarshaler = (*Either[any, any])(nil)
	_ json.Marshaler   = Result[any]{}
	_ json.Unmarshaler = (*Result[any])(nil)
)

// MarshalJSON implements json.Marshaler.
//
// It encodes e as a JSON object with a single "left" or "right" key tagging the held value.
// For the untagged representation, use [UntaggedEither].
func (e Either[L, R]) MarshalJSON() ([]byte, error) {
	if e.IsRight() {
		return json.Marshal(struct {
			R R `json:"right"`
		}{e.r})
	}
	return json.Marshal(struct {
		L L `json:"left"`
	}{e.l})
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It decodes a JSON object with a single "left" or "right" key,
// storing the tagged value on the corresponding side.
func (e *Either[L, R]) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("either: tagged Either expects a JSON object with a single %q or %q key, got %d keys", "left", "right", len(m))
	}
	for key, raw := range m {
		switch key {
		case "left":
			var l L
			if err := json.Unmarshal(raw, &l); err != nil {
				return err
			}
			*e = Left[L, R](l)
		case "right":
			var r R
			if err := json.Unmarshal(raw, &r); err != nil {
				return err
			}
			*e = Right[L](r)
		default:
			return fmt.Errorf("either: tagged Either expects a %q or %q key, got %q", "left", "right", key)
		}
	}
	return nil
}

// errMessage converts err into its message for encoding, nil error to nil.
func errMessage(err error) *string {
	if err == nil {
		return nil
	}
	return new(err.Error())
}

// errFromMessage converts a decoded message back into an error, nil message to nil.
// The concrete error type is not recoverable from an encoded form;
// the error is reconstructed with [errors.New].
func errFromMessage(msg *string) error {
	if msg == nil {
		return nil
	}
	return errors.New(*msg)
}

// MarshalJSON implements json.Marshaler.
//
// It encodes r as a JSON object with a single "ok" or "err" key tagging the held value.
// The error is encoded as its message string, or the JSON null if it is nil.
// For the untagged representation, use [UntaggedResult].
func (r Result[T]) MarshalJSON() ([]byte, error) {
	if r.IsErr() {
		return json.Marshal(struct {
			E *string `json:"err"`
		}{errMessage(r.e.r)})
	}
	return json.Marshal(struct {
		T T `json:"ok"`
	}{r.e.l})
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It decodes a JSON object with a single "ok" or "err" key.
// The "err" value must be a string or null; since the concrete error type is not
// recoverable, the error is reconstructed from the message with [errors.New].
func (r *Result[T]) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	if len(m) != 1 {
		return fmt.Errorf("either: tagged Result expects a JSON object with a single %q or %q key, got %d keys", "ok", "err", len(m))
	}
	for key, raw := range m {
		switch key {
		case "ok":
			var t T
			if err := json.Unmarshal(raw, &t); err != nil {
				return err
			}
			*r = Ok(t)
		case "err":
			var msg *string
			if err := json.Unmarshal(raw, &msg); err != nil {
				return err
			}
			*r = Err[T](errFromMessage(msg))
		default:
			return fmt.Errorf("either: tagged Result expects an %q or %q key, got %q", "ok", "err", key)
		}
	}
	return nil
}
