package twitchbot

import (
	"LineStats/bot"
	"fmt"
)

type TwitchMessage struct{}

var _ bot.IMessage = &TwitchMessage{}

func (msg *TwitchMessage) Edit(message string) error {
	if !msg.IsEditable() {
		return fmt.Errorf("say: message can not be edited")
	}
	return nil
}

func (msg *TwitchMessage) IsEditable() bool {
	return false
}
