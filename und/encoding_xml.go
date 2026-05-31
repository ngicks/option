package und

import "encoding/xml"

var (
	_ xml.Marshaler   = Und[any]{}
	_ xml.Unmarshaler = (*Und[any])(nil)
)

// MarshalXML implements xml.Marshaler.
//
// A defined Und encodes its inner value; null and undefined Und encode nothing.
func (u Und[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return u.opt.Value().MarshalXML(e, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (u *Und[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var t T
	err := d.DecodeElement(&t, &start)
	if err != nil {
		return err
	}

	*u = Defined(t)

	return nil
}
