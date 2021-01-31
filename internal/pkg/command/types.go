package command

import (
	"time"
)

type Data struct {
	Bot       IBot
	Sender    string
	Channel   string
	Args      []string
	Executor  Executor
	CreatedAt time.Time
}

type IBot interface {
	Send(string, string) (IMessage, error)
}

type IMessage interface {
	Edit(string) error
	IsEditable() bool
}

type IHandler interface {
	Handle(Data)
}

type Executor string

const (
	Discord Executor = "DISCORD"
	Twitch  Executor = "TWITCH"
)
