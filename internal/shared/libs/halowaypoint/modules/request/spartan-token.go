package halowaypointRequest

import (
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"net/http"
	"net/url"
)

func ExtractSpartanTokenPostCallback(location string) (string, error) {
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "*/*",
	}) { req.Header.Set(k, v) }

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	tokenName := "343-spartan-token"
	var tokenValue string

	for _, cookie := range cookies {
		if cookie.Name == tokenName {
			tokenValue, err = url.QueryUnescape(cookie.Value)
			if err != nil {
				return "", fmt.Errorf("%w: %s", errors.ErrSpartanTokenGrabFailure, "please retry in a few seconds...")
			}

			break
		}
	}

	if tokenValue == "" {
		return "", fmt.Errorf("%w: %s", errors.ErrSpartanTokenGrabFailure, "please retry in a few seconds...")
	}

	return tokenValue, nil
}