package authService

import (
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/libs/halowaypoint"
	halowaypointRequest "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	"infinite-bookmarker/internal/shared/libs/msa"
	msaRequest "infinite-bookmarker/internal/shared/libs/msa/modules/request"
)

func AuthenticateWithCredentials(email string, password string) (halowaypoint.UserProfileResponse, string, error) {
	resp, err := msaRequest.Authenticate(msa.LiveCredentials{
		Email: email,
		Password: password,
	}, msa.LiveClientAuthOptions{
		ClientID: "000000004C0BD2F1",
		Scope: "xboxlive.signin xboxlive.offline_access",
		ResponseType: "code",
		RedirectURI: "https://www.halowaypoint.com/sign-in/callback",
		State: "/",
	})

	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	location := resp.Header.Get("Location")

	if location == "" {
		return halowaypoint.UserProfileResponse{}, "", errors.Format("something went wrong", errors.ErrInternal)
	}

	spartanToken, err := halowaypointRequest.ExtractSpartanTokenPostCallback(location)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	profile, err := halowaypointRequest.GetUserProfile(spartanToken)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	return profile, spartanToken, nil
}
