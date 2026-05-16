// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompt_svc

import (
	"net/url"
	"regexp"
	"strings"

	identity_svc "infinite-ugc-tool/internal/application/services/auth/identity"
	"infinite-ugc-tool/internal/helpers/identity"
	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

const (
	BOOKMARK_FILES = "🔖 Bookmark files"
	CLONE_FILES = "🔄 Clone files"
	SHOW_CREDITS = "🤝 Show credits"
	SIGN_OUT = "🚫 Sign out"
	EXIT = "👋 Exit"
	GO_BACK = "← Go back"
)

const (
	FILM = "🎬 Match (Film)"
	MODE = "🎮 Mode (UgcGameVariant)"
	MAP = "🌎 Map (MapVariant)"
)

const (
	OPEN_X_1 = "Made by: Zeny IC"
	OPEN_X_2 = "Original idea: Okom"
	OPEN_X_3 = "Supporter: Grunt.API"
	OPEN_GITHUB = "Source code: GitHub"
)

const uuidPattern = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`

var uuidRegexp = regexp.MustCompile(uuidPattern)

func displayMatchGrabPrompt() (string, error) {
	var value string

	err := huh.NewInput().
		Title("Please specify a match ID or a valid match URL").
		Description("Leafapp.co, SpartanRecord.com, HaloDataHive.com and such are supported").
		Value(&value).
		Run()

	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	return extractMatchID(value)
}

func displayVariantDetailsPrompt() (string, string, error) {
	var assetID string
	var assetVersionID string
	var err error

	err = huh.NewInput().
		Title("Please specify a \"AssetID\" (GUID)").
		Description("e.g., ae4daed6-251a-4c2f-bc6f-eb25eac1bfd0").
		Value(&assetID).
		Validate(func (input string) error {
			_, err := extractUUID(input)
			if err != nil {
				return errors.New("invalid GUID")
			}

			return nil
		}).Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	err = huh.NewInput().
		Title("Please specify a \"AssetVariantID\" (GUID)").
		Description("This value is optional for published files").
		Value(&assetVersionID).
		Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	assetID = strings.TrimSpace(assetID)
	assetVersionID = strings.TrimSpace(assetVersionID)

	if assetVersionID != "" {
		_, err := extractUUID(assetVersionID)
		if err != nil {
			return "", "", err
		}
	}

	return strings.TrimSpace(assetID), strings.TrimSpace(assetVersionID), nil
}

func extractUUID(value string) (string, error) {
	match := uuidRegexp.FindString(strings.TrimSpace(value))
	if match != "" {
		return match, nil
	}

	return "", errors.Format("invalid format", errors.ErrUUIDInvalid)
}

// extractMatchID handles both raw GUIDs and URLs from common Halo tracker
// sites. For URLs we prefer the last GUID in the path because match URLs
// typically end with the match ID (e.g. /game/<guid>, /matches/<guid>).
func extractMatchID(value string) (string, error) {
	value = strings.TrimSpace(value)

	if u, err := url.Parse(value); err == nil && u.Host != "" {
		matches := uuidRegexp.FindAllString(u.Path, -1)
		if len(matches) > 0 {
			return matches[len(matches)-1], nil
		}
	}

	return extractUUID(value)
}

// runWithSpinnerAndRefresh runs fn inside a spinner. If fn fails with
// ErrSpartanTokenInvalid, the identity is refreshed *outside* the spinner —
// refresh may prompt for 2FA, and the OTC prompt would corrupt the terminal
// if it ran concurrently with a running spinner — and fn is retried once in
// a fresh spinner. fn must read the latest token via *currentIdentity each
// time it runs so the retry picks up the refreshed value.
func runWithSpinnerAndRefresh(currentIdentity *identity.Identity, title string, fn func() error) error {
	var err error
	spinner.New().Title(title).Action(func() {
		err = fn()
	}).Run()

	if err == nil || !errors.MayBe(err, errors.ErrSpartanTokenInvalid) {
		return err
	}

	refreshed, refreshErr := identity_svc.ForceRefresh()
	if refreshErr != nil {
		return err
	}
	*currentIdentity = refreshed

	spinner.New().Title(title).Action(func() {
		err = fn()
	}).Run()
	return err
}
