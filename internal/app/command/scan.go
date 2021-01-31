package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"fmt"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type Scan struct{}

var _ command.IHandler = &Scan{}

func (handler Scan) Handle(cmd command.Data) {
	if len(cmd.Args) < 1 {
		return
	}
	query := cmd.Args[0]
	targetUser := cmd.Sender
	targetChannel := cmd.Channel
	if len(cmd.Args) > 1 {
		query = cmd.Args[1]
		targetUser = irc.ToChannelName(cmd.Args[0])
		if len(cmd.Args) > 2 {
			targetChannel = irc.ToChannelName(cmd.Args[2])
		}
	}
	cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, Counting occurrences of %s for user '%s' in channel '%s'...", cmd.Sender, query, targetUser, targetChannel))
	start := time.Now()
	var line string
	var fail bool
	channel, err := postgres.GetTwitchUserByName(targetChannel)
	if err != nil {
		fmt.Println(err)
		line = fmt.Sprintf("Unable to find channel '%s'", targetChannel)
		fail = true
	}
	user, err := postgres.GetTwitchUserByName(targetUser)
	if err != nil {
		fmt.Println(err)
		line = fmt.Sprintf("Unable to find user '%s'", targetUser)
		fail = true
	}
	if !fail {
		count, err := postgres.ScanMessagesByUserID(channel.ID, user.ID, query)
		if err != nil {
			fmt.Println(err)
		}
		line = fmt.Sprintf("%s, User '%s' has said %s %s times in channel '%s'", cmd.Sender, targetUser, query, format(count), targetChannel)
	}
	duration := start.Sub(time.Now())
	if duration < time.Second*3 {
		time.Sleep(time.Second*3 - duration)
	}
	cmd.Bot.Send(cmd.Channel, line)
}
