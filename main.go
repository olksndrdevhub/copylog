package main

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

func ReadClipboard() {
	clipboard_channel := clipboard.Watch(context.TODO(), clipboard.FmtText)
	for item := range clipboard_channel {
		//  print clipboard content
		println("New item in clipboard: ", string(item))
	}
}

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	app := app.New()
	window := app.NewWindow("CopyLog")

	hello := widget.NewLabel("Welcome to  CopyLog!")

	name_input := widget.NewEntry()
  name_input.SetPlaceHolder("Enter your name...")
	name_input.Resize(fyne.NewSize(200, 30))
	name_input.Move(fyne.NewPos(0, 60))

	hello_btn := widget.NewButton("Submit", func() {
		println("Hello ", name_input.Text)
		hello.SetText("Hello " + name_input.Text + "!")
	})
	hello_btn.Resize(fyne.NewSize(100, 20))
	hello_btn.Move(fyne.NewPos(0, 120))

	window.SetContent(container.NewWithoutLayout(
		hello,
		name_input,
		hello_btn,
	))
	window.Resize(fyne.NewSize(600, 300))
	window.ShowAndRun()
	ReadClipboard()
}
