package bublyk

import (
	"fmt"
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	tests := []struct {
		year  int
		month int
		day   int
	}{
		{2021, 12, 31},
		{2001, 11, 30},
		{2021, 12, 15},
		{2031, 1, 1},
		{2127, 12, 31},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d-%d-%d", tt.year, tt.month, tt.day), func(t *testing.T) {
			newTime := time.Date(tt.year, time.Month(tt.month), tt.day, 0, 0, 0, 0, time.UTC)
			date := NewDateFromTime(&newTime)
			if uint16(tt.year) != date.Year() {
				t.Errorf("Year is incorrect.\nhave %v\nwant %v", date.Year(), tt.year)
			}

			if byte(tt.month) != date.Month() {
				t.Errorf("Month is incorrect.\nhave %v\nwant %v", date.Month(), tt.month)
			}

			if byte(tt.day) != date.Day() {
				t.Errorf("Day is incorrect.\nhave %v\nwant %v", date.Day(), tt.day)
			}

			date2 := NewDate(uint16(tt.year), byte(tt.month), byte(tt.day))
			if uint16(tt.year) != date2.Year() {
				t.Errorf("Year is incorrect.\nhave %v\nwant %v", date.Year(), tt.year)
			}

			if byte(tt.month) != date2.Month() {
				t.Errorf("Month is incorrect.\nhave %v\nwant %v", date.Month(), tt.month)
			}

			if byte(tt.day) != date2.Day() {
				t.Errorf("Day is incorrect.\nhave %v\nwant %v", date.Day(), tt.day)
			}

			t.Logf("%16b\n", date)
		})
	}

}

func TestDate_Format(t *testing.T) {
	tests := []struct {
		this   Date
		layout string
		want   string
	}{
		{maximumDate, PostgreSQLFormat, "2127-12-31"},
		{maximumDate, time.RFC822, "31 Dec 27 00:00 UTC"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.this.Format(tt.layout); got != tt.want {
				t.Errorf("Format() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextDay(t *testing.T) {
	tests := []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2022, 1, 1)},
		{NewDate(2021, 12, 1), NewDate(2021, 12, 2)},
		{NewDate(2021, 2, 28), NewDate(2021, 3, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.NextDay(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_NextWeek(t *testing.T) {
	tests := []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2022, 1, 7)},
		{NewDate(2021, 12, 1), NewDate(2021, 12, 8)},
		{NewDate(2021, 2, 28), NewDate(2021, 3, 7)},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.NextWeek(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousWeek(t *testing.T) {
	tests := []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2021, 12, 24)},
		{NewDate(2022, 1, 3), NewDate(2021, 12, 27)},
		{NewDate(2021, 12, 1), NewDate(2021, 11, 24)},
		{NewDate(2021, 3, 2), NewDate(2021, 2, 23)},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.PreviousWeek(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDate_PreviousMonth(t *testing.T) {
	tests := []struct {
		this Date
		want Date
	}{
		{NewDate(2021, 12, 31), NewDate(2021, 12, 1)},
		{NewDate(2024, 2, 29), NewDate(2024, 1, 29)},
		{NewDate(2021, 12, 1), NewDate(2021, 11, 1)},
		{NewDate(2021, 3, 2), NewDate(2021, 2, 2)},
		{NewDate(2024, 1, 1), NewDate(2023, 12, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.this.String(), func(t *testing.T) {
			if got := tt.this.PreviousMonth(); got != tt.want {
				t.Errorf("NextDay() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParse(t *testing.T) {

	tests := []struct {
		formattedDate string
		want          Date
		wantErr       bool
	}{
		{"2021-12-01", NewDate(2021, 12, 1), false},
	}
	for _, tt := range tests {
		t.Run(tt.formattedDate, func(t *testing.T) {
			got, err := Parse(tt.formattedDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
