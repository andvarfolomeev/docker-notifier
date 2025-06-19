package docker

import (
	"encoding/json"
	"fmt"
)

type Filters struct {
	values map[string][]string
}

func NewFilter() *Filters {
	return &Filters{
		values: make(map[string][]string),
	}
}

func (f *Filters) Add(key, value string) *Filters {
	f.values[key] = append(f.values[key], value)
	return f
}

func (f *Filters) Encode() (string, error) {
	b, err := json.Marshal(f.values)
	if err != nil {
		return "", fmt.Errorf("encode filters: %w", err)
	}
	return string(b), nil
}

func (f *Filters) Len() int {
	return len(f.values)
}
