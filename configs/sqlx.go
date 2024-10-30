package configs

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetReaderSqlx() *sqlx.DB {
	reader := sqlx.MustConnect("mysql", os.Getenv("DB_CONNECTION"))

	return reader
}

func GetWriterSqlx() *sqlx.DB {
	writer := sqlx.MustConnect("mysql", os.Getenv("DB_CONNECTION"))

	return writer
}
