// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package request

import (
	"net/http"
	"strings"
	"time"
)

const requestTimeout = 60 * time.Second

// Client is the shared HTTP client used for normal requests. It follows
// redirects and times out after requestTimeout — without this, a hung remote
// (e.g. an unreachable Microsoft or Halo Waypoint endpoint) would block the
// CLI indefinitely.
var Client = &http.Client{
	Timeout: requestTimeout,
}

// NoRedirectClient is used when the caller needs to inspect a redirect's
// Location header or Set-Cookie response itself — the MSA login flow and
// Spartan-token extraction both depend on this.
var NoRedirectClient = &http.Client{
	Timeout: requestTimeout,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func GetBaseHeaders(extraHeaders map[string]string) map[string]string {
	headers := map[string]string{
		"User-Agent": RequestUserAgent,
		"Accept-Encoding": "identity",
	}

	for k, v := range extraHeaders {
		headers[k] = v
	}

	return headers
}

func ComputeUrl(baseUrl string, path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return baseUrl + path
}
