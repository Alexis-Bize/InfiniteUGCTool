// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msa_req

import (
	"io"
	"net/http"

	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
)

// AuthorizePageResponse is the (followed-redirects) result of GETting the
// upstream identity-provider authorize URL handed back by halowaypoint.
// The legacy login.live.com form is reached via internal MS redirects, so
// the caller never has to know which host it ultimately landed on — just
// hand the response + body to Authenticate, which will scrape PPFT and
// urlPost out of it.
type AuthorizePageResponse struct {
	Response *http.Response
	Body     []byte
}

// GetAuthorizePage performs a GET on the authorize URL with redirects
// followed and returns the final response together with the rendered
// login page body.
func GetAuthorizePage(authorizeURL string) (*AuthorizePageResponse, error) {
	req, err := http.NewRequest("GET", authorizeURL, nil)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "*/*",
	}) { req.Header.Set(k, v) }

	resp, err := request.Client.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	return &AuthorizePageResponse{Response: resp, Body: body}, nil
}
