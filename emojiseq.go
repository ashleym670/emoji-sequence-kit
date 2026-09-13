// Package emojiseq finds emoji sequences inside UTF-8 text.
//
// Most "emoji" a person sees on screen are not one codepoint. A
// family emoji is five codepoints joined by zero-width joiners, a
// flag is two regional-indicator letters, a skin-toned wave is a
// base emoji plus a modifier, and a keycap like 3ufe0f is a digit
// glued to a combining mark. Iterating runes or bytes splits all of
// these apart. This package walks the joining rules well enough to
// treat each of them as one unit.
package emojiseq

import (
	"strings"
	"unicode/utf8"
)

const (
	zwj             = '‍'
	variationSel16  = '️'
	keycapCombining = '⃣'

	tagStart      = '\U000E0020'
	tagEnd        = '\U000E007E'
	tagTerminator = '\U000E007F'

	regionalIndicatorStart = '\U0001F1E6'
	regionalIndicatorEnd   = '\U0001F1FF'

	skinToneStart = '\U0001F3FB'
	skinToneEnd   = '\U0001F3FF'
)

// emojiRanges are the codepoint ranges treated as emoji bases. This
// is not the full Unicode emoji-data.txt table (that's a generated
// file with thousands of entries), just the blocks that cover the
// large majority of emoji seen in everyday chat text. See the
// README for known gaps.
var emojiRanges = [][2]rune{
	{0x00A9, 0x00A9},
	{0x00AE, 0x00AE},
	{0x203C, 0x203C},
	{0x2049, 0x2049},
	{0x2122, 0x2122},
	{0x2139, 0x2139},
	{0x2194, 0x21AA},
	{0x231A, 0x231B},
	{0x2328, 0x2328},
	{0x23CF, 0x23CF},
	{0x23E9, 0x23FA},
	{0x24C2, 0x24C2},
	{0x25AA, 0x25FE},
	{0x2600, 0x27BF},
	{0x2934, 0x2935},
	{0x2B05, 0x2B07},
	{0x2B1B, 0x2B1C},
	{0x2B50, 0x2B50},
	{0x2B55, 0x2B55},
	{0x3030, 0x3030},
	{0x303D, 0x303D},
	{0x3297, 0x3297},
	{0x3299, 0x3299},
	{0x1F000, 0x1FFFF}, // emoticons, dingbats, symbols & pictographs, flags
}

func isEmojiBase(r rune) bool {
	for _, rg := range emojiRanges {
		if r >= rg[0] && r <= rg[1] {
			return true
		}
	}
	return false
}

func isRegionalIndicator(r rune) bool {
	return r >= regionalIndicatorStart && r <= regionalIndicatorEnd
}

func isSkinTone(r rune) bool {
	return r >= skinToneStart && r <= skinToneEnd
}

func isTagChar(r rune) bool {
	return r >= tagStart && r <= tagEnd
}

func isKeycapBase(r rune) bool {
	return (r >= '0' && r <= '9') || r == '#' || r == '*'
}

// consumeElement advances past a single emoji "unit": a base rune
// with an optional skin-tone modifier and/or variation selector.
// The caller has already confirmed r[i] is an emoji base.
func consumeElement(r []rune, i int) int {
	j := i + 1
	if j < len(r) && isSkinTone(r[j]) {
		j++
	}
	if j < len(r) && r[j] == variationSel16 {
		j++
	}
	return j
}

// consumeAt tries to match one full emoji sequence starting at r[i]:
// a flag pair, a keycap, a tag sequence, or a ZWJ-joined chain of
// emoji units. It returns the index just past the match and whether
// anything matched at all.
func consumeAt(r []rune, i int) (int, bool) {
	n := len(r)
	if i >= n {
		return i, false
	}

	if isRegionalIndicator(r[i]) {
		if i+1 < n && isRegionalIndicator(r[i+1]) {
			return i + 2, true
		}
		return i, false
	}

	if isKeycapBase(r[i]) && i+2 < n && r[i+1] == variationSel16 && r[i+2] == keycapCombining {
		return i + 3, true
	}

	if !isEmojiBase(r[i]) {
		return i, false
	}

	j := consumeElement(r, i)

	// Tag sequence, e.g. the black-flag base followed by ASCII-in-tag
	// letters and a terminator, used for England/Scotland/Wales.
	if j < n && isTagChar(r[j]) {
		k := j
		for k < n && isTagChar(r[k]) {
			k++
		}
		if k < n && r[k] == tagTerminator {
			return k + 1, true
		}
	}

	for j < n && r[j] == zwj && j+1 < n && isEmojiBase(r[j+1]) {
		j = consumeElement(r, j+1)
	}

	return j, true
}

// Sequences returns every emoji sequence found in s, in order.
func Sequences(s string) []string {
	r := []rune(s)
	var out []string
	for i := 0; i < len(r); {
		if end, ok := consumeAt(r, i); ok {
			out = append(out, string(r[i:end]))
			i = end
			continue
		}
		i++
	}
	return out
}

// Position marks a byte range within an input string, [Start, End),
// covering one emoji sequence found by Positions.
type Position struct {
	Start int
	End   int
}

// runeByteOffsets returns the byte offset of each rune in s, plus a
// final entry equal to len(s), so offsets[i] is where rune i begins
// and offsets[len(r)] is the end of the string.
func runeByteOffsets(s string) []int {
	offsets := make([]int, 0, len(s)+1)
	bi := 0
	for _, ch := range s {
		offsets = append(offsets, bi)
		bi += utf8.RuneLen(ch)
	}
	offsets = append(offsets, bi)
	return offsets
}

// Positions returns the byte range of every emoji sequence found in
// s, in order. Each Position can be sliced directly out of s, e.g.
// s[pos.Start:pos.End], to recover the sequence text.
func Positions(s string) []Position {
	r := []rune(s)
	offsets := runeByteOffsets(s)
	var out []Position
	for i := 0; i < len(r); {
		if end, ok := consumeAt(r, i); ok {
			out = append(out, Position{Start: offsets[i], End: offsets[end]})
			i = end
			continue
		}
		i++
	}
	return out
}

// Count returns the number of emoji sequences in s.
func Count(s string) int {
	return len(Sequences(s))
}

// Strip returns s with every emoji sequence removed.
func Strip(s string) string {
	r := []rune(s)
	var b strings.Builder
	for i := 0; i < len(r); {
		if end, ok := consumeAt(r, i); ok {
			i = end
			continue
		}
		b.WriteRune(r[i])
		i++
	}
	return b.String()
}

// IsSingleSequence reports whether s is exactly one emoji sequence
// with nothing else around it.
func IsSingleSequence(s string) bool {
	r := []rune(s)
	if len(r) == 0 {
		return false
	}
	end, ok := consumeAt(r, 0)
	return ok && end == len(r)
}
