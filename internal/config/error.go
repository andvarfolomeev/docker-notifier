package config

import "fmt"

type ErrMissingArg string

func (e ErrMissingArg) Error() string {
	return fmt.Sprintf("missing required argument: --%s", string(e))
}

var ErrHelpRequested = fmt.Errorf("help requested")
