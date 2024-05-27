package prompt_svc

import (
	"infinite-ugc-tool/pkg/modules/errors"
	"regexp"
	"strings"

	"github.com/charmbracelet/huh"
)

const (
	BOOKMARK_FILES = "🔖 Bookmark files"
	CLONE_FILES = "🔄 Clone files"
	SHOW_CREDITS = "🤝 Show credits"
	SIGN_OUT = "🚫 Sign out"
	EXIT = "👋 Exit"
	GO_BACK = "← Go back"
)

const (
	BUNDLE = "📦 Bundle (MapModePair)"
	FILM = "🎬 Match (Film)"
	MODE = "🎮 Mode (UgcGameVariant)"
	MAP = "🌎 Map (MapVariant)"
)

const (
	OPEN_X_1 = "Made by: Zeny IC"
	OPEN_X_2 = "Original idea: Okom"
	OPEN_X_3 = "Supporter: Grunt.API"
	OPEN_GITHUB = "Source code: GitHub"
)

func displayMatchGrabPrompt() (string, error) {
	var value string
	var err error

	err = huh.NewInput().
		Title("Please specify a match ID or a valid match URL").
		Description("Leafapp.co, SpartanRecord.com, HaloDataHive.com and such are supported").
		Value(&value).
		Run()

	if err != nil {
		return "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	matchID, err := extractUUID(value)
	if err != nil {
		return "", err
	}

	return matchID, nil
}

func displayVariantDetailsPrompt() (string, string, error) {
	var assetID string
	var assetVersionID string
	var err error

	err = huh.NewInput().
		Title("Please specify a \"AssetID\" (GUID)").
		Description("e.g., ae4daed6-251a-4c2f-bc6f-eb25eac1bfd").
		Value(&assetID).
		Validate(func (input string) error {
			_, err := extractUUID(input)
			if err != nil {
				return errors.New("invalid GUID")
			}

			return nil
		}).Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	err = huh.NewInput().
		Title("Please specify a \"AssetVariantID\" (GUID)").
		Description("This value is optional for published files").
		Value(&assetVersionID).
		Run()

	if err != nil {
		return "", "", errors.Format(err.Error(), errors.ErrPrompt)
	}

	assetID = strings.TrimSpace(assetID)
	assetVersionID = strings.TrimSpace(assetVersionID)

	if assetVersionID != "" {
		_, err := extractUUID(assetVersionID)
		if err != nil {
			return "", "", err
		}
	}

	return strings.TrimSpace(assetID), strings.TrimSpace(assetVersionID), nil
}

func extractUUID(value string) (string, error) {
	const pattern = `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`
	re := regexp.MustCompile(pattern)
	match := re.FindString(strings.TrimSpace(value))

	if match != "" {
		return match, nil
	}

	return "", errors.Format("invalid format", errors.ErrUUIDInvalid)
}
