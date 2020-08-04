package discordbot

import (
	"LineStats/bot"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type DiscordMessage struct {
	Message *discordgo.Message
}

var _ bot.IMessage = &DiscordMessage{}

func NewMessage(msg *discordgo.Message) *DiscordMessage {
	return &DiscordMessage{
		Message: msg,
	}
}

func (msg *DiscordMessage) Edit(message string) error {
	if !msg.IsEditable() {
		return fmt.Errorf("say: message can not be edited")
	}
	_, err := discord.session.ChannelMessageEdit(msg.Message.ChannelID, msg.Message.ID, message)
	return err
}

func (msg *DiscordMessage) IsEditable() bool {
	return msg.Message.Author.ID == discord.session.State.User.ID
}
