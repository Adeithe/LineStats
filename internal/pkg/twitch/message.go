package twitch

import (
	"LineStats/internal/pkg/command"

	"errors"
)

type Message struct {
	bot *Bot
}

var _ command.IMessage = &Message{}

func (msg Message) Edit(string) error {
	return errors.New("message can not be edited")
}

func (msg Message) IsEditable() bool {
	return false
}
