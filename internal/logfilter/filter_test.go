package logfilter_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/andvarfolomeev/docker-notifier/internal/logfilter"
)

func withDockerHeader(line string) []byte {
	return append([]byte{1, 0, 0, 0, 0, 0, 0, byte(len(line))}, []byte(line)...)
}

func TestFindMatchedLines(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		input    []byte
		want     []*logfilter.MatchedLine
		wantErr  bool
	}{
		{
			name:     "empty input",
			patterns: []string{"error"},
			input:    []byte{},
			want:     []*logfilter.MatchedLine{},
			wantErr:  false,
		},
		{
			name:     "single matching line with header",
			patterns: []string{"error"},
			input:    withDockerHeader("2023-11-15T10:00:00Z Some error occurred"),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:00Z"),
					Content:   []byte("Some error occurred"),
				},
			},
			wantErr: false,
		},
		{
			name:     "single matching line without header",
			patterns: []string{"error"},
			input:    []byte("2023-11-15T10:00:00Z Some error occurred"),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:00Z"),
					Content:   []byte("Some error occurred"),
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple lines with one match with header",
			patterns: []string{"error"},
			input: bytes.Join([][]byte{
				withDockerHeader("2023-11-15T10:00:00Z Normal log line"),
				withDockerHeader("2023-11-15T10:00:01Z Some error occurred"),
				withDockerHeader("2023-11-15T10:00:02Z Another normal line"),
			}, []byte("\n")),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:01Z"),
					Content:   []byte("Some error occurred"),
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple lines with one match without header",
			patterns: []string{"error"},
			input: []byte(
				"2023-11-15T10:00:00Z Normal log line\n" +
					"2023-11-15T10:00:01Z Some error occurred\n" +
					"2023-11-15T10:00:02Z Another normal line",
			),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:01Z"),
					Content:   []byte("Some error occurred"),
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple patterns with header",
			patterns: []string{"error", "warning"},
			input: bytes.Join([][]byte{
				withDockerHeader("2023-11-15T10:00:00Z Some error occurred"),
				withDockerHeader("2023-11-15T10:00:01Z Warning message"),
				withDockerHeader("2023-11-15T10:00:02Z Normal line"),
			}, []byte("\n")),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:00Z"),
					Content:   []byte("Some error occurred"),
				},
				{
					Timestamp: []byte("2023-11-15T10:00:01Z"),
					Content:   []byte("Warning message"),
				},
			},
			wantErr: false,
		},
		{
			name:     "multiple patterns without header",
			patterns: []string{"error", "warning"},
			input: []byte(
				"2023-11-15T10:00:00Z Some error occurred\n" +
					"2023-11-15T10:00:01Z Warning message\n" +
					"2023-11-15T10:00:02Z Normal line",
			),
			want: []*logfilter.MatchedLine{
				{
					Timestamp: []byte("2023-11-15T10:00:00Z"),
					Content:   []byte("Some error occurred"),
				},
				{
					Timestamp: []byte("2023-11-15T10:00:01Z"),
					Content:   []byte("Warning message"),
				},
			},
			wantErr: false,
		},
		{
			name:     "malformed line",
			patterns: []string{"error"},
			input:    []byte("malformed"),
			want:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := make([]*regexp.Regexp, len(tt.patterns))
			for i, p := range tt.patterns {
				patterns[i] = regexp.MustCompile("(?i)" + p)
			}

			got, err := logfilter.FindMatchedLines(patterns, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMatchedLines() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("FindMatchedLines() got %d matches, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if !bytes.Equal(got[i].Timestamp, tt.want[i].Timestamp) {
					t.Errorf("FindMatchedLines() timestamp[%d] = %s, want %s", i, got[i].Timestamp, tt.want[i].Timestamp)
				}
				if !bytes.Equal(got[i].Content, tt.want[i].Content) {
					t.Errorf("FindMatchedLines() content[%d] = %s, want %s", i, got[i].Content, tt.want[i].Content)
				}
			}
		})
	}
}
