package und_test

import (
	"database/sql"
	"testing"

	"github.com/ngicks/option/und"
)

func TestUnd_FromSqlNull(t *testing.T) {
	if u := und.FromSqlNull(sql.Null[int]{Valid: true, V: 7}); !u.IsDefined() || u.Value() != 7 {
		t.Errorf("FromSqlNull(valid 7) = %v", u)
	}
	if u := und.FromSqlNull(sql.Null[int]{}); !u.IsNull() {
		t.Errorf("FromSqlNull(invalid) = state %d, want null", u.State())
	}
}

func TestUnd_SqlNull_method(t *testing.T) {
	if n := und.Defined(7).SqlNull(); !n.Valid || n.V != 7 {
		t.Errorf("Defined(7).SqlNull() = %+v, want valid 7", n)
	}
	if n := und.Null[int]().SqlNull(); n.Valid {
		t.Errorf("Null.SqlNull() = %+v, want invalid", n)
	}
	if n := und.Undefined[int]().SqlNull(); n.Valid {
		t.Errorf("Undefined.SqlNull() = %+v, want invalid", n)
	}
}

func TestUnd_SqlNull_Value(t *testing.T) {
	v, err := (und.SqlNull[int64]{Und: und.Defined[int64](5)}).Value()
	if err != nil || v != int64(5) {
		t.Errorf("defined Value() = %v, %v; want 5, nil", v, err)
	}
	v, err = (und.SqlNull[int64]{Und: und.Null[int64]()}).Value()
	if err != nil || v != nil {
		t.Errorf("null Value() = %v, %v; want nil, nil", v, err)
	}
	v, err = (und.SqlNull[int64]{Und: und.Undefined[int64]()}).Value()
	if err != nil || v != nil {
		t.Errorf("undefined Value() = %v, %v; want nil, nil", v, err)
	}
}

func TestUnd_SqlNull_Scan(t *testing.T) {
	var n und.SqlNull[int64]
	if err := n.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if !n.IsNull() {
		t.Errorf("Scan(nil): state %d, want null", n.State())
	}

	n = und.SqlNull[int64]{}
	if err := n.Scan(int64(42)); err != nil {
		t.Fatalf("Scan(42): %v", err)
	}
	// Value() on SqlNull is the driver.Valuer; reach the inner value via Und.
	if !n.IsDefined() || n.Und.Value() != 42 {
		t.Errorf("Scan(42): defined=%t value=%d", n.IsDefined(), n.Und.Value())
	}
}
