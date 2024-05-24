package halowaypointRequest

import (
	"encoding/json"
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/libs/halowaypoint"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"io"
	"net/http"
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
