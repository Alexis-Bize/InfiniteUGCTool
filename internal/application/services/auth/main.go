package auth_service

import (
	halowaypoint_request "Infinite-Bookmarker/internal/shared/libs/halowaypoint/modules/request"
	msa_request "Infinite-Bookmarker/internal/shared/libs/msa/modules/request"
	"fmt"
)

func Authenticate(email string, password string) (string, error) {
	resp, err := msa_request.Authenticate(msa_request.LiveCredentials{
		Email: email,
		Password: password,
	}, msa_request.LivePreAuthOptions{
		ClientID: "000000004C0BD2F1",
		Scope: "xboxlive.signin xboxlive.offline_access",
		ResponseType: "code",
		RedirectURI: "https://www.halowaypoint.com/sign-in/callback",
		State: "/",
	})

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