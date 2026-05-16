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
//
// The Microsoft Live pre-auth and credential POST flow in this file —
// notably the PPFT / urlPost extraction regexes and the Set-Cookie
// reduction — is inspired by the XboxReplay/xboxlive-auth Node.js library:
// https://github.com/XboxReplay/xboxlive-auth

package msa_req

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"infinite-ugc-tool/configs"
	"infinite-ugc-tool/pkg/libs/msa"
	"infinite-ugc-tool/pkg/modules/debug"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
)

// Authenticate posts the user's credentials against the login form contained
// in the supplied authorize page. The page is fetched by the caller (via
// GetAuthorizePage) so that the host / client_id / cookies all originate
// from halowaypoint's /sign-in redirect rather than a hardcoded URL.
func Authenticate(credentials msa.LiveCredentials, page *AuthorizePageResponse) (*http.Response, error) {
	preAuthResponse, err := preAuthFromPage(page)
	if err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Add("login", credentials.Email)
	form.Add("loginfmt", credentials.Email)
	form.Add("passwd", credentials.Password)
	form.Add("PPFT", preAuthResponse.Matches.PPFT)

	req, err := http.NewRequest("POST", preAuthResponse.Matches.URLPost, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Cookie": preAuthResponse.Cookie,
		"Content-Type": "application/x-www-form-urlencoded",
	}) { req.Header.Set(k, v) }

	resp, err := request.NoRedirectClient.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, errors.Format(err.Error(), errors.ErrInternal)
	}

	cookie := extractCookieHeader(resp.Header["Set-Cookie"])

	switch {
	case strings.Contains(string(body), "PROOF.Type"):
		// Legacy 2FA / OTC challenge.
		serverData, err := extractJSObjectAndConvertToJSON(body, "ServerData")
		if err != nil {
			return nil, err
		}
		return requestOneTimeCode(credentials.Email, serverData, cookie)

	case isKMSIPage(body):
		// "Keep me signed in?" interstitial on the Fluent v2 UI.
		// MS already accepted the credentials (the __Host-MSAAUTH cookies
		// are set); we just need to dismiss the prompt with LoginOptions=3
		// to receive the final auth-code redirect.
		return submitKMSI(body, cookie)

	default:
		dumpAuthFailureDebug(resp, body)
		return resp, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}
}

// isKMSIPage recognises the v2 Fluent "Keep me signed in?" page. The
// surest tell is the page's JS bundle name (`kmsi-fluent_v2…`); we also
// check the hpgid for older variants of the same prompt.
func isKMSIPage(body []byte) bool {
	s := string(body)
	return strings.Contains(s, "kmsi-fluent") || strings.Contains(s, `"hpgid":93`)
}

// submitKMSI POSTs the "no, don't stay signed in" answer (LoginOptions=3)
// to the KMSI form's urlPost so MS issues the final auth-code redirect.
func submitKMSI(body []byte, cookie string) (*http.Response, error) {
	serverData, err := extractJSObjectAndConvertToJSON(body, "ServerData")
	if err != nil {
		return nil, err
	}

	sFT := serverData.Get("sFT").Str
	urlPost := serverData.Get("urlPost").Str

	if sFT == "" || urlPost == "" {
		return nil, errors.Format("KMSI form parameters missing", errors.ErrAuthFailure)
	}

	form := url.Values{}
	form.Add("LoginOptions", "3")
	form.Add("type", "28")
	form.Add("PPFT", sFT)

	req, err := http.NewRequest("POST", urlPost, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Cookie":       cookie,
		"Content-Type": "application/x-www-form-urlencoded",
	}) { req.Header.Set(k, v) }

	resp, err := request.NoRedirectClient.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := io.ReadAll(resp.Body)
		dumpAuthFailureDebug(resp, body)
		return resp, errors.Format("KMSI submission did not redirect", errors.ErrAuthFailure)
	}

	return resp, nil
}

