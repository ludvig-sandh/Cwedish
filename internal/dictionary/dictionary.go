package dictionary

import (
	"fmt"
	"os"
	"strings"
)

type Dictionary map[string]string

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.Split(s, "\n")
}

func parseDictionaryLine(line string) (original string, swedish string, sucess bool) {
	if len(line) == 0 {
		return "", "", false
	}

	words := strings.Split(line, ":")
	if len(words) != 2 {
		fmt.Println("Couldn't parse line: ", line)
		return "", "", false
	}

	return words[0], words[1], true
}

func ParseDictionaryFile(path string) Dictionary {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	content := string(bytes[:])
	dictionary := make(Dictionary)

	lines := splitLines(content)
	for _, line := range lines {
		original, swedish, success := parseDictionaryLine(line)
		if !success {
			continue
		}

		dictionary[swedish] = original
	}

	return dictionary
}
