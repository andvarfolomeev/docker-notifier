package logfilter

import (
	"bytes"
	"fmt"
	"regexp"
)

type MatchedLine struct {
	Timestamp []byte
	Content   []byte
}

func FindMatchedLines(patterns []*regexp.Regexp, lines []byte) ([]*MatchedLine, error) {
	splitedLines := bytes.Split(lines, []byte{'\n'})
	res := make([]*MatchedLine, 0)

	for _, line := range splitedLines {
		timestamp, content, err := ParseLogLine(line)

		if err != nil {
			return nil, err
		}

		if IsMatchedLine(patterns, content) {
			res = append(res, &MatchedLine{
				Timestamp: timestamp,
				Content:   content,
			})
		}
	}

	return res, nil
}

// Line must be with headers. Read logs with timestamp option
func ParseLogLine(line []byte) (timestamp, content []byte, err error) {
	if len(line) == 0 {
		return
	}

	// strip Docker log header (first 8 bytes)
	if len(line) > 8 {
		line = line[8:]
	}

	parts := bytes.SplitN(line, []byte{' '}, 2)
	if len(parts) < 2 {
		err = fmt.Errorf("Malformed log line")
		return
	}

	timestamp = parts[0]
	content = parts[1]
	return
}

func IsMatchedLine(patterns []*regexp.Regexp, content []byte) bool {
	for _, pattern := range patterns {
		if pattern.Match(content) {
			return true
		}
	}
	return false
}