// dumpAuthFailureDebug writes diagnostic information to stderr and persists
// the response body when the credentials POST returns 200 but the page is
// neither the legacy 2FA challenge (PROOF.Type) nor a redirect to the
// callback. Common causes: a consent / "stay signed in" page, a passkey
// prompt, an "incorrect password" form re-render, or an unexpected MS UI.
func dumpAuthFailureDebug(resp *http.Response, body []byte) {
	if !debug.Enabled() {
		return
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "── credentials POST debug ─────────────────────")
	fmt.Fprintf(os.Stderr, "Status         : %d %s\n", resp.StatusCode, resp.Status)
	if resp.Request != nil && resp.Request.URL != nil {
		fmt.Fprintf(os.Stderr, "Final URL      : %s\n", resp.Request.URL.String())
	}
	if loc := resp.Header.Get("Location"); loc != "" {
		fmt.Fprintf(os.Stderr, "Location       : %s\n", loc)
	}
	if ct := resp.Header.Get("Content-Type"); ct != "" {
		fmt.Fprintf(os.Stderr, "Content-Type   : %s\n", ct)
	}

	cookies := resp.Cookies()
	if len(cookies) > 0 {
		names := make([]string, 0, len(cookies))
		for _, c := range cookies {
			names = append(names, c.Name)
		}
		fmt.Fprintf(os.Stderr, "Set-Cookie     : %s\n", strings.Join(names, ", "))
	}

	fmt.Fprintf(os.Stderr, "Body length    : %d bytes\n", len(body))

	// Surface telltale strings that hint at WHICH non-2FA page we landed on.
	hints := []string{
		"PROOF.Type",       // legacy 2FA
		"passkey",          // FIDO / passkey prompt
		"KMSI",             // "Keep me signed in"
		"AppConsent",       // OAuth consent
		"consent",          // generic consent
		"incorrect",        // bad password
		"locked",           // account locked
		"sErrTxt",          // legacy error text container
		"sErrorCode",       // legacy error code
		"$Config",          // v2.0 MSAL config blob
		"hasError",         // v2.0 error flag
		"urlEndAuth",       // v2.0 OAuth end
		"sFTName",          // legacy form token
	}
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, "Has %-13s: %v\n", "'"+h+"'", strings.Contains(string(body), h))
	}

	// Context around any error-ish strings we found.
	for _, marker := range []string{"sErrTxt", "sErrorCode", "hasError", "ErrorTitle"} {
		if idx := strings.Index(string(body), marker); idx >= 0 {
			start := max(0, idx-20)
			end := min(len(body), idx+240)
			fmt.Fprintf(os.Stderr, "%-13s ctx: ...%s...\n", marker, string(body[start:end]))
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		dir := filepath.Join(home, strings.ReplaceAll(configs.GetConfig().Name, " ", "-"))
		path := filepath.Join(dir, "auth-failure.html")
		if mkErr := os.MkdirAll(dir, 0755); mkErr == nil {
			if writeErr := os.WriteFile(path, body, 0644); writeErr == nil {
				fmt.Fprintf(os.Stderr, "Body saved     : %s\n", path)
			}
		}
	}

	fmt.Fprintln(os.Stderr, "───────────────────────────────────────────────")
	fmt.Fprintln(os.Stderr, "")
}

func requestOneTimeCode(email string, serverData gjson.Result, cookie string) (*http.Response, error) {
	sFT := serverData.Get("sFT").Str
	urlPost := serverData.Get("urlPost").Str
	proof := serverData.Get("p|@reverse|0|data").Str

	if sFT == "" || urlPost == "" || proof == "" {
		return nil, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}

	var otc string
	err := huh.NewInput().
		Title("📱 Two Factor Authentication: Please enter the code displayed in your authenticator application to continue").
		Value(&otc).
		Validate(func (input string) error {
			if len(input) == 0 {
				return errors.New("code can not be empty")
			}

			return nil
		}).Run()

	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrPrompt)
	}

	form := url.Values{}
	form.Add("type", "19")
	form.Add("SentProofIDE", proof)
	form.Add("otc", strings.TrimSpace(otc))
	form.Add("login",email)
	form.Add("PPFT", sFT)

	req, err := http.NewRequest("POST", urlPost, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Cookie": cookie,
		"Content-Type": "application/x-www-form-urlencoded",
	}) { req.Header.Set(k, v) }

	var resp *http.Response
	var doErr error
	spinner.New().Title("Validating...").Action(func() {
		resp, doErr = request.NoRedirectClient.Do(req)
	}).Run()

	if doErr != nil {
		return nil, errors.Format(doErr.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return resp, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}

	return resp, nil
}

// preAuthFromPage extracts the PPFT, urlPost, and login-page cookies from
// an already-fetched authorize page (see GetAuthorizePage). The legacy
// `preAuth` used to fetch the page itself against a hardcoded login.live.com
// URL — that path stopped working when MS retired the public-client shape
// for that endpoint, so the fetch is now driven by halowaypoint's /sign-in
// redirect and the cookies/body are passed in.
func preAuthFromPage(page *AuthorizePageResponse) (*msa.LivePreAuthResponse, error) {
	matches := msa.LivePreAuthMatchedParameters{
		URLPost: matchForIndex(string(page.Body), urlPostPattern, 1),
		PPFT:    matchForIndex(string(page.Body), ppftPattern, 1),
	}

	if matches.PPFT == "" || matches.URLPost == "" {
		dumpPreAuthDebug(page.Response, page.Body, matches)
		return nil, errors.Format("please retry in a few seconds", errors.ErrPreAuthFailure)
	}

	return &msa.LivePreAuthResponse{
		Cookie:  extractCookieHeader(page.Response.Header["Set-Cookie"]),
		Matches: matches,
	}, nil
}

