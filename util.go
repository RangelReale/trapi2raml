package trapi2raml

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

func unidentText(text string) string {
	// determine the smallest number of spaces at the beginning of a line
	min_spaces := -1
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		ct := countLeadingSpace(scanner.Text())
		if min_spaces == -1 || ct < min_spaces {
			min_spaces = ct
		}
	}

	if min_spaces <= 0 {
		return text
	}

	var ret bytes.Buffer
	scanner2 := bufio.NewScanner(strings.NewReader(text))
	is_first := true
	for scanner2.Scan() {
		if !is_first {
			ret.WriteString("\n")
		}
		is_first = false
		line := scanner2.Text()
		ret.WriteString(line[min_spaces:])
	}

	return ret.String()
}

func countLeadingSpace(line string) int {
	i := 0
	for _, runeValue := range line {
		if unicode.IsSpace(runeValue) {
			i++
		} else {
			break
		}
	}
	return i
}
