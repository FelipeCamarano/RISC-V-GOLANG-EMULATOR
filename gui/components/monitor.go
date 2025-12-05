package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

type MonitorWrapper struct {
	Container fyne.CanvasObject
	Raster    *canvas.Raster
	Mb        *motherboard.Motherboard
}

func NewMonitor(mb *motherboard.Motherboard) *MonitorWrapper {
	m := &MonitorWrapper{
		Mb: mb,
	}


	m.Raster = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		if w == 0 || h == 0 {
			return color.Black
		}

		vRamX := int((float64(x) / float64(w)) * 320.0)
		vRamY := int((float64(y) / float64(h)) * 200.0)

		if vRamX < 0 { vRamX = 0 }
		if vRamX >= 320 { vRamX = 319 }
		if vRamY < 0 { vRamY = 0 }
		if vRamY >= 200 { vRamY = 199 }

		index := uint32(vRamY*320 + vRamX)
		
		if index >= m.Mb.VRAM.Size() {
			return color.Black
		}
		
		pixelVal := m.Mb.VRAM.ReadByte(index)

		return color.RGBA{0, uint8(pixelVal), 0, 255}
	})

	m.Raster.SetMinSize(fyne.NewSize(640, 400))

	m.Container = container.NewPadded(
		widget.NewCard("VRAM Monitor", "", container.NewPadded(m.Raster)),
	)

	return m
}

func (m *MonitorWrapper) Refresh() {
	fyne.Do(func() {
		m.Raster.Refresh()
	})
}