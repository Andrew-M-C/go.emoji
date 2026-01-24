package main

import (
	"regexp"

	"github.com/Masterminds/semver/v3"
)

func listRemoteSubNodes(targetURI string) []string {
	raw, err := httpGet(targetURI)
	check(err)

	// find section like '<a href="8.0.0/">8.0.0/</a>'
	rep := regexp.MustCompile(`<a href="(\d+\.\d+\.\d+)/">(\d+\.\d+\.\d+)/</a>`)
	matches := rep.FindAllStringSubmatch(raw, -1)

	versions := make([]string, 0, len(matches))
	for _, match := range matches {
		printf("found version: %s", match[1])
		versions = append(versions, match[1])
	}

	return versions
}

func findLatestVersion(versions []string) string {
	if len(versions) == 0 {
		return ""
	}

	var latestVer *semver.Version
	for _, v := range versions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			// skip invalid version strings
			continue
		}
		if latestVer == nil || ver.GreaterThan(latestVer) {
			latestVer = ver
		}
	}

	if latestVer == nil {
		return ""
	}
	return latestVer.String()
}
