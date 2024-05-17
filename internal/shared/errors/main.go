package errors

import "errors"

var ErrInternal = errors.New("INTERNAL_ERROR")

var ErrAuthFailure = errors.New("ACCOUNT_AUTH_FAILURE")
var ErrPreAuthFailure = errors.New("PRE_AUTH_FAILURE")

var ErrSpartanTokenGrabFailure = errors.New("SPARTAN_TOKEN_GRAB_FAILURE")
var ErrSpartanTokenInvalid = errors.New("SPARTAN_TOKEN_INVALID")

var ErrCredentialsReadFailure = errors.New("CREDENTIALS_READ_ERROR")
var ErrCredentialsWriteFailure = errors.New("CREDENTIALS_WRITE_ERROR")
var ErrCredentialsDirectoryCreateFailure = errors.New("CREDENTIALS_CREATE_DIRECTORY_ERROR")