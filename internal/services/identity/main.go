package identityService

import (
	authService "infinite-bookmarker/internal/services/auth"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/libs/halowaypoint"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"os"
	"time"
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
		os.Stdout.WriteString("⏳ Refreshing your active session...\n")
		profile, spartanToken, err := authService.AuthenticateWithCredentials(currentIdentity.User.Email, currentIdentity.User.Password)
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
