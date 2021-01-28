package dif

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var connStr = "host=postgres user=postgres password=postgres dbname=postgres sslmode=disable"
var DB, DBError = sqlx.Connect("postgres", connStr)
