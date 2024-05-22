package auth

import (
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	halowaypointRequest "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	msaRequest "infinite-bookmarker/internal/shared/libs/msa/modules/request"
)

func GetAuthOptions() msaRequest.LiveClientAuthOptions {
	return msaRequest.LiveClientAuthOptions{
		ClientID: "000000004C0BD2F1",
		Scope: "xboxlive.signin xboxlive.offline_access",
		ResponseType: "code",
		RedirectURI: "https://www.halowaypoint.com/sign-in/callback",
		State: "/",
	}
}

func AuthenticateWithCredentials(email string, password string) (halowaypointRequest.UserProfileResponse, string, error) {
	resp, err := msaRequest.Authenticate(msaRequest.LiveCredentials{
		Email: email,
		Password: password,
	}, GetAuthOptions())

	if err != nil {
		return halowaypointRequest.UserProfileResponse{}, "", err
	}

	location := resp.Header.Get("Location")

	if location == "" {
		return halowaypointRequest.UserProfileResponse{}, "", fmt.Errorf("%w: %s", errors.ErrInternal, "something went wrong")
	}

	spartanToken, err := halowaypointRequest.ExtractSpartanTokenPostCallback(location)
	if err != nil {
		return halowaypointRequest.UserProfileResponse{}, "",  err
	}

	profile, err := halowaypointRequest.GetUserProfile(spartanToken)
	if err != nil {
		return halowaypointRequest.UserProfileResponse{}, "", err
	}
	
	return profile, spartanToken, nil
}