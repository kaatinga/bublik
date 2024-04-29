//go:build localtests

package bublyk

import (
	"testing"
	"time"

	"github.com/kaatinga/bochka"
	"github.com/stretchr/testify/suite"
)

func TestBublykSuite(t *testing.T) {
	suite.Run(t, new(TestsSuite))
}

type TestsSuite struct {
	suite.Suite

	db *bochka.PostgreTestHelper
}

func (suite *TestsSuite) SetupSuite() {
	suite.db = bochka.NewPostgreTestHelper(suite.T(), bochka.WithTimeout(10*time.Second))
	suite.db.Run("16")
}

func (suite *TestsSuite) TearDownSuite() {
	suite.db.Close()
}

func (suite *TestsSuite) TestDate_WithBD() {
	t := suite.T()
	t.Cleanup(func() {
		_, err := suite.db.Pool.Exec(suite.db.Context, `DROP TABLE IF EXISTS tmp1`)
		if err != nil {
			t.Error("Test table deletion failed:", err)
		}
	})

	_, err := suite.db.Pool.Exec(suite.db.Context, `
CREATE TABLE IF NOT EXISTS tmp1 (
	testDate date
); `)
	if err != nil {
		t.Error("Test table creation failed:", err)
	}

	t.Run("test date 1", func(t *testing.T) {
		inputDate := Now()
		var returnedDate Date
		err = suite.db.Pool.QueryRow(suite.db.Context, `
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
		err = suite.db.Pool.QueryRow(suite.db.Context, `
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
		err = suite.db.Pool.QueryRow(suite.db.Context, `
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
		err = suite.db.Pool.QueryRow(suite.db.Context, `
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
