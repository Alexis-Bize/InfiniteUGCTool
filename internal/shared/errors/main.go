package errors

import (
	"errors"
	"fmt"
)

var (
	ErrPrompt = errors.New("invalid prompt")
	ErrInternal = errors.New("internal error")
	ErrAuthFailure = errors.New("account auth failure")
	ErrPreAuthFailure = errors.New("pre authentication failure")
	ErrSpartanTokenGrabFailure = errors.New("spartan token grab failure")
	ErrSpartanTokenInvalid = errors.New("spartan token invalid")
	ErrIdentityReadFailure = errors.New("identity read error")
	ErrIdentityWriteFailure = errors.New("identity write error")
	ErrIdentityDirectoryCreateFailure = errors.New("indentity create directory failure")
)

func Format(message string, err error) error {
	return fmt.Errorf("%s: %s", err.Error(), message)
}
