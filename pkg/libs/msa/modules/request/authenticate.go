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
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"infinite-ugc-tool/pkg/libs/msa"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities/request"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/robertkrimen/otto"
	"github.com/tidwall/gjson"
)

func Authenticate(credentials msa.LiveCredentials, options msa.LiveClientAuthOptions) (*http.Response, error) {
	preAuthResponse, err := preAuth(&options)
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

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return resp, nil
	}

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return resp, errors.Format(err.Error(), errors.ErrInternal)
	} else if !strings.Contains(string(body), "PROOF.Type") {
		return resp, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}

	serverData, err := extractJSObjectAndConvertToJSON(body, "ServerData")
	if err != nil {
		return nil, err
	}

	return requestOneTimeCode(credentials.Email, serverData, strings.Join(resp.Header["Set-Cookie"], "; "))
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

	spinner.New().Title("Validating...").Run()

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

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return resp, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}

	return resp, nil
}

func preAuth(options *msa.LiveClientAuthOptions) (*msa.LivePreAuthResponse, error) {
	url := BuildAuthorizeUrl(
		options.ClientID,
		options.Scope,
		options.ResponseType,
		options.RedirectURI,
		options.State,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{}) {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	cookie := strings.Join(resp.Header["Set-Cookie"], "; ")
	urlPostPattern := regexp.MustCompile(`urlPost:'([^']+)'`)
	ppftPattern := regexp.MustCompile(`sFTTag:'.*value="(.*)"\/>'`)

	matches := msa.LivePreAuthMatchedParameters{
		URLPost: matchForIndex(string(body), urlPostPattern, 1),
		PPFT: matchForIndex(string(body), ppftPattern, 1),
	}

	if matches.PPFT == "" || matches.URLPost == "" {
		return nil, errors.Format("please retry in a few seconds", errors.ErrPreAuthFailure)
	}

	return &msa.LivePreAuthResponse{
		Cookie:  cookie,
		Matches: matches,
	}, nil
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