// dumpPreAuthDebug writes diagnostic information to stderr and persists the
// full response body to disk so we can inspect what Microsoft actually
// served when the PPFT / urlPost regexes fail to match. Called only on
// failure — normal users never see this.
func dumpPreAuthDebug(resp *http.Response, body []byte, matches msa.LivePreAuthMatchedParameters) {
	if !debug.Enabled() {
		return
	}

	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "── preAuth debug ──────────────────────────────")
	fmt.Fprintf(os.Stderr, "Status      : %d %s\n", resp.StatusCode, resp.Status)
	if resp.Request != nil && resp.Request.URL != nil {
		fmt.Fprintf(os.Stderr, "Final URL   : %s\n", resp.Request.URL.String())
	}
	fmt.Fprintf(os.Stderr, "Body length : %d bytes\n", len(body))
	fmt.Fprintf(os.Stderr, "Has 'PPFT'  : %v\n", bytes.Contains(body, []byte("PPFT")))
	fmt.Fprintf(os.Stderr, "Has urlPost : %v\n", bytes.Contains(body, []byte("urlPost")))
	fmt.Fprintf(os.Stderr, "PPFT match  : %q\n", matches.PPFT)
	fmt.Fprintf(os.Stderr, "URLPost mat.: %q\n", matches.URLPost)

	if idx := bytes.Index(body, []byte("PPFT")); idx >= 0 {
		start := max(0, idx-40)
		end := min(len(body), idx+160)
		fmt.Fprintf(os.Stderr, "PPFT ctx    : ...%s...\n", string(body[start:end]))
	}

	if idx := bytes.Index(body, []byte("urlPost")); idx >= 0 {
		start := max(0, idx-20)
		end := min(len(body), idx+200)
		fmt.Fprintf(os.Stderr, "urlPost ctx : ...%s...\n", string(body[start:end]))
	}

	if home, err := os.UserHomeDir(); err == nil {
		dir := filepath.Join(home, strings.ReplaceAll(configs.GetConfig().Name, " ", "-"))
		path := filepath.Join(dir, "preauth-debug.html")
		if mkErr := os.MkdirAll(dir, 0755); mkErr == nil {
			if writeErr := os.WriteFile(path, body, 0644); writeErr == nil {
				fmt.Fprintf(os.Stderr, "Body saved  : %s\n", path)
			}
		}
	}

	fmt.Fprintln(os.Stderr, "───────────────────────────────────────────────")
	fmt.Fprintln(os.Stderr, "")
}

// ppftPattern matches the PPFT input on the login page, tolerating both
// plain HTML (`name="PPFT" ... value="..."`) and the JS-escaped variant
// (`name=\"PPFT\" ... value=\"...\"`) that Microsoft serves in some
// contexts. Ported from XboxReplay/xboxlive-auth.
var ppftPattern = regexp.MustCompile(`(?i)name=\\?"PPFT\\?"[^>]*value=\\?"([^"\\]+)\\?"`)

// urlPostPattern matches the `urlPost: "..."` (or single-quoted, or
// JS-escaped) field in the ServerData blob. Ported from
// XboxReplay/xboxlive-auth — the previous Go regex only matched the
// single-quoted variant, which Microsoft no longer always serves.
var urlPostPattern = regexp.MustCompile(`(?i)\\?['"]?urlPost\\?['"]?:\s*\\?['"]([^'"\\]+)\\?['"]`)

// extractCookieHeader collapses Set-Cookie response headers into a single
// Cookie request header value, keeping only the `name=value` portion of
// each cookie and dropping attributes (Path, Domain, Expires, ...). Mirrors
// the behaviour of the reference TS implementation.
func extractCookieHeader(setCookies []string) string {
	parts := make([]string, 0, len(setCookies))
	for _, c := range setCookies {
		nameValue := strings.TrimSpace(strings.SplitN(c, ";", 2)[0])
		if nameValue != "" {
			parts = append(parts, nameValue)
		}
	}
	return strings.Join(parts, "; ")
}

func matchForIndex(body string, pattern *regexp.Regexp, index int) string {
	matches := pattern.FindStringSubmatch(body)
	if len(matches) > index {
		return matches[index]
	}

	return ""
}

func extractJSObjectAndConvertToJSON(body []byte, objName string) (gjson.Result, error) {
	bodyStr := string(body)
	serverDataRe := fmt.Sprintf(`%s(?:.*?)=(?:.*?)({(?:.*?)});`, regexp.QuoteMeta(objName))

	re, err := regexp.Compile(serverDataRe)
	if err != nil {
		return gjson.Result{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	matches := re.FindStringSubmatch(bodyStr)
	if len(matches) == 0 {
		return gjson.Result{}, errors.Format("no matches found", errors.ErrInternal)
	}

	JSObject := matches[1]
	JSCode := "JSON.stringify(" + JSObject + ");"

	vm := otto.New()
	value, err := vm.Run(JSCode)
	if err != nil {
		return gjson.Result{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	toString, err := value.ToString()
	if err != nil {
		return gjson.Result{}, errors.Format(err.Error(), errors.ErrInternal)
	}

	return gjson.Parse(toString), nil
}
