package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"database/sql"
	"fmt"
)

type LastSeen struct{}

var _ command.IHandler = &LastSeen{}

func (handler LastSeen) Handle(cmd command.Data) {
	if len(cmd.Args) < 1 {
		return
	}
	targetUser := toChannelName(cmd.Args[0])
	targetChannel := cmd.Channel
	if len(cmd.Args) > 1 {
		targetChannel = toChannelName(cmd.Args[1])
	}
	channel, err := postgres.GetTwitchUserByName(targetChannel)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			cmd.Bot.Send(cmd.Channel, fmt.Sprintf("Unable to find channel '%s'", targetChannel))
		}
		return
	}
	user, err := postgres.GetTwitchUserByName(targetUser)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			cmd.Bot.Send(cmd.Channel, fmt.Sprintf("Unable to find user '%s'", targetUser))
		}
		return
	}
	lastSeen, err := postgres.GetLastSeenByUserID(channel.ID, user.ID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, User '%s' has never typed in channel '%s'", cmd.Sender, targetUser, targetChannel))
		}
		return
	}
	cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, User '%s' was last seen in channel '%s' on %s UTC", cmd.Sender, targetUser, targetChannel, lastSeen))
}
