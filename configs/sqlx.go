package configs

import (
	"fmt"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	// Import the PostgreSQL driver for sqlx
	_ "github.com/lib/pq"
)

// GetReaderSqlx establishes a connection to the PostgreSQL database for reading
func GetReaderSqlx() *sqlx.DB {
	connStr := buildConnectionString()
	reader := sqlx.MustConnect("postgres", connStr)

	return reader
}

// GetWriterSqlx establishes a connection to the PostgreSQL database for writing
func GetWriterSqlx() *sqlx.DB {
	connStr := buildConnectionString()
	writer := sqlx.MustConnect("postgres", connStr)

	return writer
}

// buildConnectionString constructs the PostgreSQL connection string
func buildConnectionString() string {
	user := os.Getenv("DB_USER")
	password := url.QueryEscape(os.Getenv("DB_PASSWORD"))
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, dbname, sslmode)
}
