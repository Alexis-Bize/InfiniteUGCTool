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

package auth_svc

import (
	"infinite-ugc-tool/pkg/libs/halowaypoint"
	halowaypoint_req "infinite-ugc-tool/pkg/libs/halowaypoint/modules/request"
	"infinite-ugc-tool/pkg/libs/msa"
	msa_req "infinite-ugc-tool/pkg/libs/msa/modules/request"
	"infinite-ugc-tool/pkg/modules/errors"
)

func AuthenticateWithCredentials(email string, password string) (halowaypoint.UserProfileResponse, string, error) {
	// 1. Ask halowaypoint where to authenticate. The returned URL points
	//    at login.microsoftonline.com / v2.0 with halowaypoint's own
	//    client_id, scope, and redirect_uri baked in — none of which we
	//    can safely hardcode.
	authorizeURL, err := halowaypoint_req.GetSignInAuthorizeURL()
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	// 2. Fetch the resulting authorize page. MS internally redirects this
	//    to the legacy login.live.com form, which is what we then scrape
	//    for PPFT / urlPost and which sets the session cookies we need on
	//    the credentials POST.
	page, err := msa_req.GetAuthorizePage(authorizeURL)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	// 3. POST credentials against the form discovered in the page (and
	//    drive the 2FA flow if MS asks for it).
	resp, err := msa_req.Authenticate(msa.LiveCredentials{
		Email:    email,
		Password: password,
	}, page)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	// 4. Follow MS's callback to halowaypoint to swap the auth code for
	//    a Spartan token.
	location := resp.Header.Get("Location")
	if location == "" {
		return halowaypoint.UserProfileResponse{}, "", errors.Format("something went wrong", errors.ErrInternal)
	}

	spartanToken, err := halowaypoint_req.ExtractSpartanTokenPostCallback(location)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	profile, err := halowaypoint_req.GetUserProfile(spartanToken)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

	return profile, spartanToken, nil
}
