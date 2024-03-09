package main

import (
	"context"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

var selectedItemIndex = -1

func ReadClipboard(clipboardUpdates chan string, ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	clipboardChannel := clipboard.Watch(ctx, clipboard.FmtText)
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-clipboardChannel:
			clipboardUpdates <- string(item)
		}
	}
}

func trimClipboardItem(item string) string {
	if len(item) > 40 {
		// remove newlines
		item = strings.Replace(item, "\n", "", -1)
		return item[:40] + "..."
	}
	return item
}

func removeItemFromList(clipboardListData []string, index int) []string {
	if index >= 0 && index < len(clipboardListData) {
		return append((clipboardListData)[:index], (clipboardListData)[index+1:]...)
	}
	return clipboardListData
}

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	app := app.NewWithID("com.github.olksndrdevhub.copylog")
	window := app.NewWindow("CopyLog")
	// icon, error := os.ReadFile("assets/icon.png")
	// if error != nil {
	// 	println("Error reading icon file.", error)
	// }
	// iconResource := fyne.NewStaticResource("icon.png", icon)
	// window.SetIcon(iconResource)
	if desk, ok := app.(desktop.App); ok {
		menu := fyne.NewMenu("CopyLog",
			fyne.NewMenuItem("Show", func() {
				window.Show()
			}))
		desk.SetSystemTrayMenu(menu)
	}
	window.SetCloseIntercept(func() {
		window.Hide()
	})
	// menuItem1 := fyne.NewMenuItem("Show", nil)
	// menuItem2 := fyne.NewMenuItem("About", nil)
	//
	// newMenu1 := fyne.NewMenu("File", menuItem1, menuItem2)
	//
	// menu := fyne.NewMainMenu(newMenu1)
	//
	// window.SetMainMenu(menu)

	welcomeText := widget.NewLabel("Welcome to  CopyLog!")

	clipboardUpdates := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	var wait_group sync.WaitGroup
	wait_group.Add(1)
	go ReadClipboard(clipboardUpdates, ctx, &wait_group)

	clipboardListData := make([]string, 0)

	clipboardItemsList := widget.NewList(
		func() int {
			return len(clipboardListData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Clipboard Items")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			obj.(*widget.Label).SetText(trimClipboardItem(clipboardListData[id]))
		},
	)
	clipboardItemsList.Resize(fyne.NewSize(300, 0))

	itemDisplayEntry := widget.NewMultiLineEntry()
	itemDisplayEntry.PlaceHolder = "Selected item will be displayed here..."
	itemDisplayEntry.TextStyle.Symbol = true
	itemDisplayEntry.TextStyle.Monospace = true

	removeItemButton := widget.NewButton("Remove", func() {
		if len(clipboardListData) == 0 {
			return
		}
		if selectedItemIndex == -1 {
			return
		}
		clipboardListData = removeItemFromList(clipboardListData, selectedItemIndex)
		clipboardItemsList.Refresh()
		itemDisplayEntry.SetText("")
	})

	useItemButton := widget.NewButton("Use", func() {
		clipboard.Write(clipboard.FmtText, []byte(itemDisplayEntry.Text))
	})

	itemActionsLayout := container.NewHBox(
		removeItemButton,
		useItemButton,
	)
	itemActionsLayout.Hide()

	clipboardItemsList.OnSelected = func(id widget.ListItemID) {
		itemDisplayEntry.SetText(clipboardListData[id])
		itemActionsLayout.Show()
		selectedItemIndex = int(id)
	}

	go func() {
		for {
			select {
			case text := <-clipboardUpdates:
				clipboardListData = append(clipboardListData, text)
				clipboardItemsList.Refresh()
			case <-ctx.Done():
				return
			}
		}
	}()

	splitLayout := container.NewHSplit(clipboardItemsList, itemDisplayEntry)
	splitLayout.SetOffset(0.3)

	window.SetContent(container.NewBorder(
		container.NewHBox(
			welcomeText, layout.NewSpacer(), itemActionsLayout,
		),
		nil, nil, nil,
		splitLayout,
	))

	window.Resize(fyne.NewSize(640, 460))
	window.ShowAndRun()
	cancel()
	wait_group.Wait()
}
