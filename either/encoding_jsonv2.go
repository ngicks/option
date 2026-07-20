package either

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"fmt"
)

var (
	_ jsonv2.MarshalerTo     = Either[any, any]{}
	_ jsonv2.UnmarshalerFrom = (*Either[any, any])(nil)
	_ jsonv2.MarshalerTo     = Result[any]{}
	_ jsonv2.UnmarshalerFrom = (*Result[any])(nil)
)

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
//
// It encodes e as a JSON object with a single "left" or "right" key tagging the held value.
// For the untagged representation, use [UntaggedEither].
func (e Either[L, R]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if e.IsRight() {
		if err := enc.WriteToken(jsontext.String("right")); err != nil {
			return err
		}
		if err := jsonv2.MarshalEncode(enc, &e.r); err != nil {
			return err
		}
	} else {
		if err := enc.WriteToken(jsontext.String("left")); err != nil {
			return err
		}
		if err := jsonv2.MarshalEncode(enc, &e.l); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
//
// It decodes a JSON object with a single "left" or "right" key,
// storing the tagged value on the corresponding side.
func (e *Either[L, R]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	key, err := readSingleKeyObjectStart(dec, "Either", "left", "right")
	if err != nil {
		return err
	}

	// Decode into a fresh value so that e keeps a valid state
	// even if decoding fails partway through.
	var decoded Either[L, R]
	switch key {
	case "left":
		var l L
		if err := jsonv2.UnmarshalDecode(dec, &l); err != nil {
			return err
		}
		decoded = Left[L, R](l)
	case "right":
		var r R
		if err := jsonv2.UnmarshalDecode(dec, &r); err != nil {
			return err
		}
		decoded = Right[L](r)
	}

	if err := readSingleKeyObjectEnd(dec, "Either", "left", "right"); err != nil {
		return err
	}

	*e = decoded
	return nil
}

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
//
// It encodes r as a JSON object with a single "ok" or "err" key tagging the held value.
// The error is encoded as its message string, or the JSON null if it is nil.
// For the untagged representation, use [UntaggedResult].
func (r Result[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if err := enc.WriteToken(jsontext.BeginObject); err != nil {
		return err
	}
	if r.IsErr() {
		if err := enc.WriteToken(jsontext.String("err")); err != nil {
			return err
		}
		msg := errMessage(r.e.r)
		if err := jsonv2.MarshalEncode(enc, &msg); err != nil {
			return err
		}
	} else {
		if err := enc.WriteToken(jsontext.String("ok")); err != nil {
			return err
		}
		if err := jsonv2.MarshalEncode(enc, &r.e.l); err != nil {
			return err
		}
	}
	return enc.WriteToken(jsontext.EndObject)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
//
// It decodes a JSON object with a single "ok" or "err" key.
// The "err" value must be a string or null; since the concrete error type is not
// recoverable, the error is reconstructed from the message with [errors.New].
func (r *Result[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	key, err := readSingleKeyObjectStart(dec, "Result", "ok", "err")
	if err != nil {
		return err
	}

	// Decode into a fresh value so that r keeps a valid state
	// even if decoding fails partway through.
	var decoded Result[T]
	switch key {
	case "ok":
		var t T
		if err := jsonv2.UnmarshalDecode(dec, &t); err != nil {
			return err
		}
		decoded = Ok(t)
	case "err":
		var msg *string
		if err := jsonv2.UnmarshalDecode(dec, &msg); err != nil {
			return err
		}
		decoded = Err[T](errFromMessage(msg))
	}

	if err := readSingleKeyObjectEnd(dec, "Result", "ok", "err"); err != nil {
		return err
	}

	*r = decoded
	return nil
}

// readSingleKeyObjectStart reads the begin-object token and the object name of
// a single-key tagged object, validating the name is key1 or key2.
func readSingleKeyObjectStart(dec *jsontext.Decoder, typ, key1, key2 string) (string, error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return "", err
	}
	if tok.Kind() != jsontext.KindBeginObject {
		return "", fmt.Errorf("either: tagged %s expects a JSON object, got %q", typ, tok.Kind())
	}

	tok, err = dec.ReadToken()
	if err != nil {
		return "", err
	}
	if tok.Kind() == jsontext.KindEndObject {
		return "", fmt.Errorf("either: tagged %s expects a JSON object with a single %q or %q key, got an empty object", typ, key1, key2)
	}
	key := tok.String()
	if key != key1 && key != key2 {
		return "", fmt.Errorf("either: tagged %s expects a %q or %q key, got %q", typ, key1, key2, key)
	}
	return key, nil
}

// readSingleKeyObjectEnd reads the end-object token of a single-key tagged object,
// erroring out if the object has more keys.
func readSingleKeyObjectEnd(dec *jsontext.Decoder, typ, key1, key2 string) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if tok.Kind() != jsontext.KindEndObject {
		return fmt.Errorf("either: tagged %s expects a JSON object with a single %q or %q key, got more", typ, key1, key2)
	}
	return nil
}
