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
	"encoding/json"
	"io"
	"net/http"

	"infinite-ugc-tool/pkg/libs/halowaypoint"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
)

func GetUserProfile(spartanToken string) (halowaypoint.UserProfileResponse, error) {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Profile, "/users/me")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	err = OnResponse(resp)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return halowaypoint.UserProfileResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	var user halowaypoint.UserProfileResponse

	if err := json.Unmarshal(body, &user); err != nil {
		return halowaypoint.UserProfileResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	return user, nil
}
