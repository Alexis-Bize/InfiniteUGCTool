package halowaypoint_req

import (
	"encoding/json"
	"io"
	"net/http"

	"infinite-ugc-tool/internal/shared/libs/halowaypoint"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/utilities/request"
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
