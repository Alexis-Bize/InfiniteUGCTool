package halowaypointRequest

import (
	"infinite-bookmarker/internal/shared/modules/errors"
	"net/http"
)

func OnResponse(resp *http.Response) error {
	var err error

	if resp.StatusCode == 401 {
		err = errors.Format("current STv4 is invalid or has expired", errors.ErrSpartanTokenInvalid)
	} else if resp.StatusCode == 403 {
		err = errors.Format("your are not allowed to perform this action", errors.ErrForbidden)
	} else if resp.StatusCode == 404 {
		err = errors.Format("desired entity does not exist", errors.ErrNotFound)
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = errors.Format("something went wrong", errors.ErrInternal)
	}

	return err
}
