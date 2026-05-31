//go:build goexperiment.jsonv2

package und

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
)

var (
	_ jsonv2.MarshalerTo     = Und[any]{}
	_ jsonv2.UnmarshalerFrom = (*Und[any])(nil)
)

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
//
// It encodes a defined Und as its inner value and both null and undefined Und as the JSON null.
// Attach the `json:",omitzero"` option to an Und struct field to omit undefined values entirely.
func (u Und[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !u.IsDefined() {
		return enc.WriteToken(jsontext.Null)
	}
	v := u.opt.Value().Value()
	return jsonv2.MarshalEncode(enc, &v)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
//
// It decodes the JSON null into a null Und and any other JSON value into a defined Und.
func (u *Und[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() == jsontext.KindNull {
		if err := dec.SkipValue(); err != nil {
			return err
		}
		*u = Null[T]()
		return nil
	}

	// Decode into a fresh value so that u keeps a valid state
	// even if decoding fails partway through.
	var v T
	if err := jsonv2.UnmarshalDecode(dec, &v); err != nil {
		return err
	}
	*u = Defined(v)
	return nil
}
