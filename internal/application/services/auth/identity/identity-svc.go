package identity_svc

import (
	"os"
	"time"

	auth_svc "infinite-ugc-tool/internal/application/services/auth"
	"infinite-ugc-tool/internal/helpers/identity"
	"infinite-ugc-tool/pkg/libs/halowaypoint"
	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh/spinner"
)

func GetActiveIdentity() (identity.Identity, error) {
	currentIdentity, err := identity.GetOrCreateIdentity(identity.Identity{})
	if err != nil {
		return identity.Identity{}, err
	} else if (currentIdentity == (identity.Identity{})) {
		return identity.Identity{}, errors.Format("empty identity", errors.ErrIdentityMissing)
	}

	currentIdentity, err = RefreshIdentityIfRequired(currentIdentity)
	if err != nil {
		return identity.Identity{}, err
	}

	return currentIdentity, nil
}

func RefreshIdentityIfRequired(currentIdentity identity.Identity) (identity.Identity, error) {
	shouldRefresh := true

	if (currentIdentity.SpartanToken.Expiration != "") {
		parsedTime, err := time.Parse(time.RFC3339, currentIdentity.SpartanToken.Expiration)
		if err == nil && time.Now().Before(parsedTime) {
			shouldRefresh = false
		}
	}

	if shouldRefresh {
		spinner.New().Title("Refreshing your active session...").Run()
		profile, spartanToken, err := auth_svc.AuthenticateWithCredentials(currentIdentity.User.Email, currentIdentity.User.Password)
		if err != nil {
			identity.SaveIdentity(identity.Identity{})
			return identity.Identity{}, err
		}

		os.Stdout.WriteString("✅ Your active session has been refreshed with success!\n")
		return StoreIdentity(currentIdentity.User.Email, currentIdentity.User.Password, profile, spartanToken)
	}

	return currentIdentity, nil
}

func StoreIdentity(email string, password string, profile halowaypoint.UserProfileResponse, spartanToken string) (identity.Identity, error) {
	now := time.Now()
	spartanTokenEstimatedExpiration := now.Add(3 * time.Hour)
	savedIdentity := identity.Identity{
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
	}

	err := identity.SaveIdentity(savedIdentity)
	if err != nil {
		return identity.Identity{}, err
	}

	return savedIdentity, err
}
