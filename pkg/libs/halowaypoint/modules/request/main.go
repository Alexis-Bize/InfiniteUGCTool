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
	"net/http"

	"infinite-ugc-tool/pkg/modules/errors"
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
