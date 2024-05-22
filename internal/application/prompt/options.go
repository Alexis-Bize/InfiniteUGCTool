package prompt

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func ProvideOptions() {
	prompt := promptui.Select{
		Label: "Options",
		Items: []string{
			"Bookmark a file",
			"Show credits",
			"Sign out",
		},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	fmt.Printf("You choose %q\n", result)
}
