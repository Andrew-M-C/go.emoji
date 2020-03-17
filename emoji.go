// Package emoji parse emoji in strings
//
// Documents:
//
// emoji-data.txt: http://www.unicode.org/Public/emoji/5.0/emoji-data.txt
// emoji-sequences.txt: http://unicode.org/Public/emoji/5.0/emoji-sequences.txt
package emoji

import (
	"bytes"

	"github.com/Andrew-M-C/go.emoji/official"
)

// ReplaceAllEmojiFunc search string and find all emojis
func ReplaceAllEmojiFunc(s string, f func(emoji string) string) string {
	buff := bytes.Buffer{}
	nextIndex := 0

	for i, r := range s {
		if i < nextIndex {
			continue
		}

		match, length := official.AllSequences.HasEmojiPrefix(s[i:])
		if false == match {
			buff.WriteRune(r)
			continue
		}

		nextIndex = i + length
		if f != nil {
			buff.WriteString(f(s[i : i+length]))
		}
	}

	return buff.String()
}
