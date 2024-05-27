// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	ErrAuth2FARequired= errors.New("two factor authentication required")
	ErrPreAuthFailure = errors.New("pre authentication failure")
	ErrSpartanTokenGrabFailure = errors.New("spartan token grab failure")
	ErrSpartanTokenInvalid = errors.New("invalid spartan token")
	ErrIdentityReadFailure = errors.New("identity read error")
	ErrIdentityWriteFailure = errors.New("identity write error")
	ErrIdentityMissing= errors.New("missing identity")
	ErrIdentityDirectoryCreateFailure = errors.New("indentity create directory failure")
	ErrUUIDInvalid = errors.New("invalid UUID")
)

func New(message string) error {
	return errors.New(message)
}

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
