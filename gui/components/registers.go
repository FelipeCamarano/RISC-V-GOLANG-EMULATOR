package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func NewRegisterGrid(data []binding.String) fyne.CanvasObject {
	regWidgets := make([]fyne.CanvasObject, 32)
	for i := 0; i < 32; i++ {
		regWidgets[i] = widget.NewLabelWithData(data[i])
	}
	
	grid := container.NewGridWithColumns(4, regWidgets...)
	
	return widget.NewCard("Registradores (x0-x31)", "", container.NewVScroll(grid))
}