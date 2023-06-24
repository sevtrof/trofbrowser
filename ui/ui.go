package ui

import (
	"log"
	"trofbrowser/browser"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fogleman/gg"
	"rogchap.com/v8go"
)

func RunUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Trofbrowser")

	tabs := container.NewAppTabs()

	brwsr := browser.NewBrowser()

	newTabButton := widget.NewButton("New Tab", func() {
		tabContent := newTabContent(myWindow, tabs, brwsr)
		tabs.Append(container.NewTabItem("New Tab", tabContent))
	})

	tabContent := newTabContent(myWindow, tabs, brwsr)
	tabs.Append(container.NewTabItem("+", newTabButton))
	tabs.Append(container.NewTabItem("New Tab", tabContent))

	myWindow.SetContent(tabs)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}

func newTabContent(window fyne.Window, tabs *container.AppTabs, browser browser.Browser) fyne.CanvasObject {
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter URL...")

	dc := gg.NewContext(800, 600)
	js := v8go.NewContext()

	img := canvas.NewImageFromImage(dc.Image())
	scrollContainer := container.NewScroll(img)

	urlEntry.OnSubmitted = func(url string) {
		domTree, err := browser.NavigateToURL(url)
		if err != nil {
			log.Printf("Error fetching webpage: %v", err)
			return
		}

		browser.RenderDOMTree(domTree, 10, 20, dc, js, url)

		img.Image = dc.Image()
		img.Refresh()
	}

	return container.NewVBox(urlEntry, scrollContainer)
}
