package bublyk

import (
	"fmt"
	"github.com/kaatinga/assets"
	"reflect"
	"time"
)

const (
	yearMask  = 0b0000000001111111
	monthMask = 0b0000011110000000
	dayMask   = 0b1111100000000000

	maximumDate Date = 0b1111111001111111 // 2127-12-31

	PostgreSQLFormat = "2006-01-02"
)

func Now() Date {
	now := time.Now().UTC()
	return NewDateFromTime(&now)
}

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

func (this Date) YearInt() int {
	return int(this&yearMask) + 2000
}

func (this Date) MonthMonth() time.Month {
	return time.Month((this & monthMask) >> 7)
}

func (this Date) DayInt() int {
	return int((this & dayMask) >> 11)
}

func (this Date) IsSet() bool {
	return this != 0
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
	return assets.Uint162String(this.Year()) + "-" + assets.Byte2String(this.Month()) + "-" + assets.Byte2String(this.Day())
}

// DMYWithDots returns date as string in the DD.MM.YYYY format.
func (this Date) DMYWithDots() string {
	return assets.Byte2String(this.Day()) + "." + assets.Byte2String(this.Month()) + "." + assets.Uint162String(this.Year())
}

func (this *Date) Set(src interface{}) error {
	if src == nil {
		*this = 0
		return nil
	}

	switch value := src.(type) {
	case time.Time:
		*this = NewDateFromTime(&value)
	case string:
		t, err := time.ParseInLocation("2006-01-02", value, time.UTC)
		if err != nil {
			return err
		}
		*this = NewDateFromTime(&t)
	case *time.Time:
		if value == nil {
			*this = 0
		} else {
			return this.Set(*value)
		}
	case *string:
		if value == nil {
			*this = 0
		} else {
			return this.Set(*value)
		}
	default:
		if originalSrc, ok := underlyingTimeType(src); ok {
			return this.Set(originalSrc)
		}
		return fmt.Errorf("cannot convert %v to Date", value)
	}

	return nil
}

func underlyingTimeType(val interface{}) (interface{}, bool) {
	refVal := reflect.ValueOf(val)

	switch refVal.Kind() {
	case reflect.Ptr:
		if refVal.IsNil() {
			return nil, false
		}
		convVal := refVal.Elem().Interface()
		return convVal, true
	}

	timeType := reflect.TypeOf(time.Time{})
	if refVal.Type().ConvertibleTo(timeType) {
		return refVal.Convert(timeType).Interface(), true
	}

	return nil, false
}

func (this Date) Format(layout string) string {
	switch layout {
	case PostgreSQLFormat:
		return this.String()
	default:
		return makeTime(this.Year(), this.Month(), this.Day()).Format(layout)
	}
}

func (this Date) NextDay() Date {
	if this.Day() > 27 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, 1)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this>>11+1)<<11
}

func (this Date) PreviousDay() Date {
	if this.Day() == 1 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, -1)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this>>11-1)<<11
}

func (this Date) NextWeek() Date {
	if this.Day() > 21 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, 7)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this>>11+7)<<11
}

func (this Date) PreviousWeek() Date {
	if this.Day() < 8 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 0, -7)
		return NewDateFromTime(&timeDate)
	}
	return this&^dayMask | (this>>11-7)<<11
}

func (this Date) NextMonth() Date {
	if this>>11 > 28 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, 1, 0)
		return NewDateFromTime(&timeDate)
	}

	if this&^dayMask>>7 == 12 {
		if this&^monthMask&^dayMask == 127 { // we reached the maximum
			return maximumDate
		}
		return this&^monthMask&^yearMask | (1 << 7) | this&^monthMask&^dayMask + 1 // January
	}
	return this&^monthMask | (this&^dayMask>>7+1)<<7
}

func (this Date) PreviousMonth() Date {
	if this>>11 > 28 || this&^dayMask>>7 == 1 {
		timeDate := makeTime(this.Year(), this.Month(), this.Day()).AddDate(0, -1, 0)
		return NewDateFromTime(&timeDate)
	}
	return this&^monthMask | (this&^dayMask>>7-1)<<7
}

func NewDate(year uint16, month, day byte) Date {
	if year < 2000 {
		return 0
	}
	if year > 2127 {
		return maximumDate
	}
	if day > 28 || month > 12 || day == 0 || month == 0 {
		yearInt, monthMonth, dayInt := makeTime(year, month, day).Date()
		year, month, day = uint16(yearInt), byte(monthMonth), byte(dayInt)
	}
	return Date(year-2000) + (Date(month) << 7) + (Date(day) << 11)
}

func makeTime(year uint16, month, day byte) *time.Time {
	newTime := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC)
	return &newTime
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
