[![Tests](https://github.com/kaatinga/luna/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kaatinga/luna/actions/workflows/test.yml)
[![GitHub release](https://img.shields.io/github/release/kaatinga/bublyk.svg)](https://github.com/kaatinga/bublyk/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/bublyk/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/bublyk/branch/main/graph/badge.svg?token=Q34SE0KN9E)](https://codecov.io/gh/kaatinga/bublyk)
[![lint workflow](https://github.com/kaatinga/bublyk/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/bublyk/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/bublyk/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# bublyk

The package introduces the Date type, specifically designed for instances where only the date in UTC location is required, without the need for time details. In comparison to the time.Time type, Date offers several advantages:
- It consumes significantly less memory.
- It eliminates the need for boilerplate code when working with dates.
- It allows for straightforward comparisons using operators such as >, <, and others.

Additionally, it natively supports the pgx package. This means you can directly scan into the Date type and use Date as an argument in queries.

## Usage

```go
package main

import (
    "context"
    "fmt"
    "github.com/kaatinga/bublyk"
    "github.com/jackc/pgx/v4/pgxpool"
)

func main() {
    ctx := context.Background()
    pool, err := pgxpool.Connect(ctx, "postgres://postgres:postgres@localhost:5432/postgres")
    if err != nil {
        panic(err)
    }
    defer pool.Close()

    _, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS test (test_date DATE)")
    if err != nil {
        panic(err)
    }

    inputDate := bublyk.Now()
    var returnedDate bublyk.Date
    err = pool.QueryRow(ctx, "INSERT INTO test(test_date) VALUES($1) RETURNING test_date", inputDate).Scan(&returnedDate)
    if err != nil {
        panic(err)
    }

    fmt.Println(inputDate, returnedDate)
}
```

Will be happy to everyone who want to participate in the work on the `bublyk` package.
