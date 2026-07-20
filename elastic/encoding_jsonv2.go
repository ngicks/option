package elastic

import (
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"

	"github.com/ngicks/option/opt"
)

var (
	_ jsonv2.MarshalerTo     = Elastic[any]{}
	_ jsonv2.UnmarshalerFrom = (*Elastic[any])(nil)
)

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
//
// It encodes a defined Elastic as the JSON array of its options, none options as the JSON null,
// even if it was unmarshaled from a single value,
// and both null and undefined Elastic as the JSON null.
// Attach the `json:",omitzero"` option to an Elastic struct field to omit undefined values entirely.
func (e Elastic[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if !e.IsDefined() {
		return enc.WriteToken(jsontext.Null)
	}
	options := e.v.Value()
	if options == nil {
		// keep defined state encoded as an array, not null.
		options = opt.Options[T]{}
	}
	return jsonv2.MarshalEncode(enc, &options)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
//
// It decodes the JSON null into a null Elastic,
// a JSON array into a defined Elastic holding its elements as options, null elements as none,
// and any other JSON value into a defined Elastic holding it as the single some option.
func (e *Elastic[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case jsontext.KindNull:
		if err := dec.SkipValue(); err != nil {
			return err
		}
		*e = Null[T]()
		return nil
	case jsontext.KindBeginArray:
		// Decode into a fresh value so that e keeps a valid state
		// even if decoding fails partway through.
		var options opt.Options[T]
		if err := jsonv2.UnmarshalDecode(dec, &options); err != nil {
			return err
		}
		*e = FromOptions(options...)
		return nil
	default:
		var t T
		if err := jsonv2.UnmarshalDecode(dec, &t); err != nil {
			return err
		}
		*e = FromValue(t)
		return nil
	}
}
