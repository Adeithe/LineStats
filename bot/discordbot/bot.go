package discordbot

import (
	"LineStats/bot"
	"LineStats/command"
	"LineStats/command/handlers"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	ClientID string
	token    string

	session *discordgo.Session
}

var _ bot.IBot = &DiscordBot{}

var discord *DiscordBot

func New() *DiscordBot {
	return &DiscordBot{}
}

func (bot *DiscordBot) SetLogin(clientID string, token string) {
	bot.ClientID = clientID
	bot.token = token
}

func (bot *DiscordBot) Start() error {
	session, err := discordgo.New("Bot " + bot.token)
	if err != nil {
		return err
	}
	bot.session = session
	bot.session.AddHandler(bot.onMessage)
	bot.session.AddHandlerOnce(func(session *discordgo.Session, event discordgo.Disconnect) {
		panic("discord: disconnected from session")
	})
	discord = bot
	if err := bot.session.Open(); err != nil {
		return err
	}
	return bot.session.UpdateListeningStatus("Twitch Chat")
}

func (bot *DiscordBot) Reconnect() error {
	bot.Close()
	return bot.Start()
}

func (bot *DiscordBot) Close() {
	discord = nil
	if bot.session != nil {
		bot.session.Close()
	}
}

func (bot *DiscordBot) Say(channel string, msg string) (bot.IMessage, error) {
	message, err := bot.session.ChannelMessageSend(channel, msg)
	if err != nil {
		return &DiscordMessage{}, err
	}
	return NewMessage(message), err
}

func (bot *DiscordBot) Sayf(channel string, format string, args ...interface{}) (bot.IMessage, error) {
	return bot.Say(channel, fmt.Sprintf(format, args...))
}

func (bot *DiscordBot) IsSimple() bool {
	return false
}

func (bot *DiscordBot) onMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.Author.ID == session.State.User.ID {
		return
	}
	fmt.Printf("%s: %s\n", event.Author.Username, event.Content)
	channel := event.ChannelID
	sender := event.Author.Mention()
	message := event.ContentWithMentionsReplaced()

	if handlers.IsBlacklisted(event.Author.ID) {
		return
	}
	if !command.IsExecuting(channel) && command.HasPrefix(message) {
		parts := strings.Split(message, " ")
		cmd := strings.ToLower(parts[0])
		if command.IsExist(cmd) {
			go command.Execute(bot, channel, sender, cmd, parts[1:]...)
		}
	}
}
