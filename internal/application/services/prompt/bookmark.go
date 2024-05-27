package prompt_svc

import (
	"fmt"
	"os"
	"strings"

	identity_svc "infinite-ugc-tool/internal/application/services/auth/identity"
	halowaypoint_req "infinite-ugc-tool/pkg/libs/halowaypoint/modules/request"
	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func DisplayBookmarkOptions() error {
	var option string
	err := huh.NewSelect[string]().
		Title("🔖 What would like to bookmark?").
		Options(
			huh.NewOption(MAP, MAP),
			huh.NewOption(MODE, MODE),
			huh.NewOption(FILM, FILM),
			huh.NewOption(GO_BACK, GO_BACK),
		).Value(&option).Run()

	if err != nil || option == GO_BACK {
		return DisplayBaseOptions()
	}

	currentIdentity, err := identity_svc.GetActiveIdentity()
	if err != nil {
		return err
	}

	if option == FILM {
		matchID, err := displayMatchGrabPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}

			return DisplayBookmarkOptions()
		}

		spinner.New().Title("Fetching...").Run()

		stats, err := halowaypoint_req.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Invalid match ID...\n")
			return DisplayBookmarkOptions()
		}

		film, err := halowaypoint_req.GetMatchFilm(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Film not available...\n")
			return DisplayBookmarkOptions()
		}

		os.Stdout.WriteString(strings.Join([]string{
			fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
			"│ Film",
			fmt.Sprintf("└── Asset ID: %s", film.AssetID),
			"",
		}, "\n"))

		return bookmarkAsset(
			currentIdentity.XboxNetwork.Xuid,
			currentIdentity.SpartanToken.Value,
			"films",
			film.AssetID,
			"",
		)
	}

	if option == MAP || option == MODE {
		var askForAssets bool
		err := huh.NewConfirm().
			Title("🔖 Would you like to bookmark the asset from an existing match?").
			Affirmative("No, I know what I'm doing.").
			Negative("Yes please!").
			Value(&askForAssets).
			Run()

		if err != nil {
			return DisplayBookmarkOptions()
		}

		if askForAssets {
			assetID, assetVersionID, err := displayVariantDetailsPrompt()
			if err != nil {
				if !errors.MayBe(err, errors.ErrPrompt) {
					os.Stdout.WriteString("❌ Invalid input...\n")
				}

				return DisplayBookmarkOptions()
			}

			if option == MAP {
				return bookmarkAsset(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					"maps",
					assetID,
					assetVersionID,
				)
			} else if option == MODE {
				return bookmarkAsset(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					"ugcgamevariants",
					assetID,
					assetVersionID,
				)
			}

			return nil
		}

		matchID, err := displayMatchGrabPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}

			return DisplayBookmarkOptions()
		}

		spinner.New().Title("Fetching...").Run()

		stats, err := halowaypoint_req.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Invalid match ID...\n")
			return DisplayBookmarkOptions()
		}

		if option == MAP {
			os.Stdout.WriteString(strings.Join([]string{
				fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
				"│ MapVariant",
				fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.MapVariant.AssetID),
				fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.MapVariant.VersionID),
				"",
			}, "\n"))

			return bookmarkAsset(
				currentIdentity.XboxNetwork.Xuid,
				currentIdentity.SpartanToken.Value,
				"maps",
				stats.MatchInfo.MapVariant.AssetID,
				stats.MatchInfo.MapVariant.VersionID,
			)
		} else if option == MODE {
			os.Stdout.WriteString(strings.Join([]string{
				fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
				"│ UgcGameVariant",
				fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.UgcGameVariant.AssetID),
				fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.UgcGameVariant.VersionID),
				"",
			}, "\n"))

			return bookmarkAsset(
				currentIdentity.XboxNetwork.Xuid,
				currentIdentity.SpartanToken.Value,
				"ugcgamevariants",
				stats.MatchInfo.UgcGameVariant.AssetID,
				stats.MatchInfo.UgcGameVariant.VersionID,
			)
		}
	}

	return DisplayBookmarkOptions()
}

func displayAssetCloneFallbackOptions(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	var ignoreCloning bool
	err := huh.NewConfirm().
		Title("The desired asset is not published; would you like to try cloning it in your files instead?").
		Affirmative("No, that's ok.").
		Negative("Yes please!").
		Value(&ignoreCloning).
		Run()

	if err != nil || ignoreCloning {
		return DisplayBookmarkOptions()
	}

	return cloneAsset(xuid, spartanToken, category, assetID, assetVersionID)
}

func bookmarkAsset(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	var err error
	spinner.New().Title("Bookmarking...").Run()

	if category != "films" {
		err = halowaypoint_req.PingPublishedAsset(spartanToken, category, assetID)
		if err != nil {
			if errors.MayBe(err, errors.ErrNotFound) {
				if assetVersionID != "" {
					return displayAssetCloneFallbackOptions(xuid, spartanToken, category, assetID, assetVersionID)
				}

				os.Stdout.WriteString("❌ Failed to bookmark the desired file...\n")
				return DisplayBookmarkOptions()
			}

			os.Stdout.WriteString("❌ Something went wrong...\n")
			return DisplayBookmarkOptions()
		}
	}

	err = halowaypoint_req.Bookmark(xuid, spartanToken, category, assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Something went wrong...\n")
		return DisplayBookmarkOptions()
	}

	os.Stdout.WriteString("🎉 Bookmarked with success!\n")
	return DisplayBookmarkOptions()
}
