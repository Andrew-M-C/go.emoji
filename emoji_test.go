package emoji

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	cv = convey.Convey
	so = convey.So
	eq = convey.ShouldEqual

	notPanic = convey.ShouldNotPanic
)

func TestEmoji(t *testing.T) {
	cv("visualize unicode", t, func() { testUnicode((t)) })
	cv("ReplaceAllEmojiFunc()", t, func() { testReplaceAllEmojiFunc(t) })
	cv("IterateChars()", t, func() { testIterateChars(t) })
	cv("#3", t, func() { testIssue3(t) })
}

// reference: https://www.jianshu.com/p/9682f8ce1260
func testUnicode(t *testing.T) {
	printf := t.Logf
	buff := bytes.Buffer{}

	// ==== Emoji_presentation_sequence ====

	// thunder default mode
	buff.WriteRune(0x26A1)
	printf("%s", buff.String())
	buff.Reset()

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
	buff.WriteRune(0x1F467)
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

func testReplaceAllEmojiFunc(t *testing.T) {
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

func testIterateChars(t *testing.T) {
	cv("common case", func() {
		const s = "China🇨🇳US🇺🇸"

		i := 0
		expectRune := []string{"C", "h", "i", "n", "a", "🇨🇳", "U", "S", "🇺🇸"}
		expectBool := []bool{false, false, false, false, false, true, false, false, true}

		for it := IterateChars(s); it.Next(); i++ {
			so(it.Current(), eq, expectRune[i])
			so(it.CurrentIsEmoji(), eq, expectBool[i])
		}

		so(i, eq, len(expectRune))
	})

	cv("empty string", func() {
		so(func() {
			for it := IterateChars(""); it.Next(); {
				panic("should NOT iterate!")
			}
		}, notPanic)
	})

	cv("string without emoji", func() {
		const s = "Golang"

		i := 0
		for it := IterateChars(s); it.Next(); i++ {
			so(it.Current(), eq, string(s[i]))
		}

		so(i, eq, len(s))
	})
}

func testIssue3(t *testing.T) {
	str := string([]rune{0x2764, '|', 0x2764, 0xFE0F})
	cnt := 0

	final := ReplaceAllEmojiFunc(str, func(emoji string) string {
		cnt++
		return "heart"
	})

	so(cnt, eq, 2)
	so(final, eq, "heart|heart")

	t.Logf(str)
}
