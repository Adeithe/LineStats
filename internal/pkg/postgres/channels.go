package postgres

import "fmt"

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
