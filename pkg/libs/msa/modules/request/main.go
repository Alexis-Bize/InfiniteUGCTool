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

package msa_req

import (
	"fmt"
	"net/url"
)

func BuildAuthorizeUrl(clientId string, scope string, responseType string, redirectUri string, state string) string {
	return fmt.Sprintf(
		"https://login.live.com/oauth20_authorize.srf?client_id=%s&scope=%s&response_type=%s&redirect_uri=%s&state=%s&display=touch",
		url.QueryEscape(clientId),
		url.QueryEscape(scope),
		url.QueryEscape(responseType),
		url.QueryEscape(redirectUri),
		url.QueryEscape(state),
	)
}
