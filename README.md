# [go.emoji](https://github.com/Andrew-M-C/go.emoji)

[![GoDoc](https://godoc.org/github.com/Andrew-M-C/go.emoji?status.svg)](https://godoc.org/github.com/Andrew-M-C/go.emoji)
[![](https://goreportcard.com/badge/github.com/Andrew-M-C/go.emoji)](https://goreportcard.com/report/github.com/Andrew-M-C/go.emoji)
[![EmojiVer](https://img.shields.io/badge/Emoji-15.0-orange.svg)](https://www.unicode.org/Public/emoji/)
[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

This Package `emoji` is designed to recognize and parse every indivisual Unicode Emoji characters from a string.

## Example

```go
func main() {
	printf := fmt.Printf

	s := "👩‍👩‍👦🇨🇳"
	i := 0

	final := emoji.ReplaceAllEmojiFunc(s, func(emoji string) string {
		i++
		printf("%02d - %s - len %d\n", i, emoji, len(emoji))
		return fmt.Sprintf("%d-", i)
	})

	printf("final: <%s>\n", final)
}

// Output:
// 01 - 👩‍👩‍👦 - len 18
// 02 - 🇨🇳 - len 8
// final: <1-2->
```
