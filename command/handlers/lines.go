package handlers

import (
	"LineStats/command"
	"LineStats/postgres"
	"LineStats/twitch"
	"fmt"
	"time"
)

type Lines struct{}

// Handle the lines command
func (handler *Lines) Handle(cmd command.Command) {
	target := cmd.Sender
	channel := cmd.Channel
	if len(cmd.Args) > 0 {
		target = twitch.ToChannelName(cmd.Args[0])
	}
	if len(cmd.Args) > 1 {
		channel = twitch.ToChannelName(cmd.Args[1])
	}
	if !cmd.Bot.IsSimple() {
		if len(cmd.Args) < 2 {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %slines <username> <channel>", command.GetPrefix())
			return
		}
	}
	start := time.Now()
	msg, _ := cmd.Bot.Sayf(cmd.Channel, "%s, Counting chat lines for user '%s' in channel '%s'...", cmd.Sender, target, channel)
	total, err := postgres.MessageCountByUser(channel, target)
	if err != nil {
		fmt.Printf("[lines.go:33] %s\n", err)
		return
	}
	unique, err := postgres.UniqueMessageCountByUser(channel, target)
	if err != nil {
		fmt.Printf("[lines.go:38] %s\n", err)
		return
	}
	mostDate, mostCount, err := postgres.GetMostLinesByUser(channel, target)
	if err != nil {
		fmt.Printf("[lines.go:43] %s\n", err)
		return
	}
	duration := start.Sub(time.Now())
	if duration < time.Second*3 {
		time.Sleep(time.Second*3 - duration)
	}
	if total <= 0 {
		line := fmt.Sprintf("%s, User '%s' has no lines in channel '%s'", cmd.Sender, target, channel)
		if msg.IsEditable() {
			msg.Edit(line)
			return
		}
		cmd.Bot.Sayf(cmd.Channel, line)
		return
	}
	percent := (float64(unique) / float64(total)) * 100
	fmt.Printf("Lines> #%s %s: %d lines (%.2f%% unique) - Most: %s (%d lines)\n", channel, target, total, percent, mostDate, mostCount)
	line := fmt.Sprintf("%s, User '%s' has %d lines (%.2f%% unique) in channel '%s' with their most lines in %s (%d lines)", cmd.Sender, target, total, percent, channel, mostDate, mostCount)
	if msg != nil && msg.IsEditable() {
		msg.Edit(line)
		return
	}
	cmd.Bot.Say(cmd.Channel, line)
}
