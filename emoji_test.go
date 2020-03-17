package emoji

import (
	"bytes"
	"fmt"
	"testing"
)

// reference: https://www.jianshu.com/p/9682f8ce1260
func TestUnicode(t *testing.T) {
	printf := t.Logf
	buff := bytes.Buffer{}

	// ==== Emoji_presentation_sequence ====
	// thunder text mode
	buff.WriteRune(0x26A1)
	buff.WriteRune(0xFE0E)
	printf("%s", buff.String())
	buff.Reset()

	// thunder emoji mode
	buff.WriteRune(0x26A1)
	buff.WriteRune(0xFE0F)
	printf("%s", buff.String())
	buff.Reset()

	// ==== Modfier_Base_Sequence ====
	buff.WriteRune(0x270D)
	buff.WriteRune(0xFE0F)
	printf("%s", buff.String())
	buff.Reset()

	buff.WriteRune(0x270D)
	buff.WriteRune(0x1F3FF)
	printf("%s", buff.String())
	buff.Reset()

	// ==== Flag Sequence ====
	// reference: https://en.wikipedia.org/wiki/Regional_Indicator_Symbol
	// from 0x1F1E6 to 0x1F1FF
	buff.WriteRune(0x1F1E8)
	buff.WriteRune(0x1F1F3)
	// buff.WriteRune(0x1F1FA)
	// buff.WriteRune(0x1F1F8)
	printf("%s", buff.String())
	buff.Reset()

	// ==== KeyCap Sequence ====
	// 012345678*#
	buff.WriteRune('#')
	buff.WriteRune(0xFE0F)
	buff.WriteRune(0x20E3)
	printf("%s", buff.String())
	buff.Reset()

	// ==== ZWJ sequence ====
	buff.WriteRune(0x1F468)
	buff.WriteRune(0x200D)
	buff.WriteRune(0x1F469)
	buff.WriteRune(0x200D)
	buff.WriteRune(0x1F466)
	buff.WriteRune(0x1F46a)
	printf("%s", buff.String())
	buff.Reset()

	// ==== Tag Sequence ====
	buff.WriteRune(0x1F3F4)
	buff.WriteRune(0xE0067)
	buff.WriteRune(0xE0062)
	buff.WriteRune(0xE0065)
	buff.WriteRune(0xE006E)
	buff.WriteRune(0xE0067)
	buff.WriteRune(0xE007F)
	printf("%s", buff.String())
	buff.Reset()
}

func TestReplace(t *testing.T) {
	printf := t.Logf

	s := "👩‍👩‍👦🇨🇳"
	i := 0

	final := ReplaceAllEmojiFunc(s, func(emoji string) string {
		i++
		printf("%02d - %s - len %d", i, emoji, len(emoji))
		return fmt.Sprintf("%d-", i)
	})

	printf("final: <%s>", final)
}
