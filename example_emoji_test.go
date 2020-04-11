package emoji_test

import (
	"fmt"

	emoji "github.com/Andrew-M-C/go.emoji"
)

func ExampleReplaceAllEmojiFunc() {
	printf := fmt.Printf

	s := "👩‍👩‍👦🇨🇳"
	i := 0

	final := emoji.ReplaceAllEmojiFunc(s, func(emoji string) string {
		i++
		printf("%02d - %s - len %d\n", i, emoji, len(emoji))
		return fmt.Sprintf("%d-", i)
	})

	printf("final: <%s>\n", final)
	// Output:
	// 01 - 👩‍👩‍👦 - len 18
	// 02 - 🇨🇳 - len 8
	// final: <1-2->
}
