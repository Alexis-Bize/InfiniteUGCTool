// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package halowaypoint_req

import (
	"net/http"
	"net/url"

	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
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
