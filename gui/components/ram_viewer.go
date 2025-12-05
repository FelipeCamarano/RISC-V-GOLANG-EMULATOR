package components

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

type RamViewer struct {
	Container fyne.CanvasObject
	Table     *widget.Table
	Mb        *motherboard.Motherboard
}

func NewRamViewer(mb *motherboard.Motherboard) *RamViewer {
	rv := &RamViewer{Mb: mb}

	const MaxViewAddress = motherboard.IO_LIMIT
	rows := MaxViewAddress / 16

	rv.Table = widget.NewTable(
		func() (int, int) { return rows, 17 },
		createCellTemplate,
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			rv.updateCell(id, cell)
		},
	)

	configureTableColumns(rv.Table)

	rv.Container = widget.NewCard("Memory Map", "", container.NewBorder(nil, nil, nil, nil,
		container.NewBorder(createLegend(), nil, nil, nil, rv.Table),
	))

	return rv
}

func (rv *RamViewer) Refresh() {
	fyne.Do(func() {
		rv.Table.Refresh()
	})
}

func createCellTemplate() fyne.CanvasObject {
	text := canvas.NewText("00", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextSize = 11
	bg := canvas.NewRectangle(color.Transparent)
	return container.NewStack(bg, container.NewPadded(text))
}

func (rv *RamViewer) updateCell(id widget.TableCellID, cell fyne.CanvasObject) {
	stack := cell.(*fyne.Container)
	bg := stack.Objects[0].(*canvas.Rectangle)
	text := stack.Objects[1].(*fyne.Container).Objects[0].(*canvas.Text)

	applyZebraStripe(id.Row, bg)

	if id.Col == 0 {
		updateAddressCell(id.Row, text)
	} else {
		rv.updateDataCell(id.Row, id.Col, text)
	}

	bg.Refresh()
	text.Refresh()
}

func applyZebraStripe(row int, bg *canvas.Rectangle) {
	if row%2 == 0 {
		bg.FillColor = color.RGBA{30, 30, 30, 255}
	} else {
		bg.FillColor = color.Transparent
	}
}

func updateAddressCell(row int, text *canvas.Text) {
	addr := uint32(row * 16)
	text.Text = fmt.Sprintf("0x%05X", addr)
	text.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}

	if addr >= motherboard.IO_START {
		text.Color = color.RGBA{255, 165, 0, 255} // (IO)
	} else if addr >= motherboard.CARTRIDGE_START {
		text.Color = color.RGBA{50, 150, 255, 255} // (Cartucho)
	} else if addr >= motherboard.VRAM_START {
		text.Color = color.RGBA{100, 255, 100, 255} // (VRAM)
	} else {
		text.Color = theme.PrimaryColor() // RAM
	}
}

func (rv *RamViewer) updateDataCell(row, col int, text *canvas.Text) {
	byteOffset := uint32(col - 1)
	addr := uint32(row*16) + byteOffset
	val := rv.Mb.Bus.ReadByte(addr)

	text.Text = fmt.Sprintf("%02X", val)
	text.TextStyle = fyne.TextStyle{Monospace: true}

	if val == 0 {
		text.Color = color.Gray{Y: 60}
		text.TextStyle.Bold = false
	} else {
		text.Color = color.White
		text.TextStyle.Bold = true

		if addr >= motherboard.VRAM_START && addr < motherboard.CARTRIDGE_START {
			text.Color = color.RGBA{150, 255, 150, 255}
		}
	}
}

func configureTableColumns(table *widget.Table) {
	table.SetColumnWidth(0, 80)
	for i := 1; i <= 16; i++ {
		table.SetColumnWidth(i, 28)
	}
}

func createLegend() fyne.CanvasObject {
	return container.NewHBox(
		canvas.NewText("RAM", theme.PrimaryColor()),
		canvas.NewText("|", color.Gray{Y: 100}),
		canvas.NewText("VRAM", color.RGBA{100, 255, 100, 255}),
		canvas.NewText("|", color.Gray{Y: 100}),
		canvas.NewText("CART", color.RGBA{50, 150, 255, 255}),
		canvas.NewText("|", color.Gray{Y: 100}),
		canvas.NewText("IO", color.RGBA{255, 165, 0, 255}),
	)
}