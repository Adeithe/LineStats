package command

import (
	"LineStats/internal/pkg/command"
)

type Logs struct{}

var _ command.IHandler = &Logs{}

func (handler Logs) Handle(cmd command.Data) {
	if cmd.Executor != command.Discord {
		return
	}
}
