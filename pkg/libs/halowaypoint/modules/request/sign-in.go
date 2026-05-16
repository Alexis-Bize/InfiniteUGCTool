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

package halowaypoint_req

import (
	"net/http"

	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
)

const signInDiscoveryURL = "https://www.halowaypoint.com/sign-in?path=/"

// GetSignInAuthorizeURL hits halowaypoint's /sign-in endpoint without
// following the 302 and returns the upstream identity-provider URL it
// points at (currently login.microsoftonline.com / v2.0). The concrete
// client_id, scope, redirect_uri and other parameters are chosen by
// halowaypoint at request time — letting it drive avoids the
// "unauthorized_client" failures that come from hardcoding a client_id
// MS no longer accepts in the public shape.
func GetSignInAuthorizeURL() (string, error) {
	req, err := http.NewRequest("GET", signInDiscoveryURL, nil)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "*/*",
	}) { req.Header.Set(k, v) }

	resp, err := request.NoRedirectClient.Do(req)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	location := resp.Header.Get("Location")
	if location == "" {
		return "", errors.Format("no Location header on /sign-in response", errors.ErrInternal)
	}

	return location, nil
}
