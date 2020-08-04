package handlers

import (
	"LineStats/command"
	"LineStats/postgres"
	"LineStats/twitch"
	"fmt"
	"strings"
)

type Logs struct{}

// Handle the logs command
func (handler *Logs) Handle(cmd command.Command) {
	if IsBlacklisted(cmd.Sender) {
		return
	}
	if cmd.Bot.IsSimple() {
		return
	}
	if len(cmd.Args) < 2 {
		cmd.Bot.Sayf(cmd.Channel, "**Usage:** %slogs <username> <channel>", command.GetPrefix())
		return
	}
	target := twitch.ToChannelName(cmd.Args[0])
	channel := twitch.ToChannelName(cmd.Args[1])
	msg, err := cmd.Bot.Sayf(cmd.Channel, "%s, Getting recent logs for user '%s' in channel '%s'...", cmd.Sender, target, channel)
	logs, err := postgres.GetLogs(channel, target, 100, 0)
	if err != nil {
		fmt.Printf("[logs.go:27] %s\n", err)
	}
	if len(logs) < 1 {
		if msg != nil && msg.IsEditable() {
			msg.Edit(fmt.Sprintf("Unable to find any logs for user '%s' in channel '%s'", target, channel))
			return
		}
	}
	message := "```"
	for _, quote := range logs {
		line := fmt.Sprintf("[%s UTC] %s: %s\n", quote.SentAt.Format("2006-01-02 15:04:05"), quote.Sender, quote.Message)
		if len(message)+len(line) > 2000 {
			break
		}
		message = line + message
	}
	message = "```\n" + message
	if len(message) > 2000 {
		message = "```\n" + strings.Join(strings.Split(message, "\n")[2:], "\n")
	}
	msg.Edit(message)
}
