package twitchbot

import (
	"LineStats/bot"
	"LineStats/command"
	"LineStats/command/handlers"
	"LineStats/postgres"
	"LineStats/twitch"
	"fmt"
	"regexp"
	"strings"
)

type TwitchBot struct {
	Username string
	token    string

	reader *twitch.ChatClient
	writer *twitch.ChatClient

	isConnected bool
	channels    map[string]bool
}

var _ bot.IBot = &TwitchBot{}

func New() *TwitchBot {
	return &TwitchBot{
		channels: make(map[string]bool),
	}
}

func (bot *TwitchBot) SetLogin(username string, token string) {
	bot.Username = username
	bot.token = fmt.Sprintf("oauth:%s", strings.TrimPrefix(token, "oauth:"))
}

func (bot *TwitchBot) Start() error {
	if bot.isConnected {
		fmt.Println(fmt.Errorf("twitch: starting with active connection"))
		bot.Close()
	}

	bot.reader = twitch.NewAnonymousChatClient()
	bot.reader.On(twitch.ConnectEvent, func(data interface{}) {
		bot.isConnected = true
		for channel := range bot.channels {
			bot.reader.Join(channel)
		}
	})
	bot.reader.On(twitch.DisconnectEvent, func(data interface{}) {
		//bot.isConnected = false
		//bot.Reconnect()
		panic("reader: disconnected from twitch")
	})
	bot.reader.On(twitch.ChatMessageEvent, bot.onMessage)
	bot.reader.On(twitch.RawMessageEvent, func(data interface{}) {
		if msg, ok := data.(twitch.IRCMessage); ok {
			if msg.Command == twitch.PrivMessage {
				return
			}
			line := msg.Raw
			if strings.HasPrefix(msg.Raw, "@") {
				line = strings.Join(strings.Split(msg.Raw, " ")[1:], " ")
			}
			fmt.Println(line)
		}
	})

	bot.writer = twitch.NewChatClient(bot.Username, bot.token)
	bot.writer.On(twitch.DisconnectEvent, func(data interface{}) {
		panic("writer: disconnected from twitch")
		//bot.Reconnect()
	})

	if err := bot.reader.Connect(); err != nil {
		return err
	}
	if err := bot.writer.Connect(); err != nil {
		bot.Close()
		return err
	}
	return nil
}

func (bot *TwitchBot) Reconnect() error {
	if bot.isConnected {
		bot.Close()
	}
	return bot.Start()
}

func (bot *TwitchBot) Join(channel string, commands bool) {
	var joined []string
	bot.channels[channel] = commands
	for c := range bot.channels {
		joined = append(joined, c)
	}
	if bot.isConnected {
		bot.reader.Join(joined...)
	}
}

func (bot *TwitchBot) Say(channel string, msg string) (bot.IMessage, error) {
	for _, word := range BannedWords {
		if ok, _ := regexp.MatchString(fmt.Sprintf("\\b%s\\b", word), msg); ok {
			fmt.Printf("!!! Banned Word> #%s %s: %s", channel, bot.writer.Username, msg)
			return &TwitchMessage{}, fmt.Errorf("say: banned phrase '%s' in message", word)
		}
	}
	if len(msg) > MAX_LENGTH {
		msg = msg[:MAX_LENGTH-3] + "..."
	}
	bot.writer.Send(channel, msg)
	return &TwitchMessage{}, nil
}

func (bot *TwitchBot) Sayf(channel string, format string, args ...interface{}) (bot.IMessage, error) {
	return bot.Say(channel, fmt.Sprintf(format, args...))
}

func (bot *TwitchBot) IsSimple() bool {
	return true
}

func (bot *TwitchBot) Close() {
	if bot.reader != nil {
		bot.reader.Disconnect()
		bot.reader = nil
	}
	if bot.writer != nil {
		bot.writer.Disconnect()
		bot.writer = nil
	}
	bot.isConnected = false
}

func (bot *TwitchBot) onMessage(data interface{}) {
	if !bot.isConnected {
		return
	}
	if msg, ok := data.(twitch.ChatMessage); ok {
		channel := twitch.ToChannelName(msg.Channel)
		sender := twitch.ToChannelName(msg.Sender.Username)
		message := msg.Message
		if strings.HasSuffix(message, "\U000E0000") {
			message = strings.TrimSuffix(message, "\U000E0000")
		}
		message = strings.TrimSpace(message)

		fmt.Printf("#%s %s: %s\n", channel, sender, message)
		if channel != twitch.ToChannelName(bot.writer.Username) {
			if err := postgres.SaveQuote(postgres.ChatQuote{
				Channel: channel,
				Sender:  sender,
				Message: msg.Message,
				SentAt:  msg.CreatedAt,
			}); err != nil {
				fmt.Println(err)
			}
		}

		if handlers.IsBlacklisted(sender) || handlers.IsBlacklisted(fmt.Sprint(msg.Sender.UserID)) {
			return
		}
		if !command.IsExecuting(channel) && command.HasPrefix(message) {
			if enabled, ok := bot.channels[channel]; ok && enabled {
				parts := strings.Split(message, " ")
				cmd := strings.ToLower(parts[0])
				if command.IsExist(cmd) {
					go command.Execute(bot, channel, sender, cmd, parts[1:]...)
				}
			}
		}
	}
}
