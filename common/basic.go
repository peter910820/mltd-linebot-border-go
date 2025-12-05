package common

import "errors"

var (
	// The remote server returns a non-200 response status code
	ErrStatusCodeAbnormal = errors.New("response: server returned an error status code")
)
