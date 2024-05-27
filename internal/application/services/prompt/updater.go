// Copyright 2024 Alexis Bize
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompt_svc

import (
	"fmt"
	"infinite-ugc-tool/pkg/modules/errors"

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
