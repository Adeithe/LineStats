package postgres

import "time"

type User struct {
	ID        int
	Name      string
	Flags     uint32
	CreatedAt time.Time
}

type Quote struct {
	ChannelID int
	SenderID  int
	Sender    string
	Message   string
	SentAt    time.Time
}

type Lines struct {
	ChannelID int
	UserID    int
	Total     int
	Unique    int
	MostDate  string
	MostCount int
	CreatedAt time.Time
	UpdatedAt time.Time
}
