package emojiseq

import (
	"reflect"
	"testing"
)

// Codepoints are spelled out with \U escapes rather than pasted as
// literal glyphs so the exact composition of each sequence (which
// joiners, which modifiers, in what order) is visible in the diff.
const (
	man   = "\U0001F468"
	woman = "\U0001F469"
	girl  = "\U0001F467"
	boy   = "\U0001F466"
	zwjS  = "\U0000200D"

	wavingHand = "\U0001F44B"
	skinMedium = "\U0001F3FD"
	vs16       = "\U0000FE0F"
	keycapMark = "\U000020E3"
	blackFlag  = "\U0001F3F4"
	riC        = "\U0001F1E8"
	riA        = "\U0001F1E6"
	tagG       = "\U000E0067"
	tagB       = "\U000E0062"
	tagE       = "\U000E0065"
	tagN       = "\U000E006E"
	tagCancel  = "\U000E007F"
)

var (
	family      = man + zwjS + woman + zwjS + girl + zwjS + boy
	flagCanada  = riC + riA
	keycapThree = "3" + vs16 + keycapMark
	waveMedium  = wavingHand + skinMedium
	flagEngland = blackFlag + tagG + tagB + tagE + tagN + tagG + tagCancel
)

func TestSequencesSingleKinds(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"family zwj chain", family, []string{family}},
		{"flag pair", flagCanada, []string{flagCanada}},
		{"keycap", keycapThree, []string{keycapThree}},
		{"skin tone modifier", waveMedium, []string{waveMedium}},
		{"tag sequence flag", flagEngland, []string{flagEngland}},
		{"plain text", "hello world", nil},
		{"lone regional indicator", riC, nil},
		{"digit without keycap parts", "3", nil},
		{"digit with vs16 but no keycap mark", "3" + vs16, nil},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Sequences(c.in)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("Sequences(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestSequencesMixedText(t *testing.T) {
	in := "movie night \U0001F37F with the family " + family + " tonight " + flagCanada + "!"
	want := []string{"\U0001F37F", family, flagCanada}
	got := Sequences(in)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Sequences(%q) = %v, want %v", in, got, want)
	}
}

func TestCount(t *testing.T) {
	in := family + " and " + flagCanada + " and " + keycapThree
	if got := Count(in); got != 3 {
		t.Errorf("Count(%q) = %d, want 3", in, got)
	}
	if got := Count("no emoji here"); got != 0 {
		t.Errorf("Count with no emoji = %d, want 0", got)
	}
}

func TestStrip(t *testing.T) {
	in := "movie night " + family + " tonight " + flagCanada
	want := "movie night  tonight "
	if got := Strip(in); got != want {
		t.Errorf("Strip(%q) = %q, want %q", in, got, want)
	}
}

func TestIsSingleSequence(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{family, true},
		{flagCanada, true},
		{keycapThree, true},
		{waveMedium, true},
		{flagEngland, true},
		{family + "!", false},
		{"!" + flagCanada, false},
		{"", false},
		{"plain", false},
		{riC, false},
	}

	for _, c := range cases {
		if got := IsSingleSequence(c.in); got != c.want {
			t.Errorf("IsSingleSequence(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}
