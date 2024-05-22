package errors

import "errors"

var (
	ErrInternal = errors.New("INTERNAL_ERROR")
	ErrAuthFailure = errors.New("ACCOUNT_AUTH_FAILURE")
	ErrPreAuthFailure = errors.New("PRE_AUTH_FAILURE")
	ErrSpartanTokenGrabFailure = errors.New("SPARTAN_TOKEN_GRAB_FAILURE")
	ErrSpartanTokenInvalid = errors.New("SPARTAN_TOKEN_INVALID")
	ErrIdentityReadFailure = errors.New("IDENTITY_READ_ERROR")
	ErrIdentityWriteFailure = errors.New("IDENTITY_WRITE_ERROR")
	ErrIdentityDirectoryCreateFailure = errors.New("IDENTITY_CREATE_DIRECTORY_ERROR")
)