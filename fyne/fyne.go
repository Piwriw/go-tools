package fyne

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	_ "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
go get fyne.io/fyne/v2
go get -u github.com/flopp/go-findfont 支持中文

go install fyne.io/fyne/v2/cmd/fyne@latest
*/
func main1() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Form")

	nameEntry := widget.NewEntry()
	passEntry := widget.NewPasswordEntry()

	form := widget.NewForm(
		&widget.FormItem{Text: "Name", Widget: nameEntry},
		&widget.FormItem{Text: "Pass", Widget: passEntry},
	)

	form.OnSubmit = func() {
		fmt.Println("name:", nameEntry.Text, "pass:", passEntry.Text, "login in")
	}
	form.OnCancel = func() {
		fmt.Println("login canceled")
	}

	myWindow.SetContent(form)
	myWindow.Resize(fyne.NewSize(150, 150))
	myWindow.ShowAndRun()
}
