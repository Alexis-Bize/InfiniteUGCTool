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

package internal

import (
	"fmt"
	"infinite-ugc-tool/configs"
	prompt_svc "infinite-ugc-tool/internal/application/services/prompt"
	"infinite-ugc-tool/internal/helpers/release"
	"infinite-ugc-tool/pkg/modules/errors"
	"infinite-ugc-tool/pkg/modules/utilities"
	"os"
)

func Exec() {
	os.Stdout.WriteString(fmt.Sprintf("# %s (%s)\n", configs.GetConfig().Name, configs.GetConfig().Version))

	err := run(false)
	if err != nil {
		if !errors.MayBe(err, errors.ErrPrompt) {
			fmt.Println(err)
		}
	}
}

func run(isRetry bool) error {
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
				return run(true)
			}
		}

		return err
	}

	err = prompt_svc.DisplayBaseOptions()
	return err
}
