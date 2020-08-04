package handlers

import (
	"LineStats/bot"
	"LineStats/command"
	"LineStats/postgres"
	"LineStats/twitch"
	"database/sql"
	"fmt"
)

type RandomQuote struct{}

// Handle the randomquote command
func (handler *RandomQuote) Handle(cmd command.Command) {
	target := cmd.Sender
	channel := cmd.Channel
	if len(cmd.Args) > 0 {
		target = twitch.ToChannelName(cmd.Args[0])
	}
	if len(cmd.Args) > 1 {
		channel = twitch.ToChannelName(cmd.Args[1])
	}
	var msg bot.IMessage
	if !cmd.Bot.IsSimple() {
		if len(cmd.Args) < 2 {
			cmd.Bot.Sayf(cmd.Channel, "**Usage:** %srandomquote <username> <channel>", command.GetPrefix())
			return
		}
		message, err := cmd.Bot.Sayf(cmd.Channel, "Getting a random quote for user '%s' in channel '%s'...", target, channel)
		if err == nil {
			msg = message
		}
	}
	quote, err := postgres.RandomQuoteByUser(channel, target)
	if err != nil {
		if err == sql.ErrNoRows {
			if msg != nil && msg.IsEditable() {
				msg.Edit(fmt.Sprintf("Unable to find any logs for user '%s' in channel '%s'", target, channel))
			}
		}
		fmt.Printf("[randomquote.go:42] %s\n", err)
		return
	}
	line := fmt.Sprintf("[%s UTC] %s: %s", quote.SentAt.Format("2006-01-02 15:04:05"), quote.Sender, quote.Message)
	fmt.Println("Random Quote> " + line)
	if !cmd.Bot.IsSimple() {
		line = "`" + line + "`"
		if msg != nil && msg.IsEditable() {
			msg.Edit(line)
			return
		}
	}
	cmd.Bot.Say(cmd.Channel, line)
}
