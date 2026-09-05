// Command emojiseq is a thin CLI over the emojiseq library.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	emojiseq "github.com/ashleym670/emoji-sequence-kit"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: emojiseq <command> [text]

commands:
  list    print each emoji sequence found, one per line
  count   print the number of emoji sequences found
  strip   print the input with emoji sequences removed
  check   exit 0 if the input is exactly one emoji sequence, else 1

if [text] is omitted, input is read from stdin.`)
}

func input(args []string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	data, err := io.ReadAll(bufio.NewReader(os.Stdin))
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(data), "\n"), nil
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cmd := os.Args[1]
	text, err := input(os.Args[2:])
	if err != nil {
		fmt.Fprintln(os.Stderr, "emojiseq: reading input:", err)
		os.Exit(1)
	}

	switch cmd {
	case "list":
		for _, seq := range emojiseq.Sequences(text) {
			fmt.Println(seq)
		}
	case "count":
		fmt.Println(emojiseq.Count(text))
	case "strip":
		fmt.Println(emojiseq.Strip(text))
	case "check":
		if emojiseq.IsSingleSequence(text) {
			os.Exit(0)
		}
		os.Exit(1)
	default:
		usage()
		os.Exit(2)
	}
}
