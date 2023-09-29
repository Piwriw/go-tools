package gui

import (
	"fmt"
	"fyne-ict/score"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// App 构建GUI 画图
func App(rows [][]string) {
	myApp := app.New()
	myApp.Settings().SetTheme(&MyTheme{})
	myWindow := myApp.NewWindow("SGU ICT 学生绩点计算工具")

	list := widget.NewTable(
		func() (int, int) {
			return len(rows), len(rows[0])
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("wide content")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(rows[i.Row][i.Col])
		})
	result := widget.NewLabel("学生绩点平均值:")
	content := widget.NewButton("开始计算", func() {
		avscore := score.DoWork()
		result.Text = fmt.Sprintf("学生绩点平均值:%f", avscore)
		result.Refresh()
	})
	contentLay := container.New(layout.NewHBoxLayout(), content, result)
	contentLay.Resize(fyne.NewSize(500, 30))

	grid := container.New(layout.NewGridLayout(1), list, contentLay)

	myWindow.SetContent(grid)
	myWindow.Resize(fyne.NewSize(500, 500))
	myWindow.ShowAndRun()
}
