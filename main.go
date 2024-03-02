package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.New()
	window := app.NewWindow("CopyLog")

	hello := widget.NewLabel("Hi CopyLog!")
	window.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hello", func() {
			hello.SetText("Welcome :)")
		}),
	))

	window.ShowAndRun()

}
