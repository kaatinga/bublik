package bublyk

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/jackc/pgio"
	"github.com/jackc/pgtype"
)

// Value implements the database/sql/driver Valuer interface.
func (thisDate Date) Value() (driver.Value, error) {
	switch thisDate {
	case noDate:
		return thisDate.String(), nil
	default:
		return nil, nil
	}
}

func (thisDate *Date) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*thisDate = noDate
		return nil
	}

	sbuf := string(src)
	switch sbuf {
	case "infinity":
		*thisDate = maximumDate
	case "-infinity":
		*thisDate = minimumDate
	default:
		var err error
		*thisDate, err = Parse(sbuf)
		if err != nil {
			return err
		}
	}

	return nil
}

const (
	negativeInfinityDayOffset = -2147483648
	infinityDayOffset         = 2147483647
)

func (thisDate *Date) DecodeBinary(_ *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		*thisDate = noDate
		return nil
	}

	if len(src) != 4 {
		return fmt.Errorf("invalid length for date: %v", len(src))
	}

	dayOffset := int32(binary.BigEndian.Uint32(src))

	switch dayOffset {
	case infinityDayOffset:
		*thisDate = maximumDate
	case negativeInfinityDayOffset:
		*thisDate = minimumDate
	default:
		t := time.Date(2000, 1, int(1+dayOffset), 0, 0, 0, 0, time.UTC)
		*thisDate = NewDateFromTime(&t)
	}

	return nil
}

func (thisDate Date) EncodeText(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	if !thisDate.IsSet() {
		return nil, nil
	}

	return append(buf, format[[]byte](thisDate)...), nil
}

func (thisDate Date) EncodeBinary(_ *pgtype.ConnInfo, buf []byte) ([]byte, error) {
	if !thisDate.IsSet() {
		return nil, nil
	}

	var daysSinceDateEpoch int32
	tUnix := thisDate.Time().Unix()
	dateEpoch := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Unix()

	secSinceDateEpoch := tUnix - dateEpoch
	daysSinceDateEpoch = int32(secSinceDateEpoch / 86400)

	return pgio.AppendInt32(buf, daysSinceDateEpoch), nil
}
