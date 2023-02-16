// Package emoji is designed to recognize and parse
// every indivisual Unicode Emoji characters from a string.
//
// Unicode Emoji Documents: https://www.unicode.org/Public/emoji/
package emoji

import (
	"bytes"

	"github.com/Andrew-M-C/go.emoji/internal/official"
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
		if !match {
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
