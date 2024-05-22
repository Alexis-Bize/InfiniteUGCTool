package errors

import (
	"errors"
	"fmt"
	"strings"
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
	ErrMatchIdInvalid = errors.New("invalid match id")
)

func Format(message string, err error) error {
	return fmt.Errorf("%s: %s", err.Error(), message)
}

func Is(current error, expected error) bool {
	return current == expected
}

func MayBe(current error, expected error) bool {
	parts := strings.Split(current.Error(), ":")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	
	if len(parts) > 0 {
		key := parts[0]
		return key == expected.Error()
	}

	return false
}