package main

import (
	"fmt"
	"infinite-bookmarker/internal"
	promptService "infinite-bookmarker/internal/services/prompt"
	"os"
)

func main() {
	title := fmt.Sprintf("%s (%s)\n", internal.GetConfig().Title, internal.GetConfig().Version)
	os.Stdout.WriteString(title)

	err := promptService.StartAuthFlow()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = promptService.DisplayBaseOptions()
	if err != nil {
		fmt.Println(err)
		return
	}
}
