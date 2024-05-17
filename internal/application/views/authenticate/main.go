package view_authenticate

import (
	"infinite-bookmarker/internal/application/components"
	services_auth "infinite-bookmarker/internal/services/auth"
	"infinite-bookmarker/internal/shared/modules/helpers/credentials"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Render(w fyne.Window) {
	email := widget.NewEntry()
	password := widget.NewPasswordEntry()
	form := &widget.Form{}

	form = &widget.Form{
		Items: []*widget.FormItem{
			{
				Text: "Email",
				Widget: email,
			},
			{
				Text: "Password",
				Widget: password,
			},
		},
		SubmitText: "Authenticate",
	}

	form.OnSubmit = func() {
		form.Disable()
		spartanToken, err := services_auth.Authenticate(email.Text, password.Text)

		if err != nil {
			form.Enable()
			return
		}

		currentTime := time.Now()
		duration := 3*time.Hour + 30*time.Minute
		expiration := currentTime.Add(duration)

		credentials.SaveCredentials(credentials.Credentials{
			User: credentials.UserCredentials{
				Email: email.Text,
				Password: password.Text,
			},
			SpartanToken: credentials.SpartanTokenCredentials{
				Value: spartanToken,
				Expiration: expiration.Format(time.RFC3339),
			},
		})
	}

	w.SetContent(
		container.New(
			layout.NewVBoxLayout(),
			components.RenderHero("Authenticate"),
			form,
		),
	)
}
