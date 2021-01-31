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
			gqlUsers, err := twitch.GetUsers(ids...)
			if err != nil {
				fmt.Println(err)
			}
			for _, user := range gqlUsers {
				if len(user.Stream.ID) > 0 {
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
	if strings.HasSuffix(msg.Text, "\U000E0000") {
		msg.Text = strings.TrimSuffix(msg.Text, "\U000E0000")
	}
	msg.Text = strings.TrimSpace(msg.Text)
	if msg.IsAction {
		msg.Text = fmt.Sprintf("/me %s", msg.Text)
	}

	flags, ok := bot.Channels[channel]
	if !ok || !bitwise.ShouldJoinChannel(flags) {
		bot.Leave(channel)
		return
	}
	fmt.Printf("[%s UTC] #%s %s: %s\n", msg.CreatedAt.Format("2006-01-02 15:04:05"), msg.Channel, msg.Sender.Username, msg.Text)

	// Twitch devs are bad at their job and let things completely break sometimes so we need a failsafe in case the UserID doesn't exist.
	if msg.Sender.UserID < 1 {
		return
	}

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

	if bitwise.Has(flags, bitwise.RESPOND_TO_COMMANDS) && command.HasPrefix(msg.Text) {
		channel := twitch.ToChannelName(msg.Channel)
		parts := strings.Split(msg.Text, " ")
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
		user, err := postgres.GetTwitchUserByID(int64(msg.Sender.UserID))
		if err != nil {
			fmt.Println(err)
			return
		}
		if bitwise.Has(user.Flags, bitwise.BLACKLISTED) {
			fmt.Printf("User '%s' (ID: %d) is blacklisted from using the bot. Skipping...\n", msg.Sender.Username, msg.Sender.UserID)
			return
		}
		go command.Execute(command.Twitch, bot, channel, msg.Sender.Username, cmd, parts[1:]...)
	} else if bitwise.Has(flags, bitwise.BLOCK_PYRAMIDS) {
		if bitwise.Has(flags, bitwise.DONT_RESPOND_WHEN_LIVE) && bot.LiveChannels[fmt.Sprint(msg.ChannelID)] {
			return
		}
		handlePyramids(msg)
	}
}
