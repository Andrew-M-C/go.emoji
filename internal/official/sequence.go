package official

import (
	"bytes"
	"fmt"
)

// Sequences is the collection of Sequence type
type Sequences map[rune]*Sequence

// Sequence shows unicode Emoji sequences
type Sequence struct {
	Rune    rune
	End     bool
	Nexts   Sequences
	Comment string
}

// newSequence returns a sequence object
func newSequence(r rune) *Sequence {
	return &Sequence{
		Rune:    r,
		End:     false,
		Nexts:   Sequences{},
		Comment: "",
	}
}

func init() {
	initAllSequences()
}

// AllSequences indicates all specified unicode emoji sequences (including single basic emojis)
var AllSequences = Sequences{}

// AddSequence add a sequence identified by unicode slice. Notice: this function is NOT goroutine-safe.
func (seq Sequences) AddSequence(s []rune, comment string) {
	parentSeq := seq
	total := len(s)
	for i, r := range s {
		node, exist := parentSeq[r]
		if false == exist {
			node = newSequence(r)
			parentSeq[r] = node
		}

		if i == total-1 {
			node.End = true
			node.Comment = comment
		}

		parentSeq = node.Nexts
	}

	return
}

// HasEmojiPrefix checks whether a string is started with an emoji
func (seq Sequences) HasEmojiPrefix(s string) (has bool, length int) {
	nodes := seq
	lastRuneMatch := false

	for i, r := range s {
		if lastRuneMatch {
			lastRuneMatch = false
			length = i
		}

		node, exist := nodes[r]
		if false == exist {
			// log.Printf("End %v - %v - %s", has, length, s[:length])
			return
		}

		// log.Printf("match %s", string(r))
		if node.End {
			has = true
			lastRuneMatch = true
		}

		nodes = node.Nexts
	}

	if lastRuneMatch {
		lastRuneMatch = false
		length = len(s)
	}
	// log.Printf("End %v - %v - %s", has, length, s[:length])
	return
}

func (seq Sequences) printAllSequences() {
	for _, s := range seq {
		s.printLine([]rune{})
	}
	return
}

func (seq Sequence) printLine(parents []rune) {
	parents = append(parents, seq.Rune)

	if seq.End {
		buff := bytes.Buffer{}
		for i, r := range parents {
			if i > 0 {
				buff.WriteRune(' ')
			}
			buff.WriteString(fmt.Sprintf("U+%04x", r))
		}
		fmt.Printf("%s \t- %s\n", string(parents), buff.String())
	}

	for _, r := range seq.Nexts {
		r.printLine(parents)
	}
	return
}
