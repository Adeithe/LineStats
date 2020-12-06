package command

import "LineStats/internal/pkg/command"

type NameHistory struct{}

var _ command.IHandler = &NameHistory{}

func (handler NameHistory) Handle(cmd command.Data) {
	if cmd.Executor != command.Discord {
		return
	}
}
