package auth_svc

import (
	"infinite-ugc-tool/internal/shared/libs/halowaypoint"
	halowaypoint_req "infinite-ugc-tool/internal/shared/libs/halowaypoint/modules/request"
	"infinite-ugc-tool/internal/shared/libs/msa"
	msa_req "infinite-ugc-tool/internal/shared/libs/msa/modules/request"
	"infinite-ugc-tool/internal/shared/modules/errors"
)

func AuthenticateWithCredentials(email string, password string) (halowaypoint.UserProfileResponse, string, error) {
	resp, err := msa_req.Authenticate(msa.LiveCredentials{
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

	spartanToken, err := halowaypoint_req.ExtractSpartanTokenPostCallback(location)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	profile, err := halowaypoint_req.GetUserProfile(spartanToken)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	return profile, spartanToken, nil
}
