package command

import (
	"strings"
	"time"
)

type Manager struct {
	cmds      map[string]func(Data)
	executing map[string]bool
	prefix    string
}

var mgr *Manager = &Manager{
	cmds:      make(map[string]func(Data)),
	executing: make(map[string]bool),
	prefix:    "!",
}

func GetPrefix() string {
	return mgr.prefix
}

func SetPrefix(prefix string) {
	mgr.prefix = prefix
}

func HasPrefix(str string) bool {
	return strings.HasPrefix(str, mgr.prefix)
}

func Exists(cmd string) bool {
	_, ok := mgr.cmds[strings.ToLower(Trim(cmd))]
	return ok
}

func IsExecuting(ident string) bool {
	exec, ok := mgr.executing[ident]
	return ok && exec
}

func Trim(cmd string) string {
	if HasPrefix(cmd) {
		return strings.TrimPrefix(cmd, mgr.prefix)
	}
	return cmd
}

func Register(cmd string, executor IHandler, aliases ...string) {
	cmd = strings.ToLower(Trim(cmd))
	mgr.cmds[cmd] = executor.Handle
	for _, alias := range aliases {
		Register(alias, executor)
	}
}

func Execute(exec Executor, bot IBot, channel string, sender string, cmd string, args ...string) {
	command := Data{
		Bot:       bot,
		Sender:    sender,
		Channel:   channel,
		Args:      args,
		Executor:  exec,
		CreatedAt: time.Now(),
	}
	if executor, ok := mgr.cmds[cmd]; ok {
		mgr.executing[channel] = true
		executor(command)
		delete(mgr.executing, channel)
	}
}
