package msa_request

import (
	"Infinite-Bookmarker/internal/shared/modules/errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func Authenticate(credentials LiveCredentials, options LivePreAuthOptions) (*http.Response, error) {
	preAuthResponse, err := preAuth(&options)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("login", credentials.Email)
	form.Add("loginfmt", credentials.Email)
	form.Add("passwd", credentials.Password)
	form.Add("PPFT", preAuthResponse.Matches.PPFT)

	req, err := http.NewRequest("POST", preAuthResponse.Matches.URLPost, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	for k, v := range GetBaseHeaders(map[string]string{
		"Cookie": preAuthResponse.Cookie,
		"Content-Type": "application/x-www-form-urlencoded",
	}) { req.Header.Set(k, v) }

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil, fmt.Errorf("%w: %s", errors.ErrAuthFailure, "The authentication has failed")
	}

	return resp, nil
}

func getAuthorizeUrl(clientId string, scope string, responseType string, redirectUri string, state string) string {
	return fmt.Sprintf(
		"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s&state=%s&display=touch",
		url.QueryEscape(clientId),
		url.QueryEscape(scope),
		url.QueryEscape(responseType),
		url.QueryEscape(redirectUri),
		url.QueryEscape(state),
	)
}

func getMatchForIndex(body string, pattern *regexp.Regexp, index int) string {
	matches := pattern.FindStringSubmatch(body)
	if len(matches) > index {
		return matches[index]
	}

	return ""
}

func preAuth(options *LivePreAuthOptions) (*LivePreAuthResponse, error) {
	url := getAuthorizeUrl(
		options.ClientID,
		options.Scope,
		options.ResponseType,
		options.RedirectURI,
		options.State,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	for k, v := range GetBaseHeaders(map[string]string{}) {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	cookie := strings.Join(resp.Header["Set-Cookie"], "; ")
	ppftPattern := regexp.MustCompile(`sFTTag:'.*value="(.*)"\/>'`)
	urlPostPattern := regexp.MustCompile(`urlPost:'([^']+)'`)

	matches := LivePreAuthMatchedParameters{
		PPFT:    getMatchForIndex(string(body), ppftPattern, 1),
		URLPost: getMatchForIndex(string(body), urlPostPattern, 1),
	}

	if matches.PPFT == "" || matches.URLPost == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrPreAuthFailure, "please retry in a few seconds...")
	}

	return &LivePreAuthResponse{
		Cookie:  cookie,
		Matches: matches,
	}, nil
}
