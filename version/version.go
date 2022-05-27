package version

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

//go:embed version.txt
var version string

func getParsedVersion() (string, string, string) {
	split := strings.Split(version, ".")
	if len(split) != 3 {
		panic(errors.New("failed to read version from version.txt"))
	}

	return split[0], split[1], split[2]
}

func MajorMinor() string {
	major, minor, _ := getParsedVersion()
	return fmt.Sprintf("%s.%s", major, minor)
}

func Full() string {
	return version
}

func Compat(client string, server string) bool {
	return client == server
}
