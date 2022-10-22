package bublyk

import (
	"context"
	tests "github.com/kaatinga/postgreSQLtesthelper"
	"testing"
)

func TestDate_WithBD(t *testing.T) {
	ctx := context.Background()
	dbContainer, pool := tests.SetupPostgreDatabase("kaatinga", "12345", t)
	defer pool.Close()
	defer dbContainer.Terminate(ctx)

	_, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS tmp1 (
	testDate date
); `)
	if err != nil {
		t.Error("Test table creation failed:", err)
	}

	t.Run("test date 1", func(t *testing.T) {
		inputDate := Now()
		var returnedDate Date
		err = pool.QueryRow(ctx, `
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

	t.Run("test date 1", func(t *testing.T) {
		inputDate := NewDate(2022, 12, 31)
		var returnedDate Date
		err = pool.QueryRow(ctx, `
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
		err = pool.QueryRow(ctx, `
INSERT INTO tmp1(testdate) VALUES($1) RETURNING testdate IS NULL;
`, inputDate).Scan(&isNull)
		if err != nil {
			t.Error("INSERT query failed:", err)
		}
		if !isNull {
			t.Error("the value miust be null but it is not")
		}
	})

	_, err = pool.Exec(ctx, `DROP TABLE tmp1`)
	if err != nil {
		t.Error("Test table deletion failed:", err)
	}

}
