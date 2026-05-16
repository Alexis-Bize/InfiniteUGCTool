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

package halowaypoint_req

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"infinite-ugc-tool/configs"
	"infinite-ugc-tool/pkg/modules/debug"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"
)

func ExtractSpartanTokenPostCallback(location string) (string, error) {
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Accept": "*/*",
	}) { req.Header.Set(k, v) }

	resp, err := request.NoRedirectClient.Do(req)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	tokenName := "343-spartan-token"
	var tokenValue string

	for _, cookie := range cookies {
		if cookie.Name == tokenName {
			tokenValue, err = url.QueryUnescape(cookie.Value)
			if err != nil {
				dumpSpartanTokenDebug(location, resp, cookies, "cookie value failed to URL-unescape")
				return "", errors.Format("please retry in a few seconds", errors.ErrSpartanTokenGrabFailure)
			}

			break
		}
	}

	if tokenValue == "" {
		dumpSpartanTokenDebug(location, resp, cookies, fmt.Sprintf("cookie %q not present in response", tokenName))
		return "", errors.Format("please retry in a few seconds", errors.ErrSpartanTokenGrabFailure)
	}

	return tokenValue, nil
}

// dumpSpartanTokenDebug writes diagnostic information to stderr and persists
// the callback response body so we can inspect what halowaypoint.com actually
// returned when the 343-spartan-token cookie isn't found. Called only on
// failure.
func dumpSpartanTokenDebug(requestedURL string, resp *http.Response, cookies []*http.Cookie, reason string) {
	if !debug.Enabled() {
		return
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "── spartan-token debug ────────────────────────")
	fmt.Fprintf(os.Stderr, "Reason       : %s\n", reason)
	fmt.Fprintf(os.Stderr, "Requested URL: %s\n", requestedURL)
	fmt.Fprintf(os.Stderr, "Status       : %d %s\n", resp.StatusCode, resp.Status)
	if loc := resp.Header.Get("Location"); loc != "" {
		fmt.Fprintf(os.Stderr, "Location     : %s\n", loc)
	}
	if resp.Request != nil && resp.Request.URL != nil {
		fmt.Fprintf(os.Stderr, "Final URL    : %s\n", resp.Request.URL.String())
	}

	if len(cookies) == 0 {
		fmt.Fprintln(os.Stderr, "Cookies      : (none)")
	} else {
		names := make([]string, 0, len(cookies))
		for _, c := range cookies {
			names = append(names, c.Name)
		}
		fmt.Fprintf(os.Stderr, "Cookies      : %s\n", strings.Join(names, ", "))
	}

	// Echo selected response headers that often help diagnose redirects /
	// auth failures (CSP, WWW-Authenticate, X-Frame-Options, etc.).
	for _, h := range []string{"Content-Type", "WWW-Authenticate", "X-Halo-Error", "Cf-Mitigated"} {
		if v := resp.Header.Get(h); v != "" {
			fmt.Fprintf(os.Stderr, "%-13s: %s\n", h, v)
		}
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Fprintf(os.Stderr, "Body length  : %d bytes\n", len(body))

	if home, err := os.UserHomeDir(); err == nil {
		dir := filepath.Join(home, strings.ReplaceAll(configs.GetConfig().Name, " ", "-"))
		path := filepath.Join(dir, "spartan-token-debug.html")
		if mkErr := os.MkdirAll(dir, 0755); mkErr == nil {
			if writeErr := os.WriteFile(path, body, 0644); writeErr == nil {
				fmt.Fprintf(os.Stderr, "Body saved   : %s\n", path)
			}
		}
	}

	fmt.Fprintln(os.Stderr, "───────────────────────────────────────────────")
	fmt.Fprintln(os.Stderr, "")
}
