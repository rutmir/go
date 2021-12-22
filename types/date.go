package types

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"time"

	pgtypes "github.com/go-pg/pg/v10/types"
)

type Date time.Time

func (dt *Date) String() string {
	return time.Time(*dt).Format("2006-01-02")
}

func (dt *Date) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var s string
	if err := d.DecodeElement(&s, &start); err != nil {
		return err
	}

	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}

	*dt = Date(t)

	return nil
}

func (dt *Date) UnmarshalParam(param string) error {
	t, err := time.Parse("2006-01-02", param)
	if err != nil {
		return err
	}

	*dt = Date(t)

	return nil
}

// JSON.
func (dt *Date) UnmarshalJSON(input []byte) error {
	t, err := time.Parse("2006-01-02", strings.Trim(string(input), `"`))
	if err != nil {
		return err
	}

	*dt = Date(t)

	return nil
}

func (dt *Date) MarshalJSON() ([]byte, error) {
	if dt == nil || time.Time(*dt).IsZero() {
		return nil, nil
	}

	return json.Marshal(dt.String())
}

// GOB.
func (dt *Date) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		*dt = Date(time.Time{})

		return nil
	}

	return dt.UnmarshalJSON(data)
}

func (dt *Date) MarshalBinary() ([]byte, error) {
	if dt == nil || time.Time(*dt).IsZero() {
		return nil, nil
	}

	return []byte(dt.String()), nil
}

// GO-PG.
var _ pgtypes.ValueAppender = (*Date)(nil)

func (dt Date) AppendValue(b []byte, flags int) ([]byte, error) {
	if flags == 1 {
		b = append(b, '\'')
	}
	b = time.Time(dt).AppendFormat(b, "2006-01-02")
	if flags == 1 {
		b = append(b, '\'')
	}
	return b, nil
}

var _ pgtypes.ValueScanner = (*Date)(nil)

func (dt *Date) ScanValue(rd pgtypes.Reader, n int) error {
	if n <= 0 {
		*dt = Date(time.Time{})
		return nil
	}

	tmp, err := rd.ReadFullTemp()
	if err != nil {
		return err
	}

	tm, err := time.ParseInLocation("2006-01-02", string(tmp), time.UTC)
	if err != nil {
		return err
	}

	*dt = Date(tm)

	return nil
}
