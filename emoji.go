// Package emoji is designed to recognize and parse
// every indivisual Unicode Emoji characters from a string.
//
// Unicode Emoji Documents: http://www.unicode.org/Public/emoji/13.1/
package emoji

import (
	"bytes"

	"github.com/Andrew-M-C/go.emoji/official"
)

// ReplaceAllEmojiFunc searches string and find all emojis.
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
