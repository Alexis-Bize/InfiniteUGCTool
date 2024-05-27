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
	resp, err := msa_req.Authenticate(msa.LiveCredentials{
		Email: email,
		Password: password,
	}, msa.LiveClientAuthOptions{
		ClientID: "000000004C0BD2F1",
		Scope: "xboxlive.signin xboxlive.offline_access",
		ResponseType: "code",
		RedirectURI: "https://www.halowaypoint.com/sign-in/callback",
		State: "/",
	})

	if err != nil {
		return halowaypoint.UserProfileResponse{}, "", err
	}

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
