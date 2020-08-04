package postgres

import "time"

type Channel struct {
	ID      int
	Name    string
	Status  bool
	AddedBy string
	AddedAt time.Time
}

func FetchChannels() ([]Channel, error) {
	var channels []Channel
	rows, err := db.Query(`SELECT id, channel, status, added_by, created_at FROM channels`)
	if err != nil {
		return channels, err
	}
	defer rows.Close()
	for rows.Next() {
		channel := Channel{}
		if err := rows.Scan(&channel.ID, &channel.Name, &channel.Status, &channel.AddedBy, &channel.AddedAt); err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func FetchEnabledChannels() ([]Channel, error) {
	var channels []Channel
	rows, err := db.Query(`SELECT id, channel, added_by, created_at FROM channels WHERE status=TRUE`)
	if err != nil {
		return channels, err
	}
	defer rows.Close()
	for rows.Next() {
		channel := Channel{Status: true}
		if err := rows.Scan(&channel.ID, &channel.Name, &channel.AddedBy, &channel.AddedAt); err != nil {
			return channels, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}
