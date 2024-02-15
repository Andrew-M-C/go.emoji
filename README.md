# [go.emoji](https://github.com/Andrew-M-C/go.emoji)

[![Build](https://github.com/Andrew-M-C/go.emoji/actions/workflows/go_test_general.yml/badge.svg)](https://github.com/Andrew-M-C/go.emoji/actions/workflows/go_test_general.yml)
[![codecov](https://codecov.io/gh/Andrew-M-C/go.emoji/graph/badge.svg?token=9RBISZRJ3T)](https://codecov.io/gh/Andrew-M-C/go.emoji)
[![Go Report Card](https://goreportcard.com/badge/github.com/Andrew-M-C/go.emoji)](https://goreportcard.com/report/github.com/Andrew-M-C/go.emoji)
[![codebeat badge](https://codebeat.co/badges/c6f7e25f-a8fe-46a3-b4bf-a833aac65825)](https://codebeat.co/projects/github-com-andrew-m-c-go-emoji-master)

[![GoDoc](https://pkg.go.dev/badge/github.com/Andrew-M-C/go.emoji)](https://pkg.go.dev/github.com/Andrew-M-C/go.emoji@v1.1.1)
[![](https://goreportcard.com/badge/github.com/Andrew-M-C/go.emoji)](https://goreportcard.com/report/github.com/Andrew-M-C/go.emoji)
[![EmojiVer](https://img.shields.io/badge/Emoji-15.1-orange.svg)](https://www.unicode.org/Public/emoji/)
[![License](https://img.shields.io/badge/license-BSD%203--Clause-blue.svg)](https://opensource.org/licenses/BSD-3-Clause)

This Package `emoji` is designed to recognize and parse every individual Unicode Emoji characters from a string.

## Example

```go
func main() {
	s := "👩‍👩‍👦🇨🇳"
	i := 0

	final := emoji.ReplaceAllEmojiFunc(s, func(emoji string) string {
		i++
		fmt.Printf("%02d - %s - UTF-8 len %d\n", i, emoji, len(emoji))
		return fmt.Sprintf("%d-", i)
	})

	fmt.Printf("final: <%s>\n", final)
}

// Output:
// 01 - 👩‍👩‍👦 - UTF-8 len 18
// 02 - 🇨🇳 - UTF-8 len 8
// final: <1-2->
```
