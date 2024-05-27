package prompt_svc

import (
	"fmt"
	identity_svc "infinite-ugc-tool/internal/application/services/auth/identity"
	halowaypoint_req "infinite-ugc-tool/pkg/libs/halowaypoint/modules/request"
	"infinite-ugc-tool/pkg/modules/errors"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func DisplayCloneOptions() error {
	var option string
	err := huh.NewSelect[string]().
		Title("🔄 What would like to clone?").
		Options(
			huh.NewOption(MAP, MAP),
			huh.NewOption(MODE, MODE),
			huh.NewOption(GO_BACK, GO_BACK),
		).Value(&option).Run()

	if err != nil || option == GO_BACK {
		return DisplayBaseOptions()
	}

	currentIdentity, err := identity_svc.GetActiveIdentity()
	if err != nil {
		return err
	}

	if option == MAP || option == MODE {
		var askForAssets bool
		err := huh.NewConfirm().
			Title("🔄 Would you like to clone the asset from an existing match?").
			Affirmative("No, I know what I'm doing.").
			Negative("Yes please!").
			Value(&askForAssets).
			Run()

		if err != nil {
			return DisplayCloneOptions()
		}

		if askForAssets {
			assetID, assetVersionID, err := displayVariantDetailsPrompt()
			if err != nil {
				if !errors.MayBe(err, errors.ErrPrompt) {
					os.Stdout.WriteString("❌ Invalid input...\n")
				}

				return DisplayCloneOptions()
			}

			if option == MAP {
				return cloneAsset(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					"maps",
					assetID,
					assetVersionID,
				)
			} else if option == MODE {
				return cloneAsset(
					currentIdentity.XboxNetwork.Xuid,
					currentIdentity.SpartanToken.Value,
					"ugcgamevariants",
					assetID,
					assetVersionID,
				)
			}

			return DisplayCloneOptions()
		}

		matchID, err := displayMatchGrabPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}

			return DisplayCloneOptions()
		}

		spinner.New().Title("Fetching...").Run()

		stats, err := halowaypoint_req.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		if err != nil {
			os.Stdout.WriteString("❌ Invalid match ID...\n")
			return DisplayCloneOptions()
		}

		if option == MAP {
			os.Stdout.WriteString(strings.Join([]string{
				fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
				"│ MapVariant",
				fmt.Sprintf("├── Asset ID: %s", stats.MatchInfo.MapVariant.AssetID),
				fmt.Sprintf("└── Version ID: %s", stats.MatchInfo.MapVariant.VersionID),
				"",
			}, "\n"))

			return cloneAsset(
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

			return cloneAsset(
				currentIdentity.XboxNetwork.Xuid,
				currentIdentity.SpartanToken.Value,
				"ugcgamevariants",
				stats.MatchInfo.UgcGameVariant.AssetID,
				stats.MatchInfo.UgcGameVariant.VersionID,
			)
		}
	}

	return DisplayCloneOptions()
}

func cloneAsset(xuid string, spartanToken string, category string, assetID string, assetVersionID string) error {
	spinner.New().Title("Cloning...").Run()
	err := halowaypoint_req.CloneAsset(xuid, spartanToken, category, assetID, assetVersionID)
	if err != nil {
		os.Stdout.WriteString("❌ Failed to clone the desired file...\n")
		return DisplayCloneOptions()
	}

	os.Stdout.WriteString("🎉 Cloned with success!\n")
	return DisplayCloneOptions()
}
