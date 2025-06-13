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

// Line can be with or without Docker headers
func ParseLogLine(line []byte) (timestamp, content []byte, err error) {
	if len(line) == 0 {
		return
	}

	// Try to find space after timestamp
	parts := bytes.SplitN(line, []byte{' '}, 2)
	if len(parts) < 2 {
		err = fmt.Errorf("Malformed log line")
		return
	}

	// Check if first part could be Docker header (8 bytes + space)
	if len(parts[0]) == 8 {
		// This was probably a header, take next parts
		parts = bytes.SplitN(parts[1], []byte{' '}, 2)
		if len(parts) < 2 {
			err = fmt.Errorf("Malformed log line")
			return
		}
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
