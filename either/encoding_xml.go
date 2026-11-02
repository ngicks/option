package either

import (
	"encoding/xml"
	"errors"
	"fmt"
)

var (
	_ xml.Marshaler   = Either[any, any]{}
	_ xml.Unmarshaler = (*Either[any, any])(nil)
	_ xml.Marshaler   = Result[any]{}
	_ xml.Unmarshaler = (*Result[any])(nil)
)

// marshalSingleKeyXML encodes value wrapped in a key-named element inside start.
func marshalSingleKeyXML(enc *xml.Encoder, start xml.StartElement, key string, value any) error {
	if err := enc.EncodeToken(start); err != nil {
		return err
	}
	if err := enc.EncodeElement(value, xml.StartElement{Name: xml.Name{Local: key}}); err != nil {
		return err
	}
	return enc.EncodeToken(start.End())
}

// unmarshalSingleKeyXML consumes the element opened by start,
// expecting exactly one child element whose local name is a key of decode,
// and calls the matching callback with that child's start element.
func unmarshalSingleKeyXML(
	d *xml.Decoder,
	typ, key1, key2 string,
	decode map[string]func(start xml.StartElement) error,
) error {
	seen := false
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			if seen {
				return fmt.Errorf(
					"either: tagged %s expects a single <%s> or <%s> element, got more",
					typ,
					key1,
					key2,
				)
			}
			fn, ok := decode[tok.Name.Local]
			if !ok {
				return fmt.Errorf(
					"either: tagged %s expects a <%s> or <%s> element, got <%s>",
					typ,
					key1,
					key2,
					tok.Name.Local,
				)
			}
			if err := fn(tok); err != nil {
				return err
			}
			seen = true
		case xml.EndElement:
			if !seen {
				return fmt.Errorf(
					"either: tagged %s expects a single <%s> or <%s> element, got none",
					typ,
					key1,
					key2,
				)
			}
			return nil
		}
	}
}

// MarshalXML implements xml.Marshaler.
//
// It encodes the held value wrapped in a <left> or <right> element tagging the side,
// e.g. <v><left>5</left></v>.
func (e Either[L, R]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if e.IsRight() {
		return marshalSingleKeyXML(enc, start, "right", e.r)
	}
	return marshalSingleKeyXML(enc, start, "left", e.l)
}

// UnmarshalXML implements xml.Unmarshaler.
//
// It decodes an element containing a single <left> or <right> child element,
// storing the tagged value on the corresponding side.
func (e *Either[L, R]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decode into a fresh value so that e keeps a valid state
	// even if decoding fails partway through.
	var decoded Either[L, R]
	err := unmarshalSingleKeyXML(
		d,
		"Either",
		"left",
		"right",
		map[string]func(start xml.StartElement) error{
			"left": func(start xml.StartElement) error {
				var l L
				if err := d.DecodeElement(&l, &start); err != nil {
					return err
				}
				decoded = Left[L, R](l)
				return nil
			},
			"right": func(start xml.StartElement) error {
				var r R
				if err := d.DecodeElement(&r, &start); err != nil {
					return err
				}
				decoded = Right[L](r)
				return nil
			},
		},
	)
	if err != nil {
		return err
	}
	*e = decoded
	return nil
}

// MarshalXML implements xml.Marshaler.
//
// It encodes the held value wrapped in an <ok> or <err> element tagging the side.
// The error is encoded as its message string, an empty one if it is nil.
func (r Result[T]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	if r.IsErr() {
		var msg string
		if err := r.e.r; err != nil {
			msg = err.Error()
		}
		return marshalSingleKeyXML(enc, start, "err", msg)
	}
	return marshalSingleKeyXML(enc, start, "ok", r.e.l)
}

// UnmarshalXML implements xml.Unmarshaler.
//
// It decodes an element containing a single <ok> or <err> child element.
// Since the concrete error type is not recoverable,
// the error is reconstructed from the message with [errors.New].
func (r *Result[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Decode into a fresh value so that r keeps a valid state
	// even if decoding fails partway through.
	var decoded Result[T]
	err := unmarshalSingleKeyXML(
		d,
		"Result",
		"ok",
		"err",
		map[string]func(start xml.StartElement) error{
			"ok": func(start xml.StartElement) error {
				var t T
				if err := d.DecodeElement(&t, &start); err != nil {
					return err
				}
				decoded = Ok(t)
				return nil
			},
			"err": func(start xml.StartElement) error {
				var msg string
				if err := d.DecodeElement(&msg, &start); err != nil {
					return err
				}
				decoded = Err[T](errors.New(msg))
				return nil
			},
		},
	)
	if err != nil {
		return err
	}
	*r = decoded
	return nil
}
