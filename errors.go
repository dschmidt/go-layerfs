package layerfs

import (
	"fmt"
)

func newError(text string, name string) error {
	return fmt.Errorf("go-layerfs: %s: %s", text, name)
}

type layerFsError struct {
	name string
	text string
	err  error
}

func wrapError(err error, text string, name string) error {
	return &layerFsError{
		text: text,
		name: name,
		err:  err,
	}
}

func (l *layerFsError) Error() string {
	return fmt.Sprintf("go-layerfs: %s: %s", l.text, l.name)
}

func (l *layerFsError) Unwrap() error {
	return l.err
}
