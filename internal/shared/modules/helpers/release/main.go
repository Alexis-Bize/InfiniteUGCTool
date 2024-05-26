package release

import (
	"encoding/json"
	"fmt"
	"infinite-ugc-tool/internal"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/huh/spinner"
)

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
	newVersionAvailable, err := isNewVersionAvailable(internal.GetConfig().Version, latestRelease.TagName)
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
	githubUrl := internal.GetConfig().GitHub
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