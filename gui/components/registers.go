package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewRegisterGrid(bindings []binding.String) fyne.CanvasObject {
	grid := container.NewGridWithColumns(4)

	for i := 0; i < 32; i++ {
		regName := fmt.Sprintf("x%d", i)
		if i == 2 { regName = "sp" }
		if i == 1 { regName = "ra" }
		
		lblTitle := canvas.NewText(regName, theme.PrimaryColor())
		lblTitle.TextSize = 10
		lblTitle.Alignment = fyne.TextAlignCenter

		lblValue := widget.NewLabelWithData(bindings[i])
		lblValue.TextStyle = fyne.TextStyle{Monospace: true}
		lblValue.Alignment = fyne.TextAlignCenter
		
		box := container.NewVBox(
			lblTitle,
			lblValue,
		)
		
		bg := canvas.NewRectangle(color.RGBA{50, 50, 50, 255})
		cell := container.NewStack(bg, container.NewPadded(box))

		grid.Add(cell)
	}

	return container.NewVScroll(container.NewPadded(grid))
}