//go:build goexperiment.jsonv2

package opt

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

var (
	_ jsonv2.MarshalerTo     = Option[any]{}
	_ jsonv2.UnmarshalerFrom = (*Option[any])(nil)
)

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
//
// It encodes a none Option as the JSON null and a some Option as its inner value.
func (o Option[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.IsNone() {
		return enc.WriteToken(jsontext.Null)
	}
	return jsonv2.MarshalEncode(enc, &o.v)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
//
// It decodes the JSON null as a none Option and any other JSON value as a some Option.
func (o *Option[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == jsontext.KindNull {
		if err := dec.SkipValue(); err != nil {
			return err
		}
		o.some = false
		var zero T
		o.v = zero
		return nil
	}

	// Decode into a fresh value so that o keeps a valid state
	// even if decoding fails partway through.
	var v T
	if err := jsonv2.UnmarshalDecode(dec, &v); err != nil {
		return err
	}
	o.some = true
	o.v = v
	return nil
}
