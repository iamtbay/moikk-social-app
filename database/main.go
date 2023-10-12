package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn
var dbTimeout = time.Second * 10
var xctx = context.Background()
var dataPerPage = 10
var notFoundMsg = "Are you sure what are you looking for? Check parameters"

func InitDB() {
	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("database activated!")
}
