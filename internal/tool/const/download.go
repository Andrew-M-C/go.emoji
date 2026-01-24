package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func downloadEmoji() {
	// list and find latest emoji versions
	latest := emojiLatestVersion

	// check old emoji-data file
	data, err := httpGet(emojiDataFileURL())
	check(err)

	printf("latest version: %v", latest)

	// get emoji and zwj files
	basic, err := httpGet(emojiSeqFileURL(latest))
	check(err)
	zwj, err := httpGet(emojiZwjSeqFileURL(latest))
	check(err)

	err = os.WriteFile(emojiDataFile, []byte(data), 0644)
	check(err)
	err = os.WriteFile(emojiSeqFile, []byte(basic), 0644)
	check(err)
	err = os.WriteFile(emojiZwjSeqFile, []byte(zwj), 0644)
	check(err)

	referenceURL = dataFile(latest, "")
}

func httpGet(u string) (res string, err error) {
	printf("request: %s", u)

	resp, err := http.Get(u)
	if err != nil {
		return "", fmt.Errorf("http.Get(%s) error: %w", u, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP status %d, body '%s'", resp.StatusCode, string(b))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return string(b), nil
		}
		return "", fmt.Errorf("io.ReadAll error: %w", err)
	}

	return string(b), nil
}
