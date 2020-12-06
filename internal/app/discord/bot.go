package discord

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/discord"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var bot *discord.Bot

func New() *discord.Bot {
	bot = discord.New(onMessage)
	return bot
}

func Start(clientId string, token string) {
	bot.SetLogin(clientId, token)
	bot.Start()
}

func onMessage(msg *discordgo.MessageCreate) {
	parts := strings.Split(msg.ContentWithMentionsReplaced(), " ")
	cmd := strings.ToLower(command.Trim(parts[0]))
	if !command.Exists(cmd) {
		return
	}
	if command.IsExecuting(cmd) {
		return
	}
	go command.Execute(command.Discord, bot, msg.ChannelID, msg.Author.Mention(), cmd, parts[1:]...)
}
