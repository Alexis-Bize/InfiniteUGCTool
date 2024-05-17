package view_authenticate

import (
	"infinite-bookmarker/internal"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Render(w fyne.Window) {
	hello := widget.NewLabel(internal.GetConfig().Title)
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Authenticate", func() {
			hello.SetText("Welcome :)")
		}),
	))
}