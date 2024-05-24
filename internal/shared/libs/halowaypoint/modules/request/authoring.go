package halowaypointRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/libs/halowaypoint"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"io"
	"net/http"
)

func Bookmark(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Authoring, fmt.Sprintf("/hi/players/xuid(%s)/favorites/%s/%s", xuid, category, assetID))
	payload := map[string]interface{}{}

	if assetVersionID != "" {
		payload["AssetVersionId"] = assetVersionID
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Content-Type": "application/json",
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

func BookmarkFilmFromMatchID(xuid string, spartanToken string, matchID string) error {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Authoring, fmt.Sprintf("/hi/players/xuid(%s)/favorites/films/matches/%s", xuid, matchID))

	payload := map[string]interface{}{}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Content-Type": "application/json",
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

func CloneAsset(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	sessionID, err := createSession(xuid, spartanToken, category, assetID, assetVersionID)
	if err != nil {
		return err
	}

	return saveSession(xuid, spartanToken, category, assetID, sessionID)
}

func createSession(xuid string, spartanToken string, category string, assetID string, assetVersionID string) (string, error) {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Authoring, fmt.Sprintf("/hi/%s/new/sessions", category))
	payload := map[string]interface{}{
		"AssetToCopy": map[string]interface{}{
			"AssetId": assetID,
			"VersionId": assetVersionID,
		},
		"SessionOrigin": fmt.Sprintf("xuid(%s)", xuid),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Content-Type": "application/json",
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	err = OnResponse(resp)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	var session halowaypoint.NewSessionResponse
	if err := json.Unmarshal(body, &session); err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	return session.SessionID, nil
}

func saveSession(xuid string, spartanToken string, category string, assetID string, sessionID string) error {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Authoring, fmt.Sprintf("/hi/%s/%s/versions", category, assetID))
	payload := map[string]interface{}{
		"Source": "SaveAndEndSession",
		"SourceId": sessionID,
		"Player": fmt.Sprintf("xuid(%s)", xuid),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Content-Type": "application/json",
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