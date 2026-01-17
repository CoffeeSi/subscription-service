package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Connect *pgx.Conn
}

func NewDatabase(ctx context.Context, connectURL string) *Database {
	connect, err := pgx.Connect(ctx, connectURL)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	return &Database{
		Connect: connect,
	}
}
