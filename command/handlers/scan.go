package handlers

import (
	"LineStats/command"
	"LineStats/postgres"
	"LineStats/twitch"
	"fmt"
	"time"
)

type Scan struct{}

// Handle the scan command
func (handler *Scan) Handle(cmd command.Command) {
	if len(cmd.Args) < 1 {
		if !cmd.Bot.IsSimple() {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %sscan <username> <query> <channel>", command.GetPrefix())
		}
		return
	}
	target := cmd.Sender
	channel := cmd.Channel
	query := cmd.Args[0]
	if len(cmd.Args) > 1 {
		target = twitch.ToChannelName(cmd.Args[0])
		query = cmd.Args[1]
	}
	if len(cmd.Args) > 2 {
		channel = twitch.ToChannelName(cmd.Args[2])
	}
	if !cmd.Bot.IsSimple() {
		if len(cmd.Args) < 3 {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %sscan <username> <query> <channel>", command.GetPrefix())
			return
		}
	}
	start := time.Now()
	msg, _ := cmd.Bot.Sayf(cmd.Channel, "%s, Counting occurrences of %s for user '%s' in channel '%s'...", cmd.Sender, query, target, channel)
	count, err := postgres.ScanMessagesByUser(channel, target, query)
	if err != nil {
		fmt.Printf("[scan.go:41] %s\n", err)
	}
	duration := start.Sub(time.Now())
	if duration < time.Second*3 {
		time.Sleep(time.Second*3 - duration)
	}
	fmt.Printf("Scan> #%s %s: %d occurrences of '%s'\n", channel, target, count, query)
	line := fmt.Sprintf("%s, User '%s' has said %s %d times in channel '%s'", cmd.Sender, target, query, count, channel)
	if msg != nil && msg.IsEditable() {
		msg.Edit(line)
		return
	}
	cmd.Bot.Sayf(cmd.Channel, line)
}
