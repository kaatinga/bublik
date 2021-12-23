package bublyk

import (
	"github.com/kaatinga/assets"
	"time"
)

const (
	yearMask  = 0b0000000001111111
	monthMask = 0b0000011110000000
	dayMask   = 0b1111100000000000

	maximumDate Date = 0b1111111001111111 // 2127-12-31

	PostgreSQLFormat = "2006-01-02"
)

// Date represents a calendar date starting 2000 year and finishing the year 2127.
type Date uint16

func (this Date) Year() uint16 {
	return uint16(this&yearMask) + 2000
}

func (this Date) Month() byte {
	return byte((this & monthMask) >> 7)
}

func (this Date) Day() byte {
	return byte((this & dayMask) >> 11)
}

// String returns string in the default PostgreSQL date format, YYYY-MM-DD.
func (this Date) String() string {
	return assets.Uint162String(this.Year()) + "-" + assets.Byte2String(this.Month()) + "-" + assets.Byte2String(this.Day())
}

func (this Date) Format(layout string) string {
	switch layout {
	case PostgreSQLFormat:
		return assets.Uint162String(this.Year()) + "-" + assets.Byte2String(this.Month()) + "-" + assets.Byte2String(this.Day())
	default:
		return time.Date(int(this.Year()), time.Month(this.Month()), int(this.Day()), 0, 0, 0, 0, time.UTC).Format(layout)
	}
}

func NewDate(year uint16, month, day byte) Date {
	if year < 2000 {
		return 0
	}
	if year > 2127 {
		return maximumDate
	}
	return Date(year-2000) + (Date(month) << 7) + (Date(day) << 11)
}

func NewDateFromTime(t *time.Time) Date {
	year := t.Year()
	if year < 2000 {
		return 0
	}
	if year > 2127 {
		return maximumDate
	}
	return Date(year-2000) + (Date(t.Month()) << 7) + (Date(t.Day()) << 11)
}
