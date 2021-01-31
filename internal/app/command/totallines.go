package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"fmt"
	"strings"
)

type TotalLines struct{}

var _ command.IHandler = &TotalLines{}

func (handler TotalLines) Handle(cmd command.Data) {
	if cmd.Executor == command.Discord && len(cmd.Args) < 1 {
		return
	}
	targetChannel := cmd.Channel
	if len(cmd.Args) > 0 {
		targetChannel = strings.ToLower(cmd.Args[0])
	}
	n, err := postgres.GetTotalLinesByRoomName(targetChannel)
	total := format(n)
	if err != nil {
		fmt.Println(err)
		cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, Lines are not being stored for channel %s.", cmd.Sender, targetChannel))
		return
	}
	if all, err := postgres.GetTotalLinesStored(); err == nil {
		percent := (float64(n) / float64(all)) * 100
		cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, A total of %s lines have been stored for channel %s (%.2f%% of all lines)", cmd.Sender, total, targetChannel, percent))
		return
	}
	cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, A total of %s lines have been stored for channel %s", cmd.Sender, total, targetChannel))
}
