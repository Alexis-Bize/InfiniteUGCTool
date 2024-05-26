package halowaypoint_req

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"infinite-ugc-tool/internal/shared/libs/halowaypoint"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/utilities/request"
)

func GetMatchFilm(spartanToken string, matchID string) (halowaypoint.MatchSpectateResponse, error) {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Discovery, fmt.Sprintf("/hi/films/matches/%s/spectate", matchID))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return halowaypoint.MatchSpectateResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return halowaypoint.MatchSpectateResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	err = OnResponse(resp)
	if err != nil {
		return halowaypoint.MatchSpectateResponse{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return halowaypoint.MatchSpectateResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	var film halowaypoint.MatchSpectateResponse
	if err := json.Unmarshal(body, &film); err != nil {
		return halowaypoint.MatchSpectateResponse{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	return film, nil
}

func PingPublishedAsset(spartanToken string, category string, assetID string) error {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Discovery, fmt.Sprintf("/hi/%s/%s", category, assetID))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	return OnResponse(resp)
}