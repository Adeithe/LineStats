package postgres

import (
	"fmt"
	"strings"
)

// GetTwitchChannels runs a function for every user that currently has flags. Runs in batches of 100.
func GetTwitchChannels(cb func(ids []string, users []User)) error {
	rows, err := db.Query(`SELECT user_id, flags FROM permissions WHERE user_type=$1 AND flags>0`, "twitch")
	if err != nil {
		return err
	}
	defer rows.Close()
	ids := []string{}
	users := []User{}
	for rows.Next() {
		var id int
		var flags uint32
		if err := rows.Scan(&id, &flags); err != nil {
			fmt.Println(err)
			continue
		}
		user, err := GetTwitchUserByID(id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		ids = append(ids, fmt.Sprint(user.ID))
		users = append(users, user)
		if len(ids) == 100 || len(users) == 100 {
			cb(ids, users)
			ids = []string{}
			users = []User{}
		}
	}
	cb(ids, users)
	return nil
}

func GetTotalLinesByRoomName(name string) (int64, error) {
	user, err := GetTwitchUserByName(strings.ToLower(name))
	if err != nil {
		return 0, err
	}
	return GetTotalLinesByRoomID(user.ID)
}

func GetTotalLinesByRoomID(roomId int) (total int64, err error) {
	stmt, err := db.Prepare(`SELECT SUM(total) as total FROM count WHERE room_id=$1`)
	if err != nil {
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(roomId).Scan(&total)
	return
}
