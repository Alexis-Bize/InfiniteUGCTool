package prompt

import (
	"fmt"
	"infinite-bookmarker/internal/application/auth"
	"infinite-bookmarker/internal/shared/errors"
	halowaypointRequest "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"net/mail"
	"os"
	"time"

	"github.com/manifoldco/promptui"
)

func StartAuthFlow() error {
	currentIdentity, err := identity.GetOrCreateIdentity(identity.Identity{})
	if err != nil {
		return err
	}

	if currentIdentity != (identity.Identity{}) {
		os.Stdout.WriteString(fmt.Sprintf("👋 Welcome back, %s!\n", currentIdentity.XboxNetwork.Gamertag))
		shouldRefreshCredentials := true

		if (currentIdentity.SpartanToken.Expiration != "") {
			parsedTime, err := time.Parse(time.RFC3339, currentIdentity.SpartanToken.Expiration)
			if err == nil && time.Now().Before(parsedTime) {
				shouldRefreshCredentials = false
			}
		}

		if shouldRefreshCredentials {
			os.Stdout.WriteString("⏳ Refreshing your active session...\n")
			profile, spartanToken, err := auth.AuthenticateWithCredentials(currentIdentity.User.Email, currentIdentity.User.Password)
			if err != nil {
				return err
			}

			os.Stdout.WriteString("✅ Your active session has been refreshed with success!\n")
			return storeIdentity(currentIdentity.User.Email, currentIdentity.User.Password, profile, spartanToken)
		}

		return nil
	}

	email, password, err := requestIdentity()
	if err != nil {
		return err
	}

	os.Stdout.WriteString("⏳ Authenticating...\n")
	profile, spartanToken, err := auth.AuthenticateWithCredentials(email, password)
	if err != nil {
		return err
	}

	os.Stdout.WriteString(fmt.Sprintf("✅ Welcome %s!\n", profile.Gamertag))
	return storeIdentity(email, password, profile, spartanToken)
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

func storeIdentity(email string, password string, profile halowaypointRequest.UserProfileResponse, spartanToken string) error {
	now := time.Now()
	spartanTokenEstimatedExpiration := now.Add(3 * time.Hour)

	err := identity.SaveIdentity(identity.Identity{
		User: identity.UserCredentials{
			Email: email,
			Password: password,
		},
		SpartanToken: identity.SpartanTokenDetails{
			Value: spartanToken,
			Expiration: spartanTokenEstimatedExpiration.Format(time.RFC3339),
		},
		XboxNetwork: identity.XboxNetworkIdentity{
			Xuid: profile.Xuid,
			Gamertag: profile.Gamertag,
		},
	})

	return err
}
