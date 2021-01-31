package postgres

import (
	"database/sql"
	"regexp"

	_ "github.com/lib/pq"
)

var db *sql.DB
var spaces = regexp.MustCompile(`\s+`)

func Connect(dbInfo string) error {
	conn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	db = conn
	return Ping()
}

func Close() error {
	return db.Close()
}

func Ping() error {
	return db.Ping()
}

func GetDatabase() *sql.DB {
	return db
}

func GetTotalLinesStored() (total int64, err error) {
	stmt, err := db.Prepare(`SELECT SUM(total) as total FROM count`)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow().Scan(&total)
	return
}
