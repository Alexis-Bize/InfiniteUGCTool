package halowaypoint_request

import (
	"Infinite-Bookmarker/internal/shared/modules/errors"
	"fmt"
	"net/http"
	"net/url"
)

func ExtractSpartanTokenPostCallback(location string) (string, error) {
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	for k, v := range GetBaseHeaders(map[string]string{
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
				return "", fmt.Errorf("%w: %s", errors.ErrSpartanTokenGrabFailure, "Please retry in a few seconds...")
			}

			break
		}
	}

	if tokenValue == "" {
		return "", fmt.Errorf("%w: %s", errors.ErrSpartanTokenGrabFailure, "Please retry in a few seconds...")
	}

	return tokenValue, nil
}