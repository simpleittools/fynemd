package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	d2 "github.com/sqweek/dialog"
	"os"
	"strings"
)

type config struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	SaveMenuItem  *fyne.MenuItem
}

var cfg config

func main() {
	// create a fyne app
	a := app.New()

	// use custom theme
	a.Settings().SetTheme(&myTheme{})

	// create a window for the app
	win := a.NewWindow("Markdown")

	// get the user interface
	edit, preview := cfg.makeUI()
	cfg.createMenuItems(win)

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 500})
	win.CenterOnScreen()
	win.Canvas().Focus(edit)
	win.ShowAndRun()
}

// makeUI will create the user interface with the entry and rich text widgets
func (app *config) makeUI() (*widget.Entry, fyne.CanvasObject) {
	// add a new multi line input widget
	edit := widget.NewMultiLineEntry()
	edit.Wrapping = fyne.TextWrapWord

	// create a preview pane that allows for Rich Text generated by Markdown
	preview := widget.NewRichTextFromMarkdown("")
	preview.Wrapping = fyne.TextWrapWord

	previewContainer := container.NewScroll(preview)

	app.EditWidget = edit
	app.PreviewWidget = preview

	// when editing with Markdown in the Edit window, automatically show the output in the Preview pane as Rich Text
	edit.OnChanged = preview.ParseMarkdown

	return edit, previewContainer
}

// createMenuItems will create the dropdown menus
func (app *config) createMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true

	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))

	//showPreviewPane := fyne.NewMenuItem("Show Preview Pane", func() {})

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)
	//viewMenu := fyne.NewMenu("View", showPreviewPane)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)

}

func (app *config) saveFunc(win fyne.Window) func() {
	return func() {
		if app.CurrentFile != nil {
			write, err := storage.Writer(app.CurrentFile)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		filePath, err := d2.File().Filter("Markdown Files", "md").Title("Save File As").Save()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if filePath == "" {
			return
		}

		if !strings.HasSuffix(strings.ToLower(filePath), ".md") {
			filePath = filePath + ".md" // automatically add the extension of .md
		}

		err = os.WriteFile(filePath, []byte(app.EditWidget.Text), 0644)
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		app.CurrentFile, _ = storage.ParseURI("file://" + filePath)
		win.SetTitle(win.Title() + " - " + filePath)
		app.SaveMenuItem.Disabled = false
	}
}

func (app *config) openFunc(win fyne.Window) func() {
	return func() {
		filePath, err := d2.File().Filter("Markdown Files", "md").Load()
		if err != nil {
			dialog.ShowError(err, win)
			return
		}

		if filePath == "" {
			return
		}

		data, err := os.ReadFile(filePath)
		if err != nil {
			dialog.ShowError(err, win)
		}

		app.EditWidget.SetText(string(data))

		app.CurrentFile, _ = storage.ParseURI("file://" + filePath)

		win.SetTitle(win.Title() + " - " + filePath)
		app.SaveMenuItem.Disabled = false

	}
}
