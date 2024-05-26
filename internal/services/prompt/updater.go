package prompt_svc

import (
	"fmt"
	"infinite-ugc-tool/internal/shared/modules/errors"

	"github.com/charmbracelet/huh"
)

func DisplayAskToUpdate(version string) (bool, error) {
	var ignoreUpdate bool
	err := huh.NewConfirm().
		Title(fmt.Sprintf("🔥 A new version (%s) is available; would you like to download it?", version)).
		Affirmative("Later").
		Negative("Yes please!").
		Value(&ignoreUpdate).
		Run()

	if err != nil {
		return false, errors.Format(err.Error(), errors.ErrInternal)
	}

	return !ignoreUpdate, nil
}