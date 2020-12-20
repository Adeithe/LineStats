package command

import (
	"LineStats/internal/pkg/command"
	"LineStats/internal/pkg/postgres"
	"fmt"
	"strconv"
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
	total := handler.format(n)
	if err != nil {
		fmt.Println(err)
		cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, Lines are not being stored for channel %s.", cmd.Sender, targetChannel))
		return
	}
	cmd.Bot.Send(cmd.Channel, fmt.Sprintf("%s, A total of %s lines have been stored for channel %s", cmd.Sender, total, targetChannel))
}

func (handler TotalLines) format(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits--
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
