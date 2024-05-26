package halowaypoint_req

import (
	"net/http"
	"net/url"

	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/utilities/request"
)

func ExtractSpartanTokenPostCallback(location string) (string, error) {
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
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
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	tokenName := "343-spartan-token"
	var tokenValue string

	for _, cookie := range cookies {
		if cookie.Name == tokenName {
			tokenValue, err = url.QueryUnescape(cookie.Value)
			if err != nil {
				return "", errors.Format("please retry in a few seconds", errors.ErrSpartanTokenGrabFailure)
			}

			break
		}
	}

	if tokenValue == "" {
		return "", errors.Format("please retry in a few seconds", errors.ErrSpartanTokenGrabFailure)
	}

	return tokenValue, nil
}
