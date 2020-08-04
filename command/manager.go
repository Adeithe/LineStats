package command

import (
	"LineStats/bot"
	"LineStats/twitch"
	"strings"
)

type CommandManager struct {
	cmds      map[string]func(Command)
	executing map[string]bool
	prefix    string
}

var mgr *CommandManager = &CommandManager{
	cmds:      make(map[string]func(Command)),
	executing: make(map[string]bool),
	prefix:    "!",
}

func GetPrefix() string {
	return mgr.prefix
}

func SetPrefix(prefix string) {
	mgr.prefix = prefix
}

func RegisterCommand(cmd string, executor ICommandHandler, aliases ...string) {
	cmd = strings.ToLower(Trim(cmd))
	mgr.cmds[cmd] = executor.Handle
	for _, alias := range aliases {
		RegisterCommand(alias, executor)
	}
}

func UnregisterCommand(cmd string) {
	delete(mgr.cmds, cmd)
}

func IsExecuting(channel string) bool {
	executing, ok := mgr.executing[twitch.ToChannelName(channel)]
	return ok && executing
}

func HasPrefix(str string) bool {
	return strings.HasPrefix(str, mgr.prefix)
}

func IsExist(cmd string) bool {
	_, exist := mgr.cmds[strings.ToLower(Trim(cmd))]
	return exist
}

func Execute(bot bot.IBot, channel string, sender string, cmd string, args ...string) {
	channel = twitch.ToChannelName(channel)
	if executor, ok := mgr.cmds[strings.ToLower(Trim(cmd))]; ok {
		mgr.executing[channel] = true
		command := Command{
			Bot:     bot,
			Sender:  sender,
			Channel: channel,
			Args:    args,
		}
		executor(command)
		delete(mgr.executing, channel)
	}
}

func Trim(cmd string) string {
	if HasPrefix(cmd) {
		return cmd[1:]
	}
	return cmd
}
