package modules

import (
	"fmt"
	"infinite-bookmarker/internal/shared/errors"
	"os/exec"
	"runtime"
)

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func OpenBrowser(url string) error {
	var err error

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
