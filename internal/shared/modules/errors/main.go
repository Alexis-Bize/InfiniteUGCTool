package errors

import "errors"

var ErrInternal = errors.New("INTERNAL_ERROR")
var ErrAuthFailure = errors.New("AUTH_FAILURE")
var ErrPreAuthFailure = errors.New("PREAUTH_HANDSHAKE_FAILURE")
var ErrSpartanTokenGrabFailure = errors.New("SPARTAN_TOKEN_GRAB_FAILURE")