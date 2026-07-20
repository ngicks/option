package elastic

import (
	"encoding/xml"

	"github.com/ngicks/option/opt"
)

var (
	_ xml.Marshaler   = Elastic[any]{}
	_ xml.Unmarshaler = (*Elastic[any])(nil)
)

// MarshalXML implements xml.Marshaler.
//
// A defined Elastic encodes each of its some options as an element;
// none options, null and undefined Elastic encode nothing.
func (e Elastic[T]) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	for _, o := range e.Options() {
		if err := enc.EncodeElement(o, start); err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalXML implements xml.Unmarshaler.
//
// It is called once per matching element; repeated elements accumulate into e's options.
func (e *Elastic[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var o opt.Option[T]
	err := d.DecodeElement(&o, &start)
	if err != nil {
		return err
	}

	if e.IsDefined() {
		*e = FromOptions(append(e.Options(), o)...)
	} else {
		*e = FromOptions(o)
	}

	return nil
}
