[![GitHub release](https://img.shields.io/github/release/kaatinga/bublyk.svg)](https://github.com/kaatinga/bublyk/releases)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/bublyk/blob/main/LICENSE)
[![codecov](https://codecov.io/gh/kaatinga/bublyk/branch/main/graph/badge.svg?token=Q34SE0KN9E)](https://codecov.io/gh/kaatinga/bublyk)
[![lint workflow](https://github.com/kaatinga/bublyk/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/bublyk/actions?query=workflow%3Alinter)
[![help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/bublyk/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

# bublyk

The package implements type `Date` that is targeted to the cases when we do not need to work with time yet only with date
in UTC location. `Date` type has some benefits in compression to `time.Time` type as `Date` consumes much less memory, does not
require boilerplate code to work with date, also it is comparison-enabled using operators `>`, `<`, etc.

As bonus, it supports `pgx` package natively what means you can directly scan into `Date` type as well as to use
`Date` as argument in queries.

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