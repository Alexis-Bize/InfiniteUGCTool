package prompt_svc

import (
	"infinite-ugc-tool/internal"
	"infinite-ugc-tool/internal/shared/modules/errors"
	"infinite-ugc-tool/internal/shared/modules/helpers/identity"
	"infinite-ugc-tool/internal/shared/modules/utilities"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func DisplayBaseOptions() error {
	var option string
	err := huh.NewSelect[string]().
		Title("What would like to do today?").
		Options(
			huh.NewOption(BOOKMARK_FILES, BOOKMARK_FILES),
			huh.NewOption(CLONE_FILES, CLONE_FILES),
			huh.NewOption(SHOW_CREDITS, SHOW_CREDITS),
			huh.NewOption(SIGN_OUT, SIGN_OUT),
			huh.NewOption(EXIT, EXIT),
		).Value(&option).Run()

	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if option == SHOW_CREDITS {
		return DisplayCredits()
	} else if option == BOOKMARK_FILES {
		return DisplayBookmarkOptions()
	} else if option == CLONE_FILES {
		return DisplayCloneOptions()
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

	if err != nil || option == GO_BACK {
		return DisplayBaseOptions()
	}

	switch option {
		case OPEN_X_1:
			utilities.OpenBrowser("https://x.com/zeny_ic")
		case OPEN_X_2:
			utilities.OpenBrowser("https://x.com/_okom")
		case OPEN_X_3:
			utilities.OpenBrowser("https://x.com/gruntdotapi")
		case OPEN_GITHUB:
			utilities.OpenBrowser(internal.GetConfig().GitHub)
	}

	return DisplayBaseOptions()
}

func DisplayRetryAuth() bool {
	var confirm bool
	huh.NewConfirm().
		Title("❌ The authentication has failed; would you like to retry?").
		Affirmative("Yes, let's go!").
		Negative("No thanks.").
		Value(&confirm).
		Run()

	return confirm
}
