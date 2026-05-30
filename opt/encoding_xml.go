package opt

import "encoding/xml"

var (
	_ xml.Marshaler   = Option[any]{}
	_ xml.Unmarshaler = (*Option[any])(nil)
)

func (o Option[T]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if o.IsNone() {
		return nil
	}
	return e.EncodeElement(o.Value(), start)
}

func (o *Option[T]) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var t T
	err := d.DecodeElement(&t, &start)
	if err != nil {
		return err
	}

	o.some = true
	o.v = t

	return nil
}
