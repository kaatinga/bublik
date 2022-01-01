package bublyk

import (
	"github.com/kaatinga/assets"
	"time"
)

const (
	yearMask  = 0b1111111000000000
	monthMask = 0b0000000111100000
	dayMask   = 0b0000000000011111

	maximumDate Date = 0b1111111110011111 // 2127-12-31
	minimumDate Date = 0b0000000000100001 // 2000-01-01

	PostgreSQLFormat = "2006-01-02"
)

func Now() Date {
	now := time.Now().UTC()
	return NewDateFromTime(&now)
}

func CurrentMonth() Date {
	now := Now()
	return NewDate(now.Year(), now.Month(), 1)
}

// Date represents a calendar date starting 2000 year and finishing the year 2127.
type Date uint16

func (this Date) Year() uint16 {
	return uint16(this>>9) + 2000
}

func (this Date) Month() byte {
	return byte((this & monthMask) >> 5)
}

func (this Date) Day() byte {
	return byte(this & dayMask)
}

func (this Date) IsSet() bool {
	return this != 0
}

func (this Date) IsFuture() bool {
	now := Now()
	return this > now
}

func (this Date) MonthAfter(date Date) bool {
	if this.Year() == date.Year() {
		return this.Month() > date.Month()
	}
	return this.Year() > date.Year()
}

func (this Date) MonthBefore(date Date) bool {
	if this.Year() == date.Year() {
		return this.Month() < date.Month()
	}
	return this.Year() < date.Year()
}

// String returns date as string in the default PostgreSQL date format, YYYY-MM-DD.
func (this Date) String() string {
	var month = assets.Byte2String(this.Month())
	var day = assets.Byte2String(this.Day())

	// right format for month < 10
	if len(month) == 1 {
		month = "0" + month
	}

	// right format for day < 10
	if len(day) == 1 {
		day = "0" + day
	}
	return assets.Uint162String(this.Year()) + "-" + month + "-" + day
}

// DMYWithDots returns date as string in the DD.MM.YYYY format.
func (this Date) DMYWithDots() string {

	var month = assets.Byte2String(this.Month())
	var day = assets.Byte2String(this.Day())

	// right format for month < 10
	if len(month) == 1 {
		month = "0" + month
	}

	// right format for day < 10
	if len(day) == 1 {
		day = "0" + day
	}

	return day + "." + month + "." + assets.Uint162String(this.Year())
}

func (this Date) Format(layout string) string {
	switch layout {
	case PostgreSQLFormat:
		return this.String()
	default:
		return makeTime(this.Year(), this.Month(), this.Day()).Format(layout)
	}
}

func Parse(formattedDate string) (Date, error) {
	if len([]rune(formattedDate)) != len([]rune(PostgreSQLFormat)) {
		return 0, ErrUnrecognizedFormat
	}

	year, err := assets.String2Uint16(formattedDate[0:4])
	if err != nil {
		return 0, err
	}

	var month byte
	month, err = assets.String2Byte(formattedDate[5:7])
	if err != nil {
		return 0, err
	}

	var day byte
	day, err = assets.String2Byte(formattedDate[8:10])
	if err != nil {
		return 0, err
	}

	return NewDate(year, month, day), nil
}

func (this Date) Time() *time.Time {
	return makeTime(this.Year(), this.Month(), this.Day())
}

func (this Date) NextDay() Date {
	if this.Day() > 27 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, 1)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this&dayMask + 1)
}

func (this Date) PreviousDay() Date {
	if this.Day() == 1 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, -1)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this&dayMask - 1)
}

func (this Date) NextWeek() Date {
	if this.Day() > 21 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, 7)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this&dayMask + 7)
}

func (this Date) PreviousWeek() Date {
	if this.Day() < 8 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, -7)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this&dayMask - 7)
}

func (this Date) NextMonth() Date {
	if this&dayMask > 28 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 1, 0)
		return NewDateFromTime(&timeDate)
	}

	if (this&monthMask)>>5 == 12 {
		if this&yearMask == 127 { // we reached the maximum
			return maximumDate
		}
		return this&^monthMask&^yearMask | ((1 << 5) | (this&yearMask>>9+1)<<9) // January
	}
	return this&^monthMask | ((((this & monthMask) >> 5) + 1) << 5)
}

func (this Date) PreviousMonth() Date {
	if this.Day() > 28 || this.Month() == 1 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, -1, 0)
		return NewDateFromTime(&timeDate)
	}
	return composeDate(this.Year(), this.Month()-1, this.Day())
}

func NewDate(year uint16, month, day byte) Date {
	if year < 2000 {
		return minimumDate
	}
	if year > 2127 {
		return maximumDate
	}
	if day > 28 || month > 12 || day == 0 || month == 0 {
		yearInt, monthMonth, dayInt := makeTime(year, month, day).Date()
		year, month, day = uint16(yearInt), byte(monthMonth), byte(dayInt)
	}
	return composeDate(year, month, day)
}

func makeTime(year uint16, month, day byte) *time.Time {
	newTime := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	return &newTime
}

func NewDateFromTime(t *time.Time) Date {
	year := t.Year()
	if year < 2000 {
		return minimumDate
	}
	if year > 2127 {
		return maximumDate
	}
	return composeDate(uint16(year), byte(t.Month()), byte(t.Day()))
}

func composeDate(year uint16, month, day byte) Date {
	return (Date(year-2000) << 9) + (Date(month) << 5) + Date(day)
}
