package main

import (
	"fmt"
	"infinite-ugc-tool/configs"
	prompt_svc "infinite-ugc-tool/internal/application/services/prompt"
	"infinite-ugc-tool/internal/helpers/release"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities"
	"os"

	"github.com/joho/godotenv"
)

//go:generate goversioninfo -icon=assets/resource/icon.ico

func main() {
	godotenv.Load()
	os.Stdout.WriteString(fmt.Sprintf("# %s (%s)\n", configs.GetConfig().Name, configs.GetConfig().Version))

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
			return utilities.OpenBrowser(configs.GetConfig().Repository + "/releases/latest")
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
