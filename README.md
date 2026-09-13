# emoji-sequence-kit

A Go library and CLI for finding emoji sequences in text.

## The problem

Most emoji you see rendered as a single glyph are not a single
codepoint:

- 👨‍👩‍👧‍👦 (family) is four base emoji joined by three zero-width joiners
- 🇨🇦 (flag) is two regional-indicator letters, `C` and `A`
- 👋🏽 (waving hand, medium skin tone) is a base emoji plus a modifier
- 3️⃣ (keycap three) is the digit `3`, a variation selector, and a
  combining enclosing keycap mark

If you `range` over a Go string, or take `len([]rune(s))`, each of
these gets torn into two to seven pieces. Anything that counts
emoji, extracts them from a message, or checks "is this string just
one emoji" needs to know the joining rules, not just the codepoints.

## What's here

`emojiseq.go` implements the joining rules for ZWJ sequences, flag
pairs, skin-tone modifiers, keycaps, and tag-based subdivision flags
(England, Scotland, Wales). `cmd/emojiseq` wraps it in a CLI.

## Library usage

```go
package main

import (
	"fmt"

	emojiseq "github.com/ashleym670/emoji-sequence-kit"
)

func main() {
	msg := "movie night 🍿 with the family 👨‍👩‍👧‍👦 tonight!"

	fmt.Println(emojiseq.Count(msg))       // 2
	fmt.Println(emojiseq.Sequences(msg))   // [🍿 👨‍👩‍👧‍👦]
	fmt.Println(emojiseq.Strip(msg))       // "movie night  with the family  tonight!"

	fmt.Println(emojiseq.IsSingleSequence("👋🏽")) // true
	fmt.Println(emojiseq.IsSingleSequence("👋🏽!")) // false

	for _, pos := range emojiseq.Positions(msg) {
		fmt.Println(msg[pos.Start:pos.End]) // 🍿, then 👨‍👩‍👧‍👦
	}
}
```

## CLI usage

```
$ echo "good morning 🇨🇦 3️⃣ times" | go run ./cmd/emojiseq list
🇨🇦
3️⃣

$ echo "good morning 🇨🇦 3️⃣ times" | go run ./cmd/emojiseq count
2

$ go run ./cmd/emojiseq check "👋🏽"; echo $?
0
```

If no text argument is given, the CLI reads from stdin.

## Known gaps

The emoji base ranges in `emojiRanges` cover the common Unicode
blocks (emoticons, dingbats, symbols and pictographs, transport,
flags) rather than the full generated `emoji-data.txt` table from
unicode.org, so some rarer single-codepoint emoji won't be
recognized as sequence starts yet. See the roadmap below.

## License

MIT, see LICENSE.
