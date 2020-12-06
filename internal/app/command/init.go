package command

import "LineStats/internal/pkg/command"

func Init() {
	command.Register("lastseen", &LastSeen{})
	command.Register("lines", &Lines{})
	command.Register("randomquote", &RandomQuote{}, "rq")
	command.Register("scan", &Scan{})

	// Discord ONLY commands
	command.Register("logs", &Logs{})
	command.Register("namehistory", &NameHistory{}, "names", "nh")
}
