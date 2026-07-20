package either

import (
	"encoding/json"
	"encoding/json/jsontext"
	jsonv2 "encoding/json/v2"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

var (
	_ json.Marshaler         = UntaggedEither[any, any]{}
	_ json.Unmarshaler       = (*UntaggedEither[any, any])(nil)
	_ jsonv2.MarshalerTo     = UntaggedEither[any, any]{}
	_ jsonv2.UnmarshalerFrom = (*UntaggedEither[any, any])(nil)
	_ xml.Marshaler          = UntaggedEither[any, any]{}
	_ xml.Unmarshaler        = (*UntaggedEither[any, any])(nil)
	_ json.Marshaler         = UntaggedResult[any]{}
	_ json.Unmarshaler       = (*UntaggedResult[any])(nil)
	_ jsonv2.MarshalerTo     = UntaggedResult[any]{}
	_ jsonv2.UnmarshalerFrom = (*UntaggedResult[any])(nil)
	_ xml.Marshaler          = UntaggedResult[any]{}
	_ xml.Unmarshaler        = (*UntaggedResult[any])(nil)
)

// tokenReplayer replays recorded XML tokens as an xml.TokenReader.
type tokenReplayer struct {
	tokens []xml.Token
	i      int
}

func (r *tokenReplayer) Token() (xml.Token, error) {
	if r.i >= len(r.tokens) {
		return nil, io.EOF
	}
	tok := r.tokens[r.i]
	r.i++
	return tok, nil
}

// recordElement consumes the element opened by start from d,
// returning copies of all its tokens, start and its matching end element included.
func recordElement(d *xml.Decoder, start xml.StartElement) ([]xml.Token, error) {
	tokens := []xml.Token{xml.CopyToken(start)}
	for depth := 1; depth > 0; {
		tok, err := d.Token()
		if err != nil {
			return nil, err
		}
		switch tok.(type) {
		case xml.StartElement:
			depth++
		case xml.EndElement:
			depth--
		}
		tokens = append(tokens, xml.CopyToken(tok))
	}
	return tokens, nil
}

// decodeRecorded decodes an element recorded by [recordElement] into a V,
// replaying the tokens through a fresh decoder.
func decodeRecorded[V any](tokens []xml.Token) (V, error) {
	var v V
	dec := xml.NewTokenDecoder(&tokenReplayer{tokens: tokens})
	tok, err := dec.Token()
	if err != nil {
		return v, err
	}
	start := tok.(xml.StartElement)
	err = dec.DecodeElement(&v, &start)
	return v, err
}

// UntaggedEither[L, R] adapts Either[L, R] to the untagged JSON representation:
// the held value is encoded directly, without the "left" or "right" tag object.
//
// Since the encoding is untagged, unmarshaling first attempts to decode the value into L,
// then falls back to R; the first success wins.
// If L and R overlap (e.g. both decode from a JSON number), the value always lands on the left.
type UntaggedEither[L, R any] struct {
	Either[L, R]
}

// MarshalJSON implements json.Marshaler.
func (u UntaggedEither[L, R]) MarshalJSON() ([]byte, error) {
	if u.IsRight() {
		return json.Marshal(u.r)
	}
	return json.Marshal(u.l)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *UntaggedEither[L, R]) UnmarshalJSON(data []byte) error {
	var l L
	if err := json.Unmarshal(data, &l); err == nil {
		u.Either = Left[L, R](l)
		return nil
	}
	var r R
	if err := json.Unmarshal(data, &r); err != nil {
		return fmt.Errorf("either: value unmarshals as neither %T nor %T: %w", l, r, err)
	}
	u.Either = Right[L](r)
	return nil
}

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
func (u UntaggedEither[L, R]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if u.IsRight() {
		return jsonv2.MarshalEncode(enc, &u.r)
	}
	return jsonv2.MarshalEncode(enc, &u.l)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
func (u *UntaggedEither[L, R]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	var l L
	if err := jsonv2.Unmarshal(val, &l, dec.Options()); err == nil {
		u.Either = Left[L, R](l)
		return nil
	}
	var r R
	if err := jsonv2.Unmarshal(val, &r, dec.Options()); err != nil {
		return fmt.Errorf("either: value unmarshals as neither %T nor %T: %w", l, r, err)
	}
	u.Either = Right[L](r)
	return nil
}

// MarshalXML implements xml.Marshaler.
//
// It encodes the held value directly into the start element, without the tagging child element.
func (u UntaggedEither[L, R]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if u.IsRight() {
		return enc.EncodeElement(u.r, start)
	}
	return enc.EncodeElement(u.l, start)
}

// UnmarshalXML implements xml.Unmarshaler.
//
// It records the element's tokens and replays them,
// first attempting to decode the element into L, then falling back to R;
// the first success wins.
//
// Mind that XML decoding is lenient and the sides overlap more often than in JSON:
// decoding into a struct virtually never fails, and empty content decodes
// into primitives as their zero value. In such cases the value lands on the left;
// use a strict xml.Unmarshaler implementor as L to make the fallback reachable.
func (u *UntaggedEither[L, R]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tokens, err := recordElement(d, start)
	if err != nil {
		return err
	}

	if l, err := decodeRecorded[L](tokens); err == nil {
		u.Either = Left[L, R](l)
		return nil
	}
	r, err := decodeRecorded[R](tokens)
	if err != nil {
		var l L
		return fmt.Errorf("either: element unmarshals as neither %T nor %T: %w", l, r, err)
	}
	u.Either = Right[L](r)
	return nil
}

// UntaggedResult[T] adapts Result[T] to the untagged JSON representation:
// the held value is encoded directly, without the "ok" or "err" tag object.
// The error is encoded as its message string, or the JSON null if it is nil.
//
// Since the encoding is untagged, unmarshaling first attempts to decode the value into T,
// then falls back to an error message string reconstructed with [errors.New];
// the first success wins. If T also decodes from a JSON string, the value always lands on ok.
type UntaggedResult[T any] struct {
	Result[T]
}

// MarshalJSON implements json.Marshaler.
func (u UntaggedResult[T]) MarshalJSON() ([]byte, error) {
	if u.IsErr() {
		return json.Marshal(errMessage(u.Result.e.r))
	}
	return json.Marshal(u.Result.e.l)
}

// UnmarshalJSON implements json.Unmarshaler.
func (u *UntaggedResult[T]) UnmarshalJSON(data []byte) error {
	var t T
	if err := json.Unmarshal(data, &t); err == nil {
		u.Result = Ok(t)
		return nil
	}
	var msg *string
	if err := json.Unmarshal(data, &msg); err != nil {
		return fmt.Errorf("either: value unmarshals as neither %T nor an error message string: %w", t, err)
	}
	u.Result = Err[T](errFromMessage(msg))
	return nil
}

// MarshalJSONTo implements [encoding/json/v2.MarshalerTo].
func (u UntaggedResult[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if u.IsErr() {
		msg := errMessage(u.Result.e.r)
		return jsonv2.MarshalEncode(enc, &msg)
	}
	return jsonv2.MarshalEncode(enc, &u.Result.e.l)
}

// UnmarshalJSONFrom implements [encoding/json/v2.UnmarshalerFrom].
func (u *UntaggedResult[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}

	var t T
	if err := jsonv2.Unmarshal(val, &t, dec.Options()); err == nil {
		u.Result = Ok(t)
		return nil
	}
	var msg *string
	if err := jsonv2.Unmarshal(val, &msg, dec.Options()); err != nil {
		return fmt.Errorf("either: value unmarshals as neither %T nor an error message string: %w", t, err)
	}
	u.Result = Err[T](errFromMessage(msg))
	return nil
}

// MarshalXML implements xml.Marshaler.
//
// It encodes the held value directly into the start element, without the tagging child element.
// The error is encoded as its message string, an empty one if it is nil.
func (u UntaggedResult[T]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if u.IsErr() {
		var msg string
		if err := u.Result.e.r; err != nil {
			msg = err.Error()
		}
		return enc.EncodeElement(msg, start)
	}
	return enc.EncodeElement(u.Result.e.l, start)
}

// UnmarshalXML implements xml.Unmarshaler.
//
// It records the element's tokens and replays them,
// first attempting to decode the element into T,
// then falling back to an error message string reconstructed with [errors.New];
// the first success wins.
//
// Mind that XML decoding is lenient and the sides overlap more often than in JSON:
// decoding into a struct virtually never fails, and empty content decodes
// into primitives as their zero value. In such cases the value lands on ok;
// use a strict xml.Unmarshaler implementor as T to make the fallback reachable.
func (u *UntaggedResult[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	tokens, err := recordElement(d, start)
	if err != nil {
		return err
	}

	if t, err := decodeRecorded[T](tokens); err == nil {
		u.Result = Ok(t)
		return nil
	}
	msg, err := decodeRecorded[string](tokens)
	if err != nil {
		var t T
		return fmt.Errorf("either: element unmarshals as neither %T nor an error message string: %w", t, err)
	}
	u.Result = Err[T](errors.New(msg))
	return nil
}
