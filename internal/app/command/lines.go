package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"fmt"
	"time"

	"github.com/Adeithe/go-twitch/irc"
)

type Lines struct{}

var _ command.IHandler = &Lines{}

func (handler Lines) Handle(cmd command.Data) {
	targetUser := cmd.Sender
	targetChannel := cmd.Channel
	if len(cmd.Args) > 0 {
		targetUser = irc.ToChannelName(cmd.Args[0])
		if len(cmd.Args) > 1 {
			targetChannel = irc.ToChannelName(cmd.Args[1])
		}
	}
	var fail bool
	cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, Counting chat lines for user '%s' in channel '%s'...", cmd.Sender, targetUser, targetChannel))
	start := time.Now()
	channel, err := postgres.GetTwitchUserByName(targetChannel)
	if err != nil {
		fail = true
		fmt.Printf("Unable to get channel '%s'\n", targetChannel)
		return
	}
	user, err := postgres.GetTwitchUserByName(targetUser)
	if err != nil {
		fail = true
		fmt.Printf("Unable to get user '%s'\n", targetUser)
		return
	}
	lines := postgres.Lines{}
	if !fail {
		data, err := postgres.GetLinesByUserID(channel.ID, user.ID)
		if err != nil {
			fmt.Println(err)
		}
		lines = data
	}
	duration := start.Sub(time.Now())
	if duration < time.Second*3 {
		time.Sleep(time.Second*3 - duration)
	}
	if lines.Total <= 0 {
		line := fmt.Sprintf("%s, User '%s' has no lines in channel '%s'", cmd.Sender, targetUser, targetChannel)
		cmd.Bot.Send(cmd.Channel, line)
		return
	}
	percent := (float64(lines.Unique) / float64(lines.Total)) * 100
	fmt.Printf("Lines> #%s %s: %d lines (%.2f%% unique) - Most: %s (%d lines)\n", channel.Name, targetUser, lines.Total, percent, lines.MostDate, lines.MostCount)
	line := fmt.Sprintf("%s, User '%s' has %s lines (%.2f%% unique) in channel '%s' with their most lines in %s (%s lines)", cmd.Sender, targetUser, format(lines.Total), percent, targetChannel, lines.MostDate, format(int64(lines.MostCount)))
	cmd.Bot.Send(cmd.Channel, line)
}
