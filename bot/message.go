package bot

type IMessage interface {
	Edit(string) error
	IsEditable() bool
}
