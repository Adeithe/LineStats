package postgres

import (
	"fmt"
	"time"
)

type ChatQuote struct {
	Channel string
	Sender  string
	Message string
	SentAt  time.Time
}

func SaveQuote(quote ChatQuote) error {
	stmt, err := db.Prepare(`INSERT INTO messages (channel, username, message, created_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(quote.Channel, quote.Sender, quote.Message, quote.SentAt.UTC())
	return err
}

func MessageCountByUser(channel string, username string) (int, error) {
	var count int
	stmt, err := db.Prepare(`SELECT COUNT(*) FROM messages WHERE channel=$1 AND username=$2`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(channel, username).Scan(&count)
	return count, err
}

func UniqueMessageCountByUser(channel string, username string) (int, error) {
	var count int
	stmt, err := db.Prepare(`SELECT COUNT(DISTINCT message) FROM messages WHERE channel=$1 AND username=$2`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(channel, username).Scan(&count)
	return count, err
}

func GetMostLinesByUser(channel string, username string) (date string, count int, err error) {
	stmt, err := db.Prepare(`SELECT TO_CHAR(created_at, 'Month YYYY') as mon, COUNT(*) FROM messages WHERE channel=$1 AND username=$2 GROUP BY mon`)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(channel, username)
	if err != nil {
		return
	}
	for rows.Next() {
		var d string
		var c int
		if err := rows.Scan(&d, &c); err != nil {
			break
		}
		if c > count {
			date = spaces.ReplaceAllString(d, " ")
			count = c
		}
	}
	rows.Close()
	return
}

func RandomQuoteByUser(channel string, username string) (quote ChatQuote, err error) {
	quote = ChatQuote{Channel: channel}
	// SELECT username, message, created_at FROM (
	//		SELECT DISTINCT ON(message) * FROM messages TABLESAMPLE SYSTEM(1) WHERE channel=$1 AND username=$2
	//	) AS s ORDER BY RANDOM() LIMIT 1;
	stmt, err := db.Prepare(`
		SELECT DISTINCT username, message, created_at FROM messages WHERE channel=$1 AND username=$2
			OFFSET floor(random()* (
				SELECT DISTINCT COUNT(*) FROM messages WHERE channel=$1 AND username=$2
			)) LIMIT 1;
	`)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(channel, username).Scan(&quote.Sender, &quote.Message, &quote.SentAt)
	return
}

func ScanMessagesByUser(channel string, username string, query string) (int, error) {
	var count int
	stmt, err := db.Prepare(`SELECT SUM((length(message) - length(replace(message, $3, '')))::int / length($3)) AS count FROM messages WHERE channel=$1 AND username=$2`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(channel, username, query).Scan(&count)
	return count, err
}

func LastSeen(channel string, username string) (string, error) {
	var lastSeen string
	stmt, err := db.Prepare(`
		SELECT TO_CHAR(created_at, 'Day Month DD YYYY HH24:MI:SS') as last_seen 
			FROM messages WHERE channel=$1 AND username=$2 
		ORDER BY created_at DESC LIMIT 1
	`)
	if err != nil {
		return lastSeen, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(channel, username).Scan(&lastSeen)
	return spaces.ReplaceAllString(lastSeen, " "), err
}

func GetLogs(channel string, username string, count int, offset int) ([]ChatQuote, error) {
	var logs []ChatQuote
	stmt, err := db.Prepare(fmt.Sprintf(`SELECT channel, username, message, created_at FROM messages WHERE channel=$1 AND username=$2 ORDER BY created_at DESC LIMIT %d OFFSET %d`, count, offset))
	if err != nil {
		return logs, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(channel, username)
	if err != nil {
		return logs, err
	}
	for rows.Next() {
		quote := &ChatQuote{}
		if err := rows.Scan(&quote.Channel, &quote.Sender, &quote.Message, &quote.SentAt); err != nil {
			break
		}
		logs = append(logs, *quote)
	}
	rows.Close()
	return logs, nil
}
