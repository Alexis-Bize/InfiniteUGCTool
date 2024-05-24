package promptService

import (
	"fmt"
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
		Title("What would like to bookmark?").
		Options(
			huh.NewOption(BOOKMARK_MAP, BOOKMARK_MAP),
			huh.NewOption(BOOKMARK_MODE, BOOKMARK_MODE),
			huh.NewOption(BOOKMARK_FILM, BOOKMARK_FILM),
			huh.NewOption(GO_BACK, GO_BACK),
		).Value(&option).Run()

	if err != nil {
		return errors.Format(err.Error(), errors.ErrPrompt)
	} else if option == GO_BACK {
		return DisplayBaseOptions()
	}

	currentIdentity, err := identityService.GetActiveIdentity()
	if err != nil {
		return err
	}

	if option == BOOKMARK_FILM {
		matchID, err := DisplayMatchGrabPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}

			return DisplayBookmarkOptions()
		}

		spinner.New().Title("Fetching...").Run()

		stats, err := halowaypointRequest.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Invalid match ID...\n")
			return DisplayBookmarkOptions()
		}

		film, err := halowaypointRequest.GetMatchFilm(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Film not available...\n")
			return DisplayBookmarkOptions()
		}

		os.Stdout.WriteString(strings.Join([]string{
			fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
			"│ MapVariant",
			fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.MapVariant.AssetID),
			fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.MapVariant.VersionID),
			"│ UgcGameVariant",
			fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.UgcGameVariant.AssetID),
			fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.UgcGameVariant.VersionID),
			"│ Film",
			fmt.Sprintf("└── Asset ID: %s", film.AssetID),
			"",
		}, "\n"))

		return bookmarkFilm(
			currentIdentity.XboxNetwork.Xuid,
			currentIdentity.SpartanToken.Value,
			film.AssetID,
			"",
		)
	}

	if option == BOOKMARK_MAP || option == BOOKMARK_MODE {
		var askForAssets bool
		huh.NewConfirm().
			Title("Would you like to bookmark the asset from an existing match?").
			Affirmative("No, I know what I'm doing.").
			Negative("Yes please!").
			Value(&askForAssets).
			Run()

		if askForAssets {
			assetID, assetVersionID, err := DisplayBookmarkVariantPrompt()
			if err != nil {
				if !errors.MayBe(err, errors.ErrPrompt) {
					os.Stdout.WriteString("❌ Invalid input...\n")
				}

				return DisplayBookmarkOptions()
			}

			if option == BOOKMARK_MAP {
				return bookmarkMapVariant(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					assetID,
					assetVersionID,
				)
			} else if option == BOOKMARK_MODE {
				return bookmarkUgcGameVariant(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					assetID,
					assetVersionID,
				)
			}

			return nil
		}

		matchID, err := DisplayMatchGrabPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}

			return DisplayBookmarkOptions()
		}

		spinner.New().Title("Fetching...").Run()

		stats, err := halowaypointRequest.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Invalid match ID...\n")
			return DisplayBookmarkOptions()
		}

		if option == BOOKMARK_MAP {
			os.Stdout.WriteString(strings.Join([]string{
				fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
				"│ MapVariant",
				fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.MapVariant.AssetID),
				fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.MapVariant.VersionID),
				"",
			}, "\n"))

			return bookmarkMapVariant(
				currentIdentity.XboxNetwork.Xuid,
				currentIdentity.SpartanToken.Value,
				stats.MatchInfo.MapVariant.AssetID,
				stats.MatchInfo.MapVariant.VersionID,
			)
		} else if option == BOOKMARK_MODE {
			os.Stdout.WriteString(strings.Join([]string{
				fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
				"│ UgcGameVariant",
				fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.UgcGameVariant.AssetID),
				fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.UgcGameVariant.VersionID),
				"",
			}, "\n"))

			return bookmarkUgcGameVariant(
				currentIdentity.XboxNetwork.Xuid,
				currentIdentity.SpartanToken.Value,
				stats.MatchInfo.UgcGameVariant.AssetID,
				stats.MatchInfo.UgcGameVariant.VersionID,
			)
		}
	}

	return nil
}

func DisplayMatchGrabPrompt() (string, error) {
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

func DisplayBookmarkVariantPrompt() (string, string, error) {
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

	return assetID, assetVersionID, nil
}

func displayAssetCloneFallbackOptions(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	var ignoreCloning bool
	huh.NewConfirm().
		Title("The desired asset is not published; would you like to try cloning it in your files instead?").
		Affirmative("No, that's ok.").
		Negative("Yes please!").
		Value(&ignoreCloning).
		Run()

	if ignoreCloning {
		return DisplayBookmarkOptions()
	}

	err := halowaypointRequest.CloneAsset(xuid, spartanToken,category, assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Failed to clone the desired content...\n")
		return DisplayBookmarkOptions()
	}

	os.Stdout.WriteString("🎉 Cloned with success!\n")
	return DisplayBaseOptions()
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

func bookmarkFilm(xuid string, spartanToken string, assetID string, assetVersionID string) error {
	spinner.New().Title("Bookmarking your film...").Run()

	err := halowaypointRequest.Bookmark(xuid, spartanToken, "films", assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Failed to bookmark the desired content...\n")
		return DisplayBookmarkOptions()
	}

	os.Stdout.WriteString("🎉 Bookmarked with success!\n")
	return DisplayBaseOptions()
}

func bookmarkMapVariant(xuid string, spartanToken string, assetID string, assetVersionID string) error {
	var err error
	spinner.New().Title("Bookmarking your map...").Run()

	err = halowaypointRequest.PingPublishedAsset(spartanToken, "maps", assetID)
	if err != nil {
		if errors.MayBe(err, errors.ErrNotFound) && assetVersionID != "" {
			return displayAssetCloneFallbackOptions(xuid, spartanToken, "maps", assetID, assetVersionID)
		}

		os.Stdout.WriteString("❌ Something went wrong...\n")
		return DisplayBookmarkOptions()
	}

	err = halowaypointRequest.Bookmark(xuid, spartanToken, "maps", assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Failed to bookmark the desired content...\n")
		return DisplayBookmarkOptions()
	}

	os.Stdout.WriteString("🎉 Bookmarked with success!\n")
	return DisplayBaseOptions()
}

func bookmarkUgcGameVariant(xuid string, spartanToken string, assetID string, assetVersionID string) error {
	var err error
	spinner.New().Title("Bookmarking your mode...").Run()

	err = halowaypointRequest.PingPublishedAsset(spartanToken, "ugcgamevariants", assetID)
	if err != nil {
		if errors.MayBe(err, errors.ErrNotFound) && assetVersionID != "" {
			return displayAssetCloneFallbackOptions(xuid, spartanToken, "maps", assetID, assetVersionID)
		}

		os.Stdout.WriteString("❌ Failed to bookmark the desired content...\n")
		return DisplayBookmarkOptions()
	}

	err = halowaypointRequest.Bookmark(xuid, spartanToken, "ugcgamevariants", assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Something went wrong...\n")
		return DisplayBookmarkOptions()
	}

	os.Stdout.WriteString("🎉 Bookmarked with success!\n")
	return DisplayBaseOptions()
}