package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"

	"LineStats/internal/pkg/bitwise"
)

var reader *bufio.Reader

func main() {
	var perms uint32

	if ask("Respond to commands in users chat?") {
		perms = bitwise.Set(perms, bitwise.RESPOND_TO_COMMANDS)
	}

	if ask("Save messages from users chat to database?") {
		perms = bitwise.Set(perms, bitwise.RECORD_LOGS)
	}

	if ask("Should user be allowed to use commands?") {
		perms = bitwise.Set(perms, bitwise.BLACKLISTED)
	}

	if ask("Is user an administrator?") {
		perms = bitwise.Set(perms, bitwise.ADMINISTRATOR)
	}

	if ask("Only respond to commands when user is not streaming?") {
		perms = bitwise.Set(perms, bitwise.DONT_RESPOND_WHEN_LIVE)
	}

	fmt.Printf("\nPermissions: %v", perms)
}

func ask(question string) bool {
	fmt.Printf("%s [y/N] ", question)
	in := strings.ToLower(stdin())
	return (in == "y" || in == "yes")
}

func stdin() string {
	if reader == nil {
		reader = bufio.NewReader(os.Stdin)
	}
	text, _ := reader.ReadString('\n')
	lb := "\n"
	if runtime.GOOS == "windows" {
		lb = "\r\n"
	}
	return strings.Replace(text, lb, "", -1)
}
