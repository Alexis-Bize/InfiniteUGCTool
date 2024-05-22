package promptService

import (
	"fmt"
	"infinite-bookmarker/internal"
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func DisplayBaseOptions() error {
	prompt := promptui.Select{
		Label: "Options",
		Items: []string{
			BOOKMARK,
			SHOW_CREDITS,
			SIGN_OUT,
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if result == SHOW_CREDITS {
		return DisplayCredits()
	} else if result == BOOKMARK {
		return DisplayBookmarkOptions()
	} else if result == SIGN_OUT {
		os.Stdout.WriteString("👋 Good bye!")
		identity.SaveIdentity(identity.Identity{})
	}

	return nil
}

func DisplayCredits() error {
	credits := []string {
		fmt.Sprintf("%s (%s)", internal.GetConfig().Title, internal.GetConfig().Version),
		"+ Made by: Zeny IC (x.com/@Zeny_IC)",
		"+ Original idea: Okom (x.com/@_Okom)",
		"+ Supporter: Grunt.API (x.com/@GruntDotAPI)",
		"+ GitHub: github.com/Alexis-Bize/Infinite-Bookmarker",
	}

	os.Stdout.WriteString(strings.Join(credits, "\n"))

	prompt := promptui.Select{
		Label: "Options",
		Items: []string{
			GO_BACK,
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if result == GO_BACK {
		return DisplayBaseOptions()
	}

	return nil
}
