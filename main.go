package main

import (
	"fmt"
	"infinite-bookmarker/internal"
	promptService "infinite-bookmarker/internal/services/prompt"
	"infinite-bookmarker/internal/shared/errors"
	"os"
)

func main() {
	var err error
	os.Stdout.WriteString(fmt.Sprintf("%s (%s)\n", internal.GetConfig().Title, internal.GetConfig().Version))

	err = promptService.StartAuthFlow()
	if err != nil {
		if errors.MayBe(err, errors.ErrAuthFailure) {
			promptService.DisplayAskOpenAuth()
			return
		}
		
		fmt.Println(err)
		return
	}

	err = promptService.DisplayBaseOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
}
