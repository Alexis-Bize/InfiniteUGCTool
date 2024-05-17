package services_auth

import (
	"fmt"
	halowaypoint_request "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	msa_request "infinite-bookmarker/internal/shared/libs/msa/modules/request"
)

func GetAuthOptions() msa_request.LiveClientAuthOptions {
	return msa_request.LiveClientAuthOptions{
		ClientID: "000000004C0BD2F1",
		Scope: "xboxlive.signin xboxlive.offline_access",
		ResponseType: "code",
		RedirectURI: "https://www.halowaypoint.com/sign-in/callback",
		State: "/",
	}
}

func Authenticate(email string, password string) (string, error) {
	resp, err := msa_request.Authenticate(msa_request.LiveCredentials{
		Email: email,
		Password: password,
	}, GetAuthOptions())

	if err != nil {
		return "", err
	}

	location := resp.Header.Get("Location")

	if location == "" {
		return "", fmt.Errorf("something went wrong")
	}

	spartanToken, err := halowaypoint_request.ExtractSpartanTokenPostCallback(location)
	if err != nil {
		return "", err
	}
	
	return spartanToken, nil
}