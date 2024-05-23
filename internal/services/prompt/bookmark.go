package promptService

import (
	identityService "infinite-bookmarker/internal/services/identity"
	"infinite-bookmarker/internal/shared/errors"
	halowaypointRequest "infinite-bookmarker/internal/shared/libs/halowaypoint/modules/request"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func DisplayBookmarkOptions() error {
	var option string
	err := huh.NewSelect[string]().
		Title("What would like to do today?").
		Options(
			huh.NewOption(BOOKMARK_FILM, BOOKMARK_FILM),
			huh.NewOption(GO_BACK, GO_BACK),
		).Value(&option).Run()

	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}
	
	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	}

	if option == GO_BACK {
		return DisplayBaseOptions()
	}

	if option == BOOKMARK_FILM {
		matchID, err := DisplayBookmarkFilmPrompt()
		if err != nil {
			if errors.MayBe(err, errors.ErrMatchIdInvalid) {
				os.Stdout.WriteString("❌ Seems like your match ID or URL is incorrect...\n")
				return DisplayBookmarkOptions()
			}

			return err
		}

		currentIdentity, err := identityService.GetActiveIdentity()
		if err != nil {
			return err
		}

		spinner.New().Title("Bookmarking...").Run()

		err = halowaypointRequest.BookmarkFilmFromMatchID(currentIdentity.XboxNetwork.Xuid, currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			return err
		}

		os.Stdout.WriteString("🎉 Bookmarked with success!\n")
		return DisplayBaseOptions()
	}

	return nil
}

func DisplayBookmarkFilmPrompt() (string, error) {
	var value string
	err := huh.NewInput().
		Title("Please specify a match ID or a valid URL (Leafapp.co, SpartanRecord.com, HaloDataHive.com, etc.)").
		Value(&value).
		Run()

	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	matchID, err := extractMatchID(strings.TrimSpace(value))
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

	return "", errors.Format("please retry", errors.ErrMatchIdInvalid)
}
