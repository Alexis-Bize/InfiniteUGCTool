package msaRequest

import (
	"fmt"
	"infinite-bookmarker/internal/shared/libs/msa"
	"infinite-bookmarker/internal/shared/modules/errors"
	"infinite-bookmarker/internal/shared/modules/utilities/request"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
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
	if err != nil {
		return resp, errors.Format(err.Error(), errors.ErrInternal)
	} else if !strings.Contains(string(body), "PROOF.Type") {
		return resp, errors.Format("the authentication has failed", errors.ErrAuthFailure)
	}

	serverData, err := extractJSObjectAndConvertToJSON(body, "ServerData")
	if err != nil {
		return nil, err
	}

	sFT := serverData.Get("sFT")
	urlPost := serverData.Get("urlPost")
	proof := serverData.Get("p|@reverse|0|data")
	
	var otc string
	err = huh.NewInput().
		Title("Enter the code displayed in the Microsoft app (such as Authenticator or Outlook) you use for approving sign-in requests").
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

	form = url.Values{}
	form.Add("type", "19")
	form.Add("SentProofIDE", proof.Str)
	form.Add("otc", otc)
	form.Add("AddTD", "true")
	form.Add("login", credentials.Email)
	form.Add("PPFT", sFT.Str)

	req, err = http.NewRequest("POST", urlPost.Str, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, errors.Format(err.Error(), errors.ErrInternal)
	}

	for k, v := range request.GetBaseHeaders(map[string]string{
		"Cookie": strings.Join(resp.Header["Set-Cookie"], "; "),
		"Content-Type": "application/x-www-form-urlencoded",
	}) { req.Header.Set(k, v) }

	client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err = client.Do(req)
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