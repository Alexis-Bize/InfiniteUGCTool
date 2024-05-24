package promptService

import (
	"infinite-bookmarker/internal/shared/errors"
	"infinite-bookmarker/internal/shared/modules/helpers/identity"
	"infinite-bookmarker/internal/shared/modules/utilities"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func DisplayBaseOptions() error {
	var option string
	err := huh.NewSelect[string]().
		Title("What would like to do today?").
		Options(
			huh.NewOption(BOOKMARK, BOOKMARK),
			huh.NewOption(SHOW_CREDITS, SHOW_CREDITS),
			huh.NewOption(SIGN_OUT, SIGN_OUT),
			huh.NewOption(EXIT, EXIT),
		).Value(&option).Run()

	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if option == SHOW_CREDITS {
		return DisplayCredits()
	} else if option == BOOKMARK {
		return DisplayBookmarkOptions()
	} else if option == SIGN_OUT {
		var confirm bool
		huh.NewConfirm().
			Title("Are you sure?").
			Affirmative("Yes!").
			Negative("Oops, nevermind.").
			Value(&confirm).
			Run()

		if confirm {
			spinner.New().Title("Signing out...").Run()
			identity.SaveIdentity(identity.Identity{})
			return nil
		}

		return DisplayBaseOptions()
	}

	return nil
}

func DisplayCredits() error {
	var option string
	err := huh.NewSelect[string]().Title("Credits:").Options(
		huh.NewOption(OPEN_X_1, OPEN_X_1),
		huh.NewOption(OPEN_X_2, OPEN_X_2),
		huh.NewOption(OPEN_X_3, OPEN_X_3),
		huh.NewOption(OPEN_GITHUB, OPEN_GITHUB),
		huh.NewOption(GO_BACK, GO_BACK),
	).Value(&option).Run()

	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if option == GO_BACK {
		return DisplayBaseOptions()
	}

	spinner.New().Title("Attempting to open your browser...").Run()

	switch option {
		case OPEN_X_1:
			utilities.OpenBrowser("https://x.com/zeny_ic")
		case OPEN_X_2:
			utilities.OpenBrowser("https://x.com/_okom")
		case OPEN_X_3:
			utilities.OpenBrowser("https://x.com/gruntdotapi")
		case OPEN_GITHUB:
			utilities.OpenBrowser("https://github.com/Alexis-Bize/Infinite-Bookmarker")
	}

	return DisplayBaseOptions()
}

func DisplayAskOpenAuth() bool {
	var confirm bool
	huh.NewConfirm().
		Title("❌ The authentication has failed; would you like to retry?").
		Affirmative("Yes, let's go!").
		Negative("No thanks.").
		Value(&confirm).
		Run()

	return confirm
}
