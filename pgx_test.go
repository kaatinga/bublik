//go:build localtests

package bublyk

import (
	"testing"
	"time"

	"github.com/kaatinga/bochka"
)

func TestDate_WithBD(t *testing.T) {
	helper := bochka.NewPostgreTestHelper(t, bochka.WithTimeout(10*time.Second))
	helper.Run("14.5")

	t.Cleanup(func() {
		_, err := helper.Exec(helper.Context, `DROP TABLE IF EXISTS tmp1`)
		if err != nil {
			t.Error("Test table deletion failed:", err)
		}
		helper.Close()
	})

	_, err := helper.Exec(helper.Context, `
CREATE TABLE IF NOT EXISTS tmp1 (
	testDate date
); `)
	if err != nil {
		t.Error("Test table creation failed:", err)
	}

	t.Run("test date 1", func(t *testing.T) {
		inputDate := Now()
		var returnedDate Date
		err = helper.QueryRow(helper.Context, `
INSERT INTO tmp1(testdate) VALUES($1) RETURNING testdate;
`, inputDate).Scan(&returnedDate)
		if err != nil {
			t.Error("INSERT query failed:", err)
		}
		if returnedDate == inputDate {
			t.Log("success!")
			t.Logf("returned inputDate: '%s'", returnedDate)
			t.Logf("input inputDate: '%s'", inputDate)
		} else {
			t.Errorf("dates are not equal, have '%s', want '%s'", returnedDate, inputDate)
		}
	})

	t.Run("test date 2", func(t *testing.T) {
		inputDate := NewDate(2022, 12, 31)
		var returnedDate Date
		err = helper.QueryRow(helper.Context, `
INSERT INTO tmp1(testdate) VALUES($1) RETURNING testdate;
`, inputDate).Scan(&returnedDate)
		if err != nil {
			t.Error("INSERT query failed:", err)
		}
		if returnedDate == inputDate {
			t.Log("success!")
			t.Logf("returned inputDate: '%s'", returnedDate)
			t.Logf("input inputDate: '%s'", inputDate)
		} else {
			t.Errorf("dates are not equal, have '%s', want '%s'", returnedDate, inputDate)
		}
	})

	t.Run("test date 3", func(t *testing.T) {
		inputDate := NewDate(2000, 1, 1)
		var returnedDate Date
		err = helper.QueryRow(helper.Context, `
INSERT INTO tmp1(testdate) VALUES($1) RETURNING testdate;
`, inputDate).Scan(&returnedDate)
		if err != nil {
			t.Error("INSERT query failed:", err)
		}
		if returnedDate == inputDate {
			t.Log("success!")
			t.Logf("returned inputDate: '%s'", returnedDate)
			t.Logf("input inputDate: '%s'", inputDate)
		} else {
			t.Errorf("dates are not equal, have '%s', want '%s'", returnedDate, inputDate)
		}
	})

	t.Run("test null date", func(t *testing.T) {
		var inputDate Date = 0
		var isNull bool
		err = helper.QueryRow(helper.Context, `
INSERT INTO tmp1(testdate) VALUES($1) RETURNING testdate IS NULL;
`, inputDate).Scan(&isNull)
		if err != nil {
			t.Error("INSERT query failed:", err)
		}
		if !isNull {
			t.Error("the value must be null but it is not")
		}
	})
}
