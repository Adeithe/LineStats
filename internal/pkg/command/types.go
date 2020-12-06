package command

type Data struct {
	Bot      IBot
	Sender   string
	Channel  string
	Args     []string
	Executor Executor
}

type IBot interface {
	Send(string, string) (IMessage, error)
}

type IMessage interface {
	Edit(string) error
	IsEditable() bool
}

type IHandler interface {
	Handle(Data)
}

type Executor string

const (
	Discord Executor = "DISCORD"
	Twitch  Executor = "TWITCH"
)
