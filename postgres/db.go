package postgres

import (
	"database/sql"
	"regexp"

	_ "github.com/lib/pq"
)

type DBSize struct {
	Disk         string
	Uncompressed string
	Messages     int
	Channels     int
}

var db *sql.DB
var spaces = regexp.MustCompile(`\s+`)

// Connect to the postgres database using string provided
func Connect(dbInfo string) error {
	conn, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	db = conn
	return nil
}

// Close the sql connection
func Close() {
	if db != nil {
		db.Close()
	}
}

// GetDB will return the current sql.DB client
func GetDB() *sql.DB {
	return db
}

// CreateTables will create required tables in the database if they do not already exist
func CreateTables() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS channels (
			id SERIAL NOT NULL,
			channel VARCHAR(255) NOT NULL DEFAULT '',
			status BOOLEAN NOT NULL DEFAULT FALSE,
			added_by VARCHAR(255) NOT NULL DEFAULT 'system',
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(channel)
		);

		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL NOT NULL,
			channel VARCHAR(255) NOT NULL DEFAULT '',
			username VARCHAR(255) NOT NULL DEFAULT '',
			message VARCHAR(500) NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(channel, username, message, created_at)
		);
	`)
	return err
}

// Ping the database
func Ping() error {
	return db.Ping()
}

// GetSize gets various information about the database.
// NOTE: DBSize.Messages and DBSize.Channels are not always 100% accurate.
func GetSize() (DBSize, error) {
	size := DBSize{}
	if err := db.QueryRow(`
		SELECT
			pg_size_pretty(pg_total_relation_size(relid)) AS uncompressed,
			pg_size_pretty(pg_total_relation_size(relid) - pg_relation_size(relid)) as disk
		FROM pg_catalog.pg_statio_user_tables ORDER BY pg_total_relation_size(relid) DESC
	`).Scan(&size.Uncompressed, &size.Disk); err != nil {
		return size, err
	}
	size.Disk = spaces.ReplaceAllString(size.Disk, "")
	size.Uncompressed = spaces.ReplaceAllString(size.Uncompressed, "")
	if err := db.QueryRow(`SELECT reltuples::BIGINT AS estimate FROM pg_class WHERE relname='messages'`).Scan(&size.Messages); err != nil {
		return size, err
	}
	if err := db.QueryRow(`SELECT reltuples::BIGINT AS estimate FROM pg_class WHERE relname='channels'`).Scan(&size.Channels); err != nil {
		return size, err
	}
	return size, nil
}
