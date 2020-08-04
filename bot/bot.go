package bot

type IBot interface {
	SetLogin(string, string)
	Start() error
	Reconnect() error
	Close()

	Say(string, string) (IMessage, error)
	Sayf(string, string, ...interface{}) (IMessage, error)

	IsSimple() bool
}
