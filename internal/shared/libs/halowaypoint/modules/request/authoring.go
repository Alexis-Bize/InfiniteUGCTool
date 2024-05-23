package halowaypointRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/libs/halowaypoint"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"net/http"
)

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

	if resp.StatusCode == 401 {
		errors.Format("current STv4 is invalid or has expired", errors.ErrSpartanTokenInvalid)
	} else if resp.StatusCode == 403 {
		errors.Format("your are not allowed to bookmark this film", errors.ErrForbidden)
	} else if resp.StatusCode == 404 {
		errors.Format("desired film does not exist", errors.ErrNotFound)
	} else if resp.StatusCode != 200 {
		errors.Format("something went wrong", errors.ErrInternal)
	}

	return nil
}

