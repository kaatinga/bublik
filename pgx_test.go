package bublyk

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net"
	"testing"
	"time"
)

func TestDate_WithBD(t *testing.T) {
	pool, err := NewDBConnection()
	if err != nil {
		t.Error(err)
	}
	defer pool.Close()

	ctx := context.Background()
	_, err = pool.Exec(ctx, `
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

func NewDBConnection() (*pgxpool.Pool, error) {
	var ConnPool *pgxpool.Pool
	initCtx := context.Background()

	dbConnectionConfig, err := pgxpool.ParseConfig("postgres://*:*@localhost:5432/postgres")
	if err != nil {
		return nil, err
	}

	dbConnectionConfig.AfterRelease = func(*pgx.Conn) bool {
		log.Println("db connection released")
		return true
	}

	dbConnectionConfig.MaxConns = 8
	dbConnectionConfig.MinConns = dbConnectionConfig.MaxConns >> 1

	dbConnectionConfig.HealthCheckPeriod = time.Minute

	dbConnectionConfig.MaxConnLifetime = time.Hour

	dbConnectionConfig.MaxConnIdleTime = time.Hour

	dbConnectionConfig.ConnConfig.ConnectTimeout = 1 * time.Second

	dbConnectionConfig.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: dbConnectionConfig.HealthCheckPeriod,
		Timeout:   dbConnectionConfig.ConnConfig.ConnectTimeout,
	}).DialContext

	var conn *pgx.Conn
	ctx := context.Background()
	conn, err = pgx.ConnectConfig(ctx, dbConnectionConfig.ConnConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	var maxConnInDB int32
	err = conn.
		QueryRow(initCtx, "SELECT CAST(setting AS integer) FROM pg_settings WHERE name='max_connections';").
		Scan(&maxConnInDB)
	if err != nil {
		return nil, err
	}
	if dbConnectionConfig.MaxConns > maxConnInDB {
		return nil, fmt.Errorf("incorrect connection pool size %d is set in db", maxConnInDB)
	}

	ConnPool, err = pgxpool.ConnectConfig(initCtx, dbConnectionConfig)
	if err != nil {
		return nil, err
	}

	return ConnPool, ConnPool.Ping(initCtx)
}
