package promptService

import (
	"fmt"
	authService "infinite-bookmarker/internal/services/auth"
	identityService "infinite-bookmarker/internal/services/identity"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"net/mail"
	"os"

	"github.com/manifoldco/promptui"
)

func StartAuthFlow() error {
	currentIdentity, err := identity.GetOrCreateIdentity(identity.Identity{})
	if err != nil {
		return err
	}

	if currentIdentity != (identity.Identity{}) {
		os.Stdout.WriteString(fmt.Sprintf("👋 Welcome back, %s!\n", currentIdentity.XboxNetwork.Gamertag))
		_, err := identityService.RefreshIdentityIfRequired(currentIdentity)
		return err
	}

	email, password, err := requestIdentity()
	if err != nil {
		return err
	}

	os.Stdout.WriteString("⏳ Authenticating...\n")
	profile, spartanToken, err := authService.AuthenticateWithCredentials(email, password)
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("✅ Welcome %s!\n", profile.Gamertag))
	identityService.StoreIdentity(email, password, profile, spartanToken)
	return nil
}

func requestIdentity() (string, string, error) {
	os.Stdout.WriteString("👋 Hey there! Please authenticate using your Microsoft credentials to continue.\n")

	prompt := promptui.Prompt{
		Label: "📧 Email address",
		Validate: func(input string) error {
			_, err := mail.ParseAddress(input)
			if err != nil {
				return errors.Format("specified email address is invalid", errors.ErrPrompt)
			}

			return nil
		},
	}

	email, err := prompt.Run()
	if err != nil {
		return "", "", err
	}

	prompt = promptui.Prompt{
		Label: "🔑 Password",
		Mask: '*',
		Validate: func(input string) error {
			if len(input) == 0 {
				return errors.Format("password can not be empty", errors.ErrPrompt)
			}

			return nil
		},
	}

	password, err := prompt.Run()
	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	return email, password, nil
}
