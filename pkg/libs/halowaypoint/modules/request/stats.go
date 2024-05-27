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
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"infinite-ugc-tool/pkg/libs/halowaypoint"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
)

func GetMatchStats(spartanToken string, matchID string) (halowaypoint.MatchStatsResponse, error) {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Stats, fmt.Sprintf("/hi/matches/%s/stats", matchID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return halowaypoint.MatchStatsResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return halowaypoint.MatchStatsResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	err = OnResponse(resp)
	if err != nil {
		return halowaypoint.MatchStatsResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return halowaypoint.MatchStatsResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	var stats halowaypoint.MatchStatsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		return halowaypoint.MatchStatsResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	return stats, nil
}
