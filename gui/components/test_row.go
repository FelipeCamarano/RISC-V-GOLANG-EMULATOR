package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TestRow struct {
	Container *fyne.Container
	
	icon  *widget.Icon
	label *canvas.Text
	
	Name    string
	Details string
}

func NewTestRow(name string, parentWindow fyne.Window, onDebug func()) *TestRow {
	row := &TestRow{Name: name}

	row.icon = widget.NewIcon(theme.MediaRecordIcon())
	row.label = canvas.NewText(name, theme.ForegroundColor())
	row.label.TextSize = 14

	btnDetails := row.createDetailsButton(parentWindow)
	btnDebug := row.createDebugButton(onDebug)

	rightButtons := container.NewHBox(btnDebug, btnDetails)
	row.Container = container.NewBorder(nil, nil, row.icon, rightButtons, row.label)

	return row
}

func (t *TestRow) createDetailsButton(parentWindow fyne.Window) *widget.Button {
	return widget.NewButtonWithIcon("", theme.InfoIcon(), func() {
		content := container.NewVBox(
			widget.NewLabelWithStyle("Detalhes: "+t.Name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			widget.NewLabel(t.Details),
			widget.NewButton("Fechar", func() {
				if len(parentWindow.Canvas().Overlays().List()) > 0 {
					parentWindow.Canvas().Overlays().Top().Hide()
				}
			}),
		)
		widget.NewModalPopUp(content, parentWindow.Canvas()).Show()
	})
}

func (t *TestRow) createDebugButton(onDebug func()) *widget.Button {
	return widget.NewButtonWithIcon("", theme.SearchIcon(), onDebug)
}


func (t *TestRow) SetRunning() {
	t.icon.SetResource(theme.MediaPlayIcon())
	t.label.Color = theme.ForegroundColor()
	t.label.TextStyle = fyne.TextStyle{Bold: true}
	t.Refresh()
}

func (t *TestRow) SetPass(cycles int, result int32) {
	t.icon.SetResource(theme.ConfirmIcon())
	t.label.Color = color.RGBA{0, 180, 0, 255}
	t.label.TextStyle = fyne.TextStyle{Bold: false}
	t.Details = fmt.Sprintf("PASS\nCiclos: %d\nResult (x10): %d", cycles, result)
	t.Refresh()
}

func (t *TestRow) SetFail(msg string, cycles int, result int32, testCase int32) {
	t.icon.SetResource(theme.CancelIcon())
	t.label.Color = color.RGBA{200, 0, 0, 255}
	t.label.TextStyle = fyne.TextStyle{Bold: false}
	t.Details = fmt.Sprintf("FAIL: %s\nCiclos: %d\nVal: %d | Case: %d", msg, cycles, result, testCase)
	t.Refresh()
}

func (t *TestRow) SetPending() {
	t.icon.SetResource(theme.MediaRecordIcon())
	t.label.Color = theme.ForegroundColor()
	t.Details = "Aguardando execução..."
	t.Refresh()
}

func (t *TestRow) Refresh() {
	t.icon.Refresh()
	t.label.Refresh()
}