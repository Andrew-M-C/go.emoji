package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	emojiDataFile   = "../../official/emoji-sequences.txt"
	emojiZwjSeqFile = "../../official/emoji-zwj-sequences.txt"

	emojiDataGoFile = "../../official/emoji-sequences.go"

	typeEmojiBasic            = " Basic_Emoji "
	typeEmojiKeycapSequence   = " Emoji_Keycap_Sequence "
	typeEmojiFlagSequence     = " RGI_Emoji_Flag_Sequence "
	typeEmojiTagSequence      = " RGI_Emoji_Tag_Sequence "
	typeEmojiModifierSequence = " RGI_Emoji_Modifier_Sequence "
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
	return
}

func main() {
	parseEmojiData(emojiDataFile)
	parseEmojiData(emojiZwjSeqFile)
	printEmojiDataParam()
	return
}

func parseEmojiData(filepath string) {
	f, err := os.Open(filepath)
	check(err)
	defer f.Close()

	buff := bufio.NewReader(f)
	done := false
	for false == done {
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
	if 0 == len(line) {
		return
	}

	line = strings.TrimSpace(line)
	if '#' == line[0] {
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
	if len(digits) > 1 {
		// basic emoji unicode list with range
		for _, d := range digits {
			i, err := strconv.ParseInt(d, 16, 32)
			check(err)
			comment := parts[1] + fmt.Sprintf(" ==> %c", rune(i))
			addRecord([]rune{rune(i)}, comment)
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

	return
}

func printEmojiDataParam() {
	f, err := os.OpenFile(emojiDataGoFile, os.O_RDWR|os.O_CREATE, 0755)
	check(err)
	defer f.Close()

	err = f.Truncate(0)
	check(err)

	printRecords := func(varName string, records []record, funcName string) {
		f.WriteString("func " + funcName + "() {\n")
		defer f.WriteString("}\n\n")

		for _, rec := range records {
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
			f.WriteString(s)
		}
		return
	}

	// file heading
	f.WriteString("// Package official indicates official unicode emoji variables\n")
	f.WriteString("// Reference: https://www.unicode.org/Public/emoji/13.0/emoji-sequences.txt\n")
	f.WriteString("// " + emojiDataFileVer + "\n")
	f.WriteString("package official\n\n")

	// varaibles
	printRecords("AllSequences", sequences, "initAllSequences")

	return
}
