package trapi2ramlgen

import (
	"bufio"
	"bytes"
	"strings"
	"unicode"
)

func (w *wrapWriter) unidentText(text string) string {
	// determine the smallest number of spaces at the beginning of a line
	min_spaces := -1
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		ct := w.countLeadingSpace(scanner.Text())
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

func (w *wrapWriter) unidentTypedText(contenttype string, text string) string {
	// this changes key order, they should be kept in the same source order
	/*
		if contenttype == "json" || contenttype == "application/json" {
			var jdata interface{}
			err := json.Unmarshal([]byte(text), &jdata)
			if err == nil {
				ret, err := json.MarshalIndent(jdata, "", "    ")
				if err == nil {
					return string(ret)
				} else {
					w.warnings = append(w.warnings, NewErrWarning(fmt.Sprintf("Could not parse json: %v [%s]\n", err, text)))
				}
			} else {
				w.warnings = append(w.warnings, NewErrWarning(fmt.Sprintf("Could not parse json: %v [%s]\n", err, text)))
			}
		}
	*/
	// fallback
	return w.unidentText(text)
}

func (w *wrapWriter) countLeadingSpace(line string) int {
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
