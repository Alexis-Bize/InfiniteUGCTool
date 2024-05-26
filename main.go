package main

import (
	"fmt"
	"os"

	"infinite-ugc-haven/internal"
	prompt_svc "infinite-ugc-haven/internal/services/prompt"
	"infinite-ugc-haven/internal/shared/modules/errors"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
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
