package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"database/sql"
	"fmt"

	"github.com/Adeithe/go-twitch/irc"
)

type RandomQuote struct{}

var _ command.IHandler = &RandomQuote{}

func (handler RandomQuote) Handle(cmd command.Data) {
	targetUser := cmd.Sender
	targetChannel := cmd.Channel
	if len(cmd.Args) > 0 {
		targetUser = irc.ToChannelName(cmd.Args[0])
		if len(cmd.Args) > 1 {
			targetChannel = irc.ToChannelName(cmd.Args[1])
		}
	}
	channel, err := postgres.GetTwitchUserByName(targetChannel)
	if err != nil {
		fmt.Printf("Unable to get channel '%s'\n", targetChannel)
		return
	}
	user, err := postgres.GetTwitchUserByName(targetUser)
	if err != nil {
		fmt.Printf("Unable to get user '%s'\n", targetUser)
		return
	}
	quote, err := postgres.GetQuoteByUserID(channel.ID, user.ID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			fmt.Printf("Unable to find any logs for user '%s' in channel '%s'\n", user.Name, channel.Name)
		}
		return
	}
	line := fmt.Sprintf("[%s UTC] %s: %s", quote.SentAt.Format("2006-01-02 15:04:05"), user.Name, quote.Message)
	fmt.Printf("RandomQuote> %s\n", line)
	cmd.Bot.Send(cmd.Channel, line)
}
