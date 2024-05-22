package main

import (
	"fmt"
	"infinite-bookmarker/internal"
	"infinite-bookmarker/internal/application/prompt"
	"os"
)

func main() {
	title := fmt.Sprintf("%s (%s)\n", internal.GetConfig().Title, internal.GetConfig().Version)
	os.Stdout.WriteString(title)

	err := prompt.StartAuthFlow()
	if err != nil {
		fmt.Println(err)
		return
	}

	prompt.ProvideOptions()
}
