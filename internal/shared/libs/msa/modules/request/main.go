package msa_request

import (
	"fmt"
	"net/url"
)

func BuildAuthorizeUrl(clientId string, scope string, responseType string, redirectUri string, state string) string {
	return fmt.Sprintf(
		"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s&state=%s&display=touch",
		url.QueryEscape(clientId),
		url.QueryEscape(scope),
		url.QueryEscape(responseType),
		url.QueryEscape(redirectUri),
		url.QueryEscape(state),
	)
}