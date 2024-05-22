package msaRequest

import (
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

func Authenticate(credentials LiveCredentials, options LiveClientAuthOptions) (*http.Response, error) {
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

	for k, v := range request.GetBaseHeaders(map[string]string{
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
		return nil, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil, fmt.Errorf("%w: %s", errors.ErrAuthFailure, "the authentication has failed")
	}

	return resp, nil
}

func getMatchForIndex(body string, pattern *regexp.Regexp, index int) string {
	matches := pattern.FindStringSubmatch(body)
	if len(matches) > index {
		return matches[index]
	}

	return ""
}

func preAuth(options *LiveClientAuthOptions) (*LivePreAuthResponse, error) {
	url := BuildAuthorizeUrl(
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

	for k, v := range request.GetBaseHeaders(map[string]string{}) {
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
	urlPostPattern := regexp.MustCompile(`urlPost:'([^']+)'`)
	ppftPattern := regexp.MustCompile(`sFTTag:'.*value="(.*)"\/>'`)

	matches := LivePreAuthMatchedParameters{
		URLPost: getMatchForIndex(string(body), urlPostPattern, 1),
		PPFT: getMatchForIndex(string(body), ppftPattern, 1),
	}

	if matches.PPFT == "" || matches.URLPost == "" {
		return nil, fmt.Errorf("%w: %s", errors.ErrPreAuthFailure, "please retry in a few seconds...")
	}

	return &LivePreAuthResponse{
		Cookie:  cookie,
		Matches: matches,
	}, nil
}
