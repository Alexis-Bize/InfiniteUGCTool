package main

import (
	"embed"
	"fmt"
	"infinite-ugc-tool/internal"
	prompt_svc "infinite-ugc-tool/internal/services/prompt"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/helpers/release"
	"infinite-ugc-tool/internal/shared/modules/utilities"
	"os"

	"github.com/joho/godotenv"
)

//go:embed config.txt
var f embed.FS

func main() {
	godotenv.Load()
	internal.LoadConfig(f)

	os.Stdout.WriteString(fmt.Sprintf("# %s (%s)\n", internal.GetConfig().Name, internal.GetConfig().Version))

	err := exec(false)
	if err != nil {
		if !errors.MayBe(err, errors.ErrPrompt) {
			fmt.Println(err)
		}
	}
}

func exec(isRetry bool) error {
	var err error

	latestVersion, _ := release.CheckForUpdates()
	if latestVersion != "" {
		downloadLatestRelease, _ := prompt_svc.DisplayAskToUpdate(latestVersion)
		if downloadLatestRelease {
			return utilities.OpenBrowser(internal.GetConfig().GitHub + "/releases/latest")
		}
	}

	err = prompt_svc.StartAuthFlow(isRetry)
	if err != nil {
		if errors.MayBe(err, errors.ErrAuthFailure) {
			if prompt_svc.DisplayRetryAuth() {
				return exec(true)
			}
		}

		return err
	}

	err = prompt_svc.DisplayBaseOptions()
	return err
}
