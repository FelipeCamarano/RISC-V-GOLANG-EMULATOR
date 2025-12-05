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

	// Vamos mostrar até o endereço 0x000A0000 (655KB)
	// Isso cobre:
	// 0x00000 -> RAM
	// 0x80000 -> VRAM
	// 0x90000 -> Cartucho
	// 0x9FC00 -> IO
	const MaxViewAddress = 0x000A0000
	rows := MaxViewAddress / 16

	rv.Table = widget.NewTable(
		func() (int, int) { return rows, 17 },
		func() fyne.CanvasObject {
			text := canvas.NewText("00", color.White)
			text.Alignment = fyne.TextAlignCenter
			text.TextSize = 12
			return container.NewPadded(text)
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			text := cell.(*fyne.Container).Objects[0].(*canvas.Text)
			
			// Coluna de Endereço
			if id.Col == 0 {
				addr := uint32(id.Row * 16)
				text.Text = fmt.Sprintf("0x%05X", addr)
				text.Color = theme.PrimaryColor()
				text.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				// Colunas de Dados (Lê do BARRAMENTO agora!)
				addr := uint32(id.Row*16 + (id.Col - 1))
				val := mb.Bus.ReadByte(addr) // <--- AQUI A MÁGICA
				
				text.Text = fmt.Sprintf("%02X", val)
				text.TextStyle = fyne.TextStyle{Monospace: true}

				// --- Color Coding por Região ---
				if addr >= motherboard.IO_START {
					text.Color = color.RGBA{255, 165, 0, 255} // Laranja (IO)
				} else if addr >= motherboard.CARTRIDGE_START {
					text.Color = color.RGBA{100, 200, 255, 255} // Azul Claro (Cartucho)
				} else if addr >= motherboard.VRAM_START {
					text.Color = color.RGBA{100, 255, 100, 255} // Verde (VRAM)
				} else {
					// RAM Normal
					if val == 0 {
						text.Color = color.Gray{Y: 60}
					} else {
						text.Color = color.White
					}
				}
			}
			text.Refresh()
		},
	)

	// Ajustes de largura
	rv.Table.SetColumnWidth(0, 90)
	for i := 1; i <= 16; i++ {
		rv.Table.SetColumnWidth(i, 32)
	}

	rv.Container = widget.NewCard("System Memory Map (Bus View)", "", container.NewPadded(rv.Table))
	return rv
}

func (rv *RamViewer) Refresh() {
	fyne.Do(func() {
		rv.Table.Refresh()
	})
}