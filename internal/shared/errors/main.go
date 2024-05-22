package errors

import (
	"errors"
	"fmt"
)

var (
	ErrPrompt = errors.New("invalid prompt")
	ErrInternal = errors.New("internal error")
	ErrNotFound = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
	ErrAuthFailure = errors.New("account auth failure")
	ErrPreAuthFailure = errors.New("pre authentication failure")
	ErrSpartanTokenGrabFailure = errors.New("spartan token grab failure")
	ErrSpartanTokenInvalid = errors.New("invalid spartan token")
	ErrIdentityReadFailure = errors.New("identity read error")
	ErrIdentityWriteFailure = errors.New("identity write error")
	ErrIdentityMissing= errors.New("missing identity")
	ErrIdentityDirectoryCreateFailure = errors.New("indentity create directory failure")
)

func Format(message string, err error) error {
	return fmt.Errorf("%s: %s", err.Error(), message)
}
