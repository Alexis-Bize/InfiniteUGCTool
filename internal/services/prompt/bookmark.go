package promptService

import (
	identityService "infinite-bookmarker/internal/services/identity"
	"infinite-bookmarker/internal/shared/errors"
	halowaypointRequest "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	"os"
	"regexp"

	"github.com/manifoldco/promptui"
)

func DisplayBookmarkOptions() error {
	prompt := promptui.Select{
		Label: "Options",
		Items: []string{
			BOOKMARK_FILM,
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

	if result == BOOKMARK_FILM {
		matchID, err := DisplayBookmarkFilmPrompt()
		if err != nil {
			return err
		}

		currentIdentity, err := identityService.GetActiveIdentity()
		if err != nil {
			return err
		}

		err = halowaypointRequest.BookmarkFilmFromMatchID(currentIdentity.XboxNetwork.Xuid, currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			return err
		}

		os.Stdout.WriteString("🎉 Bookmarked with success!\n")
		return DisplayBookmarkOptions()
	}

	return nil
}

func DisplayBookmarkFilmPrompt() (string, error) {
	prompt := promptui.Prompt{
		Label: "Match ID or URL (e.g., d6f60558-a14f-40c3-9016-d9085f6ec152)",
		Validate: func(input string) error {
			_, err := extractMatchID(input)
			return err
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	matchID, err := extractMatchID(result)
	if err != nil {
		return "", err
	}

	return matchID, nil
}

func extractMatchID(value string) (string, error) {
	const pattern = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	re := regexp.MustCompile(pattern)
	match := re.FindString(value)

	if match != "" {
		return match, nil
	}

	return "", errors.Format("invalid match ID or URL", errors.ErrPrompt)
}
