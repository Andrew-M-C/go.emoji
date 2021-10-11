package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func downloadEmoji() {
	// list and find latest emoji versions
	s, err := httpGet(emojiOfficialURL)
	var lines []string
	check(err)

	if find := strings.SplitN(s, "Parent Directory", 2); len(find) == 1 {
		printf("error, cannot find emoji sub directories")
		os.Exit(-1)
	} else {
		lines = strings.Split(find[1], "\n")
	}

	latest := ""
	rep := regexp.MustCompile(`(\d+\.\d+)\/.+20\d\d-[01]\d-\d\d`)
	for _, line := range lines {
		match := rep.FindStringSubmatch(line)
		if len(match) < 2 {
			continue
		}

		v := match[1]
		printf("found version: %s", v)

		if vercmp(latest, v) < 0 {
			latest = v
		}
	}

	printf("latest version: %v", latest)

	// get emoji and zwj files
	basic, err := httpGet(emojiOfficialURL + latest + "/" + emojiOfficialSeqFile)
	check(err)
	zwj, err := httpGet(emojiOfficialURL + latest + "/" + emojiOfficialZwjSeqFile)
	check(err)

	err = ioutil.WriteFile(emojiDataFile, []byte(basic), 0644)
	check(err)
	err = ioutil.WriteFile(emojiZwjSeqFile, []byte(zwj), 0644)
	check(err)

	referenceURL = emojiOfficialURL + latest + "/"
}

func vercmp(left, right string) int {
	if left == "" {
		if right == "" {
			return 0
		}
		return -1
	}
	if right == "" {
		return 1
	}

	toCode := func(s string) uint32 {
		nums := strings.Split(s, ".")
		if len(nums) == 1 {
			n, _ := strconv.ParseUint(s, 10, 32)
			return uint32(n) << 16
		}

		hi, _ := strconv.ParseUint(nums[0], 10, 32)
		lo, _ := strconv.ParseUint(nums[1], 10, 32)
		return (uint32(hi) << 16) + uint32(lo)
	}

	leftCode := toCode(left)
	rightCode := toCode(right)

	if leftCode < rightCode {
		return -1
	}
	if leftCode > rightCode {
		return 1
	}
	return 0
}

func httpGet(u string) (res string, err error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", fmt.Errorf("http.Get(%s) error: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP status %d, body '%s'", resp.StatusCode, string(b))
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return string(b), nil
		}
		return "", fmt.Errorf("ioutil.ReadAll error: %w", err)
	}

	return string(b), nil
}
