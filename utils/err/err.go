package err

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

type Error struct {
	e       error
	code    int
	message string
}

// JSON returns the json encoded array of bytes
func (e *Error) JSON() ([]byte, error) {
	em := struct {
		Message string
		Code    int
		Error   string
	}{e.message, e.code, e.e.Error()}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(&em)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create json in err.Error.JSON()")
	}

	return buf.Bytes(), nil
}

func (e *Error) Code() int { return e.code }

func (e *Error) Error() error { return e.e }

// NewError returns a common error object used across Overpass
func NewError(
	err error,
	code int,
	message string,
) *Error {
	return &Error{err, code, message}
}
