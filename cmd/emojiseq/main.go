// Command emojiseq is a thin CLI over the emojiseq library.
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	emojiseq "github.com/ashleym670/emoji-sequence-kit"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: emojiseq <command> [--json] [text]

commands:
  list    print each emoji sequence found, one per line
  count   print the number of emoji sequences found
  strip   print the input with emoji sequences removed
  check   exit 0 if the input is exactly one emoji sequence, else 1

--json switches the output of any command to a single JSON value on
stdout instead of the plain-text form.

if [text] is omitted, input is read from stdin.`)
}

// sequenceJSON is what "list --json" prints for each match: the
// sequence text plus the byte range it occupies in the input, so
// callers can slice the original string without re-scanning it.
type sequenceJSON struct {
	Text  string `json:"text"`
	Start int    `json:"start"`
	End   int    `json:"end"`
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
	jsonOut := false
	var rest []string
	for _, a := range os.Args[2:] {
		if a == "--json" {
			jsonOut = true
			continue
		}
		rest = append(rest, a)
	}

	text, err := input(rest)
	if err != nil {
		fmt.Fprintln(os.Stderr, "emojiseq: reading input:", err)
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)

	switch cmd {
	case "list":
		if jsonOut {
			seqs := make([]sequenceJSON, 0, len(emojiseq.Positions(text)))
			for _, pos := range emojiseq.Positions(text) {
				seqs = append(seqs, sequenceJSON{Text: text[pos.Start:pos.End], Start: pos.Start, End: pos.End})
			}
			if err := enc.Encode(seqs); err != nil {
				fmt.Fprintln(os.Stderr, "emojiseq: encoding output:", err)
				os.Exit(1)
			}
			return
		}
		for _, seq := range emojiseq.Sequences(text) {
			fmt.Println(seq)
		}
	case "count":
		if jsonOut {
			if err := enc.Encode(struct {
				Count int `json:"count"`
			}{emojiseq.Count(text)}); err != nil {
				fmt.Fprintln(os.Stderr, "emojiseq: encoding output:", err)
				os.Exit(1)
			}
			return
		}
		fmt.Println(emojiseq.Count(text))
	case "strip":
		if jsonOut {
			if err := enc.Encode(struct {
				Result string `json:"result"`
			}{emojiseq.Strip(text)}); err != nil {
				fmt.Fprintln(os.Stderr, "emojiseq: encoding output:", err)
				os.Exit(1)
			}
			return
		}
		fmt.Println(emojiseq.Strip(text))
	case "check":
		ok := emojiseq.IsSingleSequence(text)
		if jsonOut {
			if err := enc.Encode(struct {
				OK bool `json:"ok"`
			}{ok}); err != nil {
				fmt.Fprintln(os.Stderr, "emojiseq: encoding output:", err)
				os.Exit(1)
			}
		}
		if ok {
			os.Exit(0)
		}
		os.Exit(1)
	default:
		usage()
		os.Exit(2)
	}
}
