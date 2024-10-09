package mem

import (
	"strconv"
	"time"
)

const (
	Layout      = "01/02 03:04:05PM '06 -0700" // The reference time, in numerical order.
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	Kitchen     = "3:04PM"

	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
	DateTime   = "2006-01-02 15:04:05"
	DateOnly   = "2006-01-02"
	TimeOnly   = "15:04:05"
)

type Data struct {
	s  string
	i  int64
	f  float64
	b  bool
	t  time.Time
	is string
	o  *Data
	a  []Data
}

func FromString(s string) *Data {
	data := Data{s: s}
	data.is = "string"
	b, err := strconv.ParseBool(s)
	if err == nil {
		data.is = "bool"
		data.b = b
		return &data
	}

	i, err := strconv.ParseInt(s, 10, 64)

	if err == nil {
		data.i = i
		data.is = "int"
		data.f = float64(i)
		return &data
	}

	f, err := strconv.ParseFloat(s, 64)

	if err == nil {
		data.f = f
		data.is = "float"
		data.i = int64(f)
		return &data
	}

	t, err := time.Parse(RFC3339, s)

	if err == nil {
		data.t = t
		data.is = "datetime"
		return &data
	}

	t, err = time.Parse(DateTime, s)

	if err == nil {
		data.t = t
		data.is = "datetime"
		return &data
	}

	t, err = time.Parse(DateOnly, s)

	if err == nil {
		data.t = t
		data.is = "date"
		return &data
	}

	return &data
}

func (d Data) Str() string {
	return d.s
}

func (d Data) Is() string {
	return d.is
}

func (d Data) Val() any {
	if d.is == "object" {
		return *d.o
	}
	if d.is == "array" {
		return d.a
	}

	return d.s
}
