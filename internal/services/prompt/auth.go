package promptService

import (
	"fmt"
	authService "infinite-bookmarker/internal/services/auth"
	identityService "infinite-bookmarker/internal/services/identity"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"net/mail"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func StartAuthFlow(isRetry bool) error {
	currentIdentity, err := identity.GetOrCreateIdentity(identity.Identity{})
	if err != nil {
		return err
	}

	if currentIdentity != (identity.Identity{}) {
		os.Stdout.WriteString(fmt.Sprintf("👋 Welcome back, %s!\n", currentIdentity.XboxNetwork.Gamertag))
		_, err := identityService.RefreshIdentityIfRequired(currentIdentity)
		return err
	}

	email, password, err := requestIdentity(isRetry)
	if err != nil {
		return err
	}

	spinner.New().Title("Authenticating...").Run()
	
	profile, spartanToken, err := authService.AuthenticateWithCredentials(email, password)
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("✅ Welcome %s!\n", profile.Gamertag))
	identityService.StoreIdentity(email, password, profile, spartanToken)
	return nil
}

func requestIdentity(isRetry bool) (string, string, error) {
	var err error
	var email string
	var password string

	if !isRetry {
		os.Stdout.WriteString("👋 Hey there! Please authenticate using your Microsoft credentials to continue.\n")
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

	return email, password, nil
}
