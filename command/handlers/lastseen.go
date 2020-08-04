package handlers

import (
	"LineStats/bot"
	"LineStats/command"
	"LineStats/postgres"
	"LineStats/twitch"
	"database/sql"
	"fmt"
)

type LastSeen struct{}

// Handle the lastseen command
func (handler *LastSeen) Handle(cmd command.Command) {
	if len(cmd.Args) < 1 {
		if !cmd.Bot.IsSimple() {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %slastseen <username> <channel>", command.GetPrefix())
		}
		return
	}
	target := twitch.ToChannelName(cmd.Args[0])
	channel := cmd.Channel
	if len(cmd.Args) > 1 {
		channel = twitch.ToChannelName(cmd.Args[1])
	}
	var msg bot.IMessage
	if !cmd.Bot.IsSimple() {
		if len(cmd.Args) < 2 {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %slastseen <username> <channel>", command.GetPrefix())
			return
		}
		message, err := cmd.Bot.Sayf(cmd.Channel, "Checking to see when '%s' was last seen in in channel '%s'...", target, channel)
		if err == nil {
			msg = message
		}
	}
	lastSeen, err := postgres.LastSeen(channel, target)
	if err != nil {
		if err == sql.ErrNoRows {
			line := fmt.Sprintf("%s, User '%s' has never typed in channel '%s'", cmd.Sender, target, channel)
			if msg != nil && msg.IsEditable() {
				msg.Edit(line)
			} else {
				cmd.Bot.Say(cmd.Channel, line)
			}
		}
		fmt.Printf("[lastseen.go:48] %s\n", err)
		return
	}
	line := fmt.Sprintf("%s, User '%s' was last seen in channel '%s' on %s UTC", cmd.Sender, target, channel, lastSeen)
	if msg != nil && msg.IsEditable() {
		msg.Edit(line)
		return
	}
	cmd.Bot.Say(cmd.Channel, line)
}
