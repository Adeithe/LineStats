package discord

import (
	"LineStats/internal/pkg/command"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	ClientID  string
	token     string
	onMessage func(*discordgo.MessageCreate)

	session *discordgo.Session
}

var _ command.IBot = &Bot{}

func New(onMessage func(*discordgo.MessageCreate)) (bot *Bot) {
	bot = &Bot{onMessage: onMessage}
	return
}

func (bot *Bot) SetLogin(clientId string, token string) {
	bot.ClientID = clientId
	bot.token = token
}

func (bot *Bot) Start() error {
	session, err := discordgo.New("Bot " + bot.token)
	if err != nil {
		return err
	}
	bot.session = session
	bot.session.AddHandler(func(session *discordgo.Session, event *discordgo.MessageCreate) { bot.onMessage(event) })
	bot.session.AddHandlerOnce(func(session *discordgo.Session, event discordgo.Disconnect) {
		panic("discord: disconnected from session")
	})
	if err := bot.session.Open(); err != nil {
		return err
	}
	return bot.session.UpdateListeningStatus("Twitch Chat")
}

func (bot *Bot) Send(channel string, message string) (command.IMessage, error) {
	msg, err := bot.session.ChannelMessageSend(channel, message)
	if err != nil {
		return &Message{}, err
	}
	return &Message{bot, msg}, nil
}
