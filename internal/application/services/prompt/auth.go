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

package prompt_svc

import (
	"fmt"
	"net/mail"
	"os"
	"strings"

	auth_svc "infinite-ugc-tool/internal/application/services/auth"
	identity_svc "infinite-ugc-tool/internal/application/services/auth/identity"
	"infinite-ugc-tool/internal/helpers/identity"
	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func StartAuthFlow(isRetry bool) error {
	var err error
	var email string
	var password string

	currentIdentity, _ := identity.GetOrCreateIdentity(identity.Identity{})
	if currentIdentity != (identity.Identity{}) {
		os.Stdout.WriteString(fmt.Sprintf("👋 Welcome back, %s!\n", currentIdentity.XboxNetwork.Gamertag))
		_, err := identity_svc.RefreshIdentityIfRequired(currentIdentity)
		return err
	}

	email, password, err = requestIdentity(isRetry)
	if err != nil {
		return err
	}

	spinner.New().Title("Authenticating...").Run()
	profile, spartanToken, err := auth_svc.AuthenticateWithCredentials(email, password)
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("✅ Welcome %s!\n", profile.Gamertag))
	identity_svc.StoreIdentity(email, password, profile, spartanToken)
	return nil
}

func requestIdentity(isRetry bool) (string, string, error) {
	var err error
	var email string
	var password string

	if !isRetry {
		os.Stdout.WriteString("| Hey there! Please authenticate using your Microsoft credentials to continue\n")
		os.Stdout.WriteString("└── You must have authenticated on HaloWaypoint.com at least once before!\n")
	}

	err = huh.NewInput().
		Title("What's your email?").
		Value(&email).
		Validate(func (input string) error {
			_, err := mail.ParseAddress(input)
			if err != nil {
				return errors.New("specified email is invalid")
			}

			return nil
		}).Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	err = huh.NewInput().
		Title("What's your password?").
		Password(true).
		Value(&password).
		Validate(func (input string) error {
			if len(input) == 0 {
				return errors.New("password can not be empty")
			}

			return nil
		}).Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	return strings.TrimSpace(email), strings.TrimSpace(password), nil
}
