// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompt_svc

import (
	"fmt"
	"os"
	"strings"

	identity_svc "infinite-ugc-tool/internal/application/services/auth/identity"
	"infinite-ugc-tool/internal/helpers/identity"
	"infinite-ugc-tool/pkg/libs/halowaypoint"
	halowaypoint_req "infinite-ugc-tool/pkg/libs/halowaypoint/modules/request"
	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh"
)

func DisplayCloneOptions() error {
	for {
		var option string
		err := huh.NewSelect[string]().
			Title("🔄 What would like to clone?").
			Options(
				huh.NewOption(MAP, MAP),
				huh.NewOption(MODE, MODE),
				huh.NewOption(GO_BACK, GO_BACK),
			).Value(&option).Run()

		if err != nil || option == GO_BACK {
			return nil
		}

		currentIdentity, err := identity_svc.GetActiveIdentity()
		if err != nil {
			return err
		}

		switch option {
		case MAP:
			handleAssetClone(&currentIdentity, MAP, "maps")
		case MODE:
			handleAssetClone(&currentIdentity, MODE, "ugcgamevariants")
		}
	}
}

func handleAssetClone(currentIdentity *identity.Identity, option string, category string) {
	var manualEntry bool
	err := huh.NewConfirm().
		Title("🔄 Would you like to clone the asset from an existing match?").
		Affirmative("No, I know what I'm doing.").
		Negative("Yes please!").
		Value(&manualEntry).
		Run()

	if err != nil {
		return
	}

	if manualEntry {
		assetID, assetVersionID, err := displayVariantDetailsPrompt()
		if err != nil {
			if !errors.MayBe(err, errors.ErrPrompt) {
				os.Stdout.WriteString("❌ Invalid input...\n")
			}
			return
		}
		cloneAsset(currentIdentity, category, assetID, assetVersionID)
		return
	}

	matchID, err := displayMatchGrabPrompt()
	if err != nil {
		if !errors.MayBe(err, errors.ErrPrompt) {
			os.Stdout.WriteString("❌ Invalid input...\n")
		}
		return
	}

	var stats halowaypoint.MatchStatsResponse
	statsErr := runWithSpinnerAndRefresh(currentIdentity, "Fetching...", func() error {
		var err error
		stats, err = halowaypoint_req.GetMatchStats(currentIdentity.SpartanToken.Value, matchID)
		return err
	})

	if statsErr != nil {
		os.Stdout.WriteString("❌ Invalid match ID...\n")
		return
	}

	var label, assetID, versionID string
	switch option {
	case MAP:
		label = "MapVariant"
		assetID = stats.MatchInfo.MapVariant.AssetID
		versionID = stats.MatchInfo.MapVariant.VersionID
	case MODE:
		label = "UgcGameVariant"
		assetID = stats.MatchInfo.UgcGameVariant.AssetID
		versionID = stats.MatchInfo.UgcGameVariant.VersionID
	}

	os.Stdout.WriteString(strings.Join([]string{
		fmt.Sprintf("Match Details (ID: %s)", stats.MatchID),
		"│ " + label,
		fmt.Sprintf("├── Asset ID: %s", assetID),
		fmt.Sprintf("└── Version ID: %s", versionID),
		"",
	}, "\n"))

	cloneAsset(currentIdentity, category, assetID, versionID)
}

func cloneAsset(currentIdentity *identity.Identity, category string, assetID string, assetVersionID string) {
	err := runWithSpinnerAndRefresh(currentIdentity, "Cloning...", func() error {
		return halowaypoint_req.CloneAsset(currentIdentity.XboxNetwork.Xuid, currentIdentity.SpartanToken.Value, category, assetID, assetVersionID)
	})

	if err != nil {
		os.Stdout.WriteString("❌ Failed to clone the desired file...\n")
		return
	}

	os.Stdout.WriteString("🎉 Cloned with success!\n")
}
