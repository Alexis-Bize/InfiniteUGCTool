package prompt_svc

import (
	"fmt"
	"net/mail"
	"os"
	"strings"

	auth_svc "infinite-ugc-tool/internal/services/auth"
	identity_svc "infinite-ugc-tool/internal/services/identity"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/helpers/identity"

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

	email = os.Getenv("ACCOUNT_EMAIL")
	password = os.Getenv("ACCOUNT_PASSWORD")

	if email == "" || password == "" {
		email, password, err = requestIdentity(isRetry)
		if err != nil {
			return err
		}
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
