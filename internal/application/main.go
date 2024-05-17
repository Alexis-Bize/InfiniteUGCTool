package application

import (
	"fmt"
	"infinite-bookmarker/internal"
	view_authenticate "infinite-bookmarker/internal/application/views/authenticate"
	"infinite-bookmarker/internal/shared/modules/helpers/credentials"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Start() {
	currentCredentials, err := credentials.GetOrCreateCredentials(credentials.Credentials{})
	if err != nil {
		log.Fatal(err)
	}

	hasUserCredentials := currentCredentials.User != (credentials.UserCredentials{})
	fmt.Println(hasUserCredentials)

	a := app.New()
	w := a.NewWindow(internal.GetConfig().Title + " (" + internal.GetConfig().Version + ")")

	w.Resize(fyne.Size{
		Width: 720,
		Height: 440,
	})

	w.SetPadded(true)
	w.SetFixedSize(true)
	w.CenterOnScreen()

	view_authenticate.Render(w)

	w.ShowAndRun()
}