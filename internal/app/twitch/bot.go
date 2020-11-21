package twitch

import (
	"fmt"
	"strings"
	"time"

	"LineStats/internal/pkg/bitwise"
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"LineStats/internal/pkg/twitch"

	ttv "github.com/Adeithe/go-twitch/irc"
)

var bot *twitch.Bot
var started bool

func New() *twitch.Bot {
	bot = twitch.New(onMessage)
	return bot
}

func Start(username string, token string) {
	bot.Start(username, token)
	if started {
		return
	}
	go ticker()
	started = true
}

func ticker() {
	for {
		live := []string{}
		postgres.GetTwitchChannels(func(ids []string, users []postgres.User) {
			if users, err := twitch.GetUsers(ids...); err == nil {
				for _, user := range users {
					userID := fmt.Sprint(user.ID)
					bot.LiveChannels[userID] = true
					live = append(live, userID)
				}
			}
			for _, user := range users {
				if !bitwise.ShouldJoinChannel(user.Flags) {
					if bot.InChannel(user.Name) {
						bot.Leave(user.Name)
					}
					continue
				}
				bot.Join(user.Name, user.Flags)
			}
		})
		keys := []string{}
		for channel := range bot.LiveChannels {
			var found bool
			for _, id := range live {
				if id == channel {
					found = true
					break
				}
			}
			if !found {
				keys = append(keys, channel)
			}
		}
		for _, key := range keys {
			delete(bot.LiveChannels, key)
		}
		time.Sleep(5 * time.Minute)
	}
}

func onMessage(msg ttv.ChatMessage) {
	channel := msg.Channel
	//sender := msg.Sender.Username
	if strings.HasSuffix(msg.Message, "\U000E0000") {
		msg.Message = strings.TrimSuffix(msg.Message, "\U000E0000")
	}
	msg.Message = strings.TrimSpace(msg.Message)
	if msg.IsAction {
		msg.Message = fmt.Sprintf("/me %s", msg.Message)
	}

	flags, ok := bot.Channels[channel]
	if !ok || !bitwise.ShouldJoinChannel(flags) {
		bot.Leave(channel)
		return
	}

	fmt.Printf("[%s UTC] #%s %s: %s\n", msg.CreatedAt.Format("2006-01-02 15:04:05"), msg.Channel, msg.Sender.Username, msg.Message)
	if bitwise.Has(flags, bitwise.RECORD_LOGS) {
		if err := postgres.SaveQuote(msg); err != nil {
			fmt.Println(err)
		}
	} else {
		// We need to save users anyway for compatibility and name lookup
		if err := postgres.SaveUser(msg.Sender); err != nil {
			fmt.Println(err)
		}
	}

	if command.HasPrefix(msg.Message) && bitwise.Has(flags, bitwise.RESPOND_TO_COMMANDS) {
		channel := twitch.ToChannelName(msg.Channel)
		parts := strings.Split(msg.Message, " ")
		cmd := strings.ToLower(command.Trim(parts[0]))
		if !command.Exists(cmd) {
			return
		}
		if command.IsExecuting(channel) {
			fmt.Printf("Channel '%s' is already executing a command. Skipping...\n", channel)
			return
		}
		if bitwise.Has(flags, bitwise.DONT_RESPOND_WHEN_LIVE) && bot.LiveChannels[fmt.Sprint(msg.ChannelID)] {
			fmt.Printf("Channel '%s' is live. Skipping command as requested by the streamer...\n", channel)
			return
		}
		user, err := postgres.GetTwitchUserByID(msg.Sender.UserID)
		if err != nil {
			fmt.Println(err)
			return
		}
		if bitwise.Has(user.Flags, bitwise.BLACKLISTED) {
			fmt.Printf("User '%s' (ID: %d) is blacklisted from using the bot. Skipping...\n", msg.Sender.Username, msg.Sender.UserID)
			return
		}
		go command.Execute(command.Twitch, bot, channel, msg.Sender.Username, cmd, parts[1:]...)
	}
}
