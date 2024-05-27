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

package utilities

import (
	"fmt"
	"os/exec"
	"runtime"

	"infinite-ugc-tool/pkg/modules/errors"

	"github.com/charmbracelet/huh/spinner"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) error {
	var err error
	spinner.New().Title("Attempting to open your browser...").Run()

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return errors.Format(err.Error(), errors.ErrInternal)
	}

	return nil
}
