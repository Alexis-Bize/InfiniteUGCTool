package errors

import "errors"

var (
	ErrInternal = errors.New("INTERNAL_ERROR")
	ErrAuthFailure = errors.New("ACCOUNT_AUTH_FAILURE")
	ErrPreAuthFailure = errors.New("PRE_AUTH_FAILURE")
	ErrSpartanTokenGrabFailure = errors.New("SPARTAN_TOKEN_GRAB_FAILURE")
	ErrSpartanTokenInvalid = errors.New("SPARTAN_TOKEN_INVALID")
	ErrCredentialsReadFailure = errors.New("CREDENTIALS_READ_ERROR")
	ErrCredentialsWriteFailure = errors.New("CREDENTIALS_WRITE_ERROR")
	ErrCredentialsDirectoryCreateFailure = errors.New("CREDENTIALS_CREATE_DIRECTORY_ERROR")
)