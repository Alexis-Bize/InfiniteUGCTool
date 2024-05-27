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

package release

import (
	"encoding/json"
	"fmt"
	"infinite-ugc-tool/configs"
	"infinite-ugc-tool/pkg/modules/errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/huh/spinner"
)

type Release struct {
	TagName string `json:"tag_name"`
}

type Releases []Release

func CheckForUpdates() (string, error) {
	spinner.New().Title("Checking for updates...").Run()

	owner, repo := extractOwnerAndRepo()
	if owner == "" || repo == "" {
		return "", errors.Format("something went wrong", errors.ErrInternal)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", owner, repo)
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Format("something went wrong", errors.ErrInternal)
	}

	var releases Releases
	err = json.NewDecoder(resp.Body).Decode(&releases)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	}

	latestRelease := releases[0]
	newVersionAvailable, err := isNewVersionAvailable(configs.GetConfig().Version, latestRelease.TagName)
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrInternal)
	} else if !newVersionAvailable {
		return "", nil
	}

	return latestRelease.TagName, nil
}

func isNewVersionAvailable(currentVersion, latestVersion string) (bool, error) {
	v1, err := semver.NewVersion(currentVersion)
	if err != nil {
		return false, err
	}

	v2, err := semver.NewVersion(latestVersion)
	if err != nil {
		return false, err
	}

	return v2.Compare(v1) == 1, nil
}

func extractOwnerAndRepo() (string, string) {
	githubUrl := configs.GetConfig().Repository
	u, err := url.Parse(githubUrl)
	if err != nil {
		return "", ""
	}

	path := strings.TrimPrefix(u.Path, "/")
	components := strings.Split(path, "/")
	if len(components) != 2 {
		return "", ""
	}

	return components[0], components[1]
}
