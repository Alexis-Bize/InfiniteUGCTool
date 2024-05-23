package main

import (
	"fmt"
	"infinite-bookmarker/internal"
	promptService "infinite-bookmarker/internal/services/prompt"
	"infinite-bookmarker/internal/shared/errors"
	"os"
)

func main() {
	os.Stdout.WriteString(fmt.Sprintf("# %s (%s)\n", internal.GetConfig().Title, internal.GetConfig().Version))

	err := exec(false)
	if err != nil {
		if !errors.MayBe(err, errors.ErrPrompt) {
			fmt.Println(err)
		}
	}
}

func exec(isRetry bool) error {
	var err error

	err = promptService.StartAuthFlow(isRetry)
	if err != nil {
		if errors.MayBe(err, errors.ErrAuthFailure) {
			if promptService.DisplayAskOpenAuth() {
				return exec(true)
			}
		}

		return err
	}

	err = promptService.DisplayBaseOptions()
	return err
}
