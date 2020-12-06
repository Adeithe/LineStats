package discord

import (
	"LineStats/internal/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type Message struct {
	bot *Bot
	msg *discordgo.Message
}

var _ command.IMessage = &Message{}

func (msg Message) Edit(new string) error {
	_, err := msg.bot.session.ChannelMessageEdit(msg.msg.ChannelID, msg.msg.ID, new)
	if err != nil {
		return err
	}
	return nil
}

func (msg Message) IsEditable() bool {
	return msg.IsEditable()
}
