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

func GetUserProfile(spartanToken string) (UserProfileResponse, error) {
	url := request.ComputeUrl(halowaypoint.GetConfig().Urls.Profile, "/users/me")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return UserProfileResponse{}, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "application/json",
		"X-343-Authorization-Spartan": spartanToken,
	}) { req.Header.Set(k, v) }

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return UserProfileResponse{}, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return UserProfileResponse{}, fmt.Errorf("%w: %s", errors.ErrSpartanTokenInvalid, "current STv4 is invalid or has expired")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserProfileResponse{}, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	var user UserProfileResponse

	if err := json.Unmarshal(body, &user); err != nil {
		return UserProfileResponse{}, fmt.Errorf("%w: %s", errors.ErrInternal, err.Error())
	}

	return user, nil
}