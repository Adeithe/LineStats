package postgres

import (
	"database/sql"
	"strings"

	twitch "github.com/Adeithe/go-twitch/irc"
)

func SaveQuote(msg twitch.ChatMessage) error {
	stmt, err := db.Prepare(`INSERT INTO messages(room_id, user_id, username, message, created_at) VALUES($1, $2, $3, $4, $5)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(msg.ChannelID, msg.Sender.ID, msg.Sender.Username, msg.Text, msg.CreatedAt.UTC())
	return err
}

func SaveUser(user twitch.ChatSender) error {
	stmt, err := db.Prepare(`INSERT INTO users(id, name) VALUES($1, $2)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID, user.Username)
	return err
}

func GetTwitchUserByID(userId int64) (User, error) {
	user := User{ID: -1}
	stmt, err := db.Prepare(`SELECT id, name, created_at FROM users WHERE id=$1 ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(userId).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		return user, err
	}
	user.Flags, _ = _GetFlagsForUserByID("twitch", user.ID)
	return user, nil
}

func GetTwitchUserByName(login string) (User, error) {
	user := User{}
	stmt, err := db.Prepare(`SELECT id, name, created_at FROM users WHERE name=$1 ORDER BY created_at DESC LIMIT 1`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(login).Scan(&user.ID, &user.Name, &user.CreatedAt); err != nil {
		return user, err
	}
	userByID, _ := GetTwitchUserByID(user.ID)
	if userByID.ID > 0 {
		user = userByID
	}
	return user, nil
}

func GetLastSeenByUserID(roomId int64, userId int64) (string, error) {
	var lastSeen string
	stmt, err := db.Prepare(`SELECT TO_CHAR(updated_at, 'Day Month DD YYYY HH24:MI:SS') FROM count WHERE room_id=$1 AND user_id=$2`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	err = stmt.QueryRow(roomId, userId).Scan(&lastSeen)
	return spaces.ReplaceAllString(lastSeen, " "), err
}

func GetLinesByUserID(roomId int64, userId int64) (Lines, error) {
	lines := &Lines{}
	count, err := db.Prepare(`SELECT room_id, user_id, total, created_at, updated_at FROM count WHERE room_id=$1 AND user_id=$2`)
	if err != nil {
		return *lines, err
	}
	defer count.Close()
	if err := count.QueryRow(roomId, userId).Scan(&lines.ChannelID, &lines.UserID, &lines.Total, &lines.CreatedAt, &lines.UpdatedAt); err != nil {
		return *lines, err
	}
	unique, err := db.Prepare(`SELECT COUNT(DISTINCT message) FROM messages WHERE room_id=$1 AND user_id=$2`)
	if err != nil {
		return *lines, err
	}
	defer unique.Close()
	if err := unique.QueryRow(roomId, userId).Scan(&lines.Unique); err != nil {
		return *lines, err
	}
	most, err := db.Prepare(`SELECT TO_CHAR(created_at, 'Month YYYY') as mon, COUNT(*) FROM messages WHERE room_id=$1 AND user_id=$2 GROUP BY mon`)
	if err != nil {
		return *lines, err
	}
	defer most.Close()
	rows, err := most.Query(roomId, userId)
	if err != nil {
		return *lines, err
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		var count int64
		if err := rows.Scan(&date, &count); err != nil {
			break
		}
		if count > lines.MostCount {
			lines.MostDate = spaces.ReplaceAllString(date, " ")
			lines.MostCount = count
		}
	}
	return *lines, nil
}

func ScanMessagesByUserID(roomId int64, userId int64, query string) (int64, error) {
	var count int64
	stmt, err := db.Prepare(`SELECT SUM((length(message) - length(replace(message, $3, '')))::int / length($3)) AS count FROM messages WHERE room_id=$1 AND user_id=$2`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(roomId, userId, query).Scan(&count)
	return count, err
}

func GetQuoteByUserID(roomId int64, userId int64) (Quote, error) {
	quote := &Quote{}
	stmt, err := db.Prepare(`
		SELECT room_id, user_id, username, message, created_at FROM messages WHERE room_id=$1 AND user_id=$2 
			OFFSET floor(random() * (
				SELECT total FROM count WHERE room_id=$1 AND user_id=$2
			)) LIMIT 1;
	`)
	if err != nil {
		return *quote, err
	}
	defer stmt.Close()
	if err := stmt.QueryRow(roomId, userId).Scan(&quote.ChannelID, &quote.SenderID, &quote.Sender, &quote.Message, &quote.SentAt); err != nil {
		return *quote, err
	}
	return *quote, nil
}

func _GetFlagsForUserByID(userType string, userID int64) (uint32, error) {
	userType = strings.ToLower(userType)
	stmt, err := db.Prepare(`SELECT flags FROM permissions WHERE user_type=$1 AND user_id=$2`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var flags uint32
	if err := stmt.QueryRow(userType, userID).Scan(&flags); err != nil {
		if err == sql.ErrNoRows {
			if _, err := db.Exec(`INSERT INTO permissions(user_type, user_id) VALUES($1, $2)`, userType, userID); err == nil {
				return _GetFlagsForUserByID(userType, userID)
			}
		}
		return 0, err
	}
	return flags, nil
}
