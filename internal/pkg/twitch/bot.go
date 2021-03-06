package twitch

import (
	"LineStats/internal/pkg/bitwise"
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/prometheus"

	twitch "github.com/Adeithe/go-twitch/irc"
)

type Bot struct {
	Username string
	token    string

	Channels     map[string]uint32
	LiveChannels map[string]bool

	reader *twitch.Client
	writer *twitch.Conn

	onMessage func(msg twitch.ChatMessage)
}

var _ command.IBot = &Bot{}

func New(onMessage func(msg twitch.ChatMessage)) *Bot {
	return &Bot{
		Channels:     make(map[string]uint32),
		LiveChannels: make(map[string]bool),
		onMessage:    onMessage,
	}
}

func (bot *Bot) Start(username string, token string) {
	bot.reader = twitch.New()
	bot.reader.OnShardMessage(func(shardID int, msg twitch.ChatMessage) {
		prometheus.TwitchMessagesIn.Inc()
		bot.onMessage(msg)
	})
	bot.reader.OnShardDisconnect(func(shardID int) {
		panic("reader: disconnected from twitch irc")
	})

	bot.writer = &twitch.Conn{}
	bot.writer.SetLogin(username, token)
	bot.writer.OnDisconnect(func() {
		panic("writer: disconnected from twitch irc")
	})
	if err := bot.writer.Connect(); err != nil {
		panic("writer: failed to connect to twitch irc")
	}
}

func (bot *Bot) Join(channel string, flags uint32) {
	bot.reader.Join(channel)
	bot.Channels[channel] = flags
}

func (bot *Bot) Leave(channel string) {
	bot.reader.Leave(channel)
	delete(bot.Channels, channel)
}

func (bot *Bot) InChannel(channel string) bool {
	if flags, ok := bot.Channels[channel]; ok {
		return bitwise.ShouldJoinChannel(flags)
	}
	return false
}

func (bot *Bot) Send(channel string, message string) (command.IMessage, error) {
	bot.writer.Say(channel, message)
	prometheus.TwitchMessagesOut.Inc()
	return &Message{bot}, nil
}
