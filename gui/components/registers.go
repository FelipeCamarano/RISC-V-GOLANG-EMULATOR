package components

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewRegisterGrid(bindings []binding.String) fyne.CanvasObject {
	grid := container.NewGridWithColumns(4)

	lastValues := make([]string, 32)

	for i := 0; i < 32; i++ {
		cell := createRegisterCell(i, bindings, lastValues)
		grid.Add(cell)
	}

	return container.NewPadded(grid)
}

func createRegisterCell(index int, bindings []binding.String, lastValues []string) fyne.CanvasObject {
	regName := getRegisterName(index)

	lblTitle := canvas.NewText(regName, theme.PrimaryColor())
	lblTitle.TextSize = 10
	lblTitle.Alignment = fyne.TextAlignCenter

	// Valor (Hex)
	lblValue := widget.NewLabel("") 
	lblValue.TextStyle = fyne.TextStyle{Monospace: true}
	lblValue.Alignment = fyne.TextAlignCenter

	defaultColor := color.RGBA{40, 40, 40, 255}
	flashColor := color.RGBA{0, 100, 0, 255}

	bg := canvas.NewRectangle(defaultColor)

	setupRegisterListener(index, bindings, lastValues, lblValue, bg, defaultColor, flashColor)

	content := container.NewVBox(lblTitle, lblValue)
	
	return container.NewStack(bg, container.NewPadded(content))
}

func getRegisterName(i int) string {
	// Nomes ABI padrão RISC-V
	names := []string{
		"zero", "ra", "sp", "gp", "tp", "t0", "t1", "t2",
		"s0/fp", "s1", "a0", "a1", "a2", "a3", "a4", "a5",
		"a6", "a7", "s2", "s3", "s4", "s5", "s6", "s7",
		"s8", "s9", "s10", "s11", "t3", "t4", "t5", "t6",
	}
	if i >= 0 && i < len(names) {
		return fmt.Sprintf("x%d (%s)", i, names[i])
	}
	return fmt.Sprintf("x%d", i)
}

func setupRegisterListener(index int, bindings []binding.String, lastValues []string, lblValue *widget.Label, bg *canvas.Rectangle, defaultColor, flashColor color.Color) {
	bindings[index].AddListener(binding.NewDataListener(func() {
		val, _ := bindings[index].Get()

		fyne.Do(func() {
			lblValue.SetText(val)

			if val != lastValues[index] && lastValues[index] != "" {
				flashBackground(bg, defaultColor, flashColor)
			}
			lastValues[index] = val
		})
	}))
}

func flashBackground(bg *canvas.Rectangle, defaultColor, flashColor color.Color) {
	bg.FillColor = flashColor
	bg.Refresh()

	go func() {
		time.Sleep(150 * time.Millisecond)
		fyne.Do(func() {
			bg.FillColor = defaultColor
			bg.Refresh()
		})
	}()
}