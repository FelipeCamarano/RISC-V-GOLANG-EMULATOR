package components

import (
	"image/color"
	"strings"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

type Terminal struct {
	Container fyne.CanvasObject
	Output    *widget.TextGrid
	Buffer    strings.Builder
	Mb        *motherboard.Motherboard
}

func NewTerminal(mb *motherboard.Motherboard) *Terminal {
	t := &Terminal{
		Mb:     mb,
		Output: widget.NewTextGrid(),
	}

	t.Output.SetText("RISC-V OS Serial Terminal v1.0\n_")

	bg := canvas.NewRectangle(color.Black)
	stack := container.NewStack(bg, container.NewPadded(t.Output))
	
	t.Container = widget.NewCard("UART Terminal", "", stack)
	
	return t
}

func (t *Terminal) CheckOutput() {
	val := t.Mb.Bus.ReadByte(motherboard.IO_START)
	
	if val != 0 {
		fmt.Printf("Terminal: Lendo byte da porta IO: 0x%X\n", val)
		char := rune(val)
		
		fyne.Do(func() {
			t.Buffer.WriteRune(char)
			t.Output.SetText(t.Buffer.String() + "_")
		})
		
		t.Mb.Bus.WriteByte(motherboard.IO_START, 0)
	}
}