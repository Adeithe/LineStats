package postgres

import "time"

type User struct {
	ID        int64
	Name      string
	Flags     uint32
	CreatedAt time.Time
}

type Quote struct {
	ChannelID int64
	SenderID  int64
	Sender    string
	Message   string
	SentAt    time.Time
}

type Lines struct {
	ChannelID int64
	UserID    int64
	Total     int64
	Unique    int64
	MostDate  string
	MostCount int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
