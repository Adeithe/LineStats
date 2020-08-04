package command

import "LineStats/bot"

type Command struct {
	Bot     bot.IBot
	Sender  string
	Channel string
	Args    []string
}

type ICommandHandler interface {
	Handle(Command)
}
