package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// references
// https://xbuba.com/questions/54450823
// https://www.jianshu.com/p/65f5139f73ad
// https://juejin.im/post/5b47564b51882519ec07e9ec
// https://www.jianshu.com/p/9682f8ce1260

const (
	emojiOfficialURL = "https://www.unicode.org/Public/emoji/"

	emojiOldDataFileVer = "12.1"

	emojiOfficialDataFile   = "emoji-data.txt"
	emojiOfficialSeqFile    = "emoji-sequences.txt"
	emojiOfficialZwjSeqFile = "emoji-zwj-sequences.txt"

	emojiDataFile   = "../../../internal/official/emoji-data.txt"
	emojiSeqFile    = "../../../internal/official/emoji-sequences.txt"
	emojiZwjSeqFile = "../../../internal/official/emoji-zwj-sequences.txt"

	emojiDataGoFile = "../../../internal/official/emoji-sequences.go"

	// typeEmojiBasic            = " Basic_Emoji "
	// typeEmojiKeycapSequence   = " Emoji_Keycap_Sequence "
	// typeEmojiFlagSequence     = " RGI_Emoji_Flag_Sequence "
	// typeEmojiTagSequence      = " RGI_Emoji_Tag_Sequence "
	// typeEmojiModifierSequence = " RGI_Emoji_Modifier_Sequence "
)

var (
	referenceURL = ""
)

type record struct {
	runes   []rune
	comment string
}

var (
	printf = log.Printf

	emojiDataFileVer string
	sequences        []record
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func addRecord(digits []rune, comment string) {
	printf("Add record: %+v, %s", digits, comment)
	if len(digits) > 0 {
		sequences = append(sequences, record{
			runes:   digits,
			comment: comment,
		})
	}
}

func main() {
	downloadEmoji()
	parseEmojiData(emojiDataFile)
	parseEmojiData(emojiSeqFile)
	parseEmojiData(emojiZwjSeqFile)
	printEmojiDataParam()
}

func parseEmojiData(filepath string) {
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	buff := bufio.NewReader(f)
	done := false
	for !done {
		line, _, err := buff.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				err = nil
				break
			}
			panic(err)
		}

		parseEmojiDataLine(string(line))
	}
}

func parseEmojiDataLine(line string) {
	if len(line) == 0 {
		return
	}

	line = strings.TrimSpace(line)
	if line[0] == '#' {
		if strings.HasPrefix(line, "# Date: ") {
			line := strings.TrimLeft(line, "# ")
			emojiDataFileVer = line
		}
		return
	}

	// printf("line ==> %s", line)
	parts := strings.Split(line, ";")
	if len(parts) < 2 {
		printf("skip line '%s'", line)
		return
	}

	// digit part
	parts[0] = strings.TrimSpace(parts[0])

	digits := strings.Split(parts[0], "..")
	if len(digits) >= 2 {
		// basic emoji unicode list with range
		min, err := strconv.ParseInt(digits[0], 16, 32)
		check(err)
		max, err := strconv.ParseInt(digits[1], 16, 32)
		check(err)

		for i := min; i <= max; i++ {
			comment := parts[1] + fmt.Sprintf(" ==> %c", rune(i))
			addRecord([]rune{rune(i)}, comment)
			addRecord([]rune{rune(i), 0xFE0E}, comment)
		}
		return
	}

	// emoji sequences
	digits = strings.Split(parts[0], " ")
	runes := make([]rune, 0, len(digits))
	for _, d := range digits {
		i, err := strconv.ParseInt(d, 16, 32)
		check(err)
		runes = append(runes, rune(i))
	}

	comment := parts[1] + fmt.Sprintf(" ==> %s", string(runes))
	addRecord(runes, comment)

	// Add U+HHHH U+FE0E sequences, which are not defined in official sequence
	if len(runes) == 1 {
		runes = append(runes, 0xFE0E)
		addRecord(runes, comment)
	}
}

func printEmojiDataParam() {
	f, err := os.OpenFile(emojiDataGoFile, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	defer f.Close()

	err = f.Truncate(0)
	check(err)

	printRecords := func(varName string, records []record, funcName string) {
		_, _ = f.WriteString("func " + funcName + "() {\n")
		defer func() {
			_, _ = f.WriteString("}\n\n")
		}()

		for _, rec := range records {
			// Fix #4: previous version also treat some ascii as emoji, which
			// should be ignored
			if len(rec.runes) == 1 && rec.runes[0] < 0x7F {
				continue
			}

			runesStr := []string{}
			for _, r := range rec.runes {
				runesStr = append(runesStr, fmt.Sprintf("0x%04x", r))
			}

			s := fmt.Sprintf(
				"\t%s.AddSequence([]rune{%s}, \"%s\")\n",
				varName,
				strings.Join(runesStr, ", "),
				strings.TrimSpace(rec.comment),
			)
			_, _ = f.WriteString(s)
		}
	}

	// file heading
	_, _ = f.WriteString("// Code generated by internal/tool/const. DO NOT EDIT.\n\n")
	_, _ = f.WriteString("// Package official indicates official unicode emoji variables.\n")
	_, _ = f.WriteString("//\n")
	_, _ = f.WriteString("// Reference: " + referenceURL + "\n")
	_, _ = f.WriteString("//\n")
	_, _ = f.WriteString("// " + emojiDataFileVer + "\n")
	_, _ = f.WriteString("package official\n\n")

	// varaibles
	printRecords("AllSequences", sequences, "initAllSequences")
}
