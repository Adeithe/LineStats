package command

import (
	"LineStats/internal/pkg/command"
	"strconv"
	"strings"
)

func Init() {
	command.Register("lastseen", &LastSeen{})
	command.Register("lines", &Lines{})
	command.Register("randomquote", &RandomQuote{}, "rq")
	command.Register("scan", &Scan{})
	command.Register("totallines", &TotalLines{})

	// Discord ONLY commands
	command.Register("logs", &Logs{})
	command.Register("namehistory", &NameHistory{}, "names", "nh")
}

func toChannelName(str string) string {
	return strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(strings.TrimSuffix(str, ","), "@"), "#"))
}

func format(n int64) string {
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
