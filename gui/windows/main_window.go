package windows

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/components"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

type MainWindow struct {
	Window    fyne.Window
	App       fyne.App
	Mb        *motherboard.Motherboard

	PcBinding   binding.String
	RegBindings []binding.String

	IsRunning bool
}

func NewMainWindow(a fyne.App, mb *motherboard.Motherboard) *MainWindow {
	w := &MainWindow{
		App:         a,
		Window:      a.NewWindow("RISC-V Emulator Monitor"),
		Mb:          mb,
		PcBinding:   binding.NewString(),
		RegBindings: make([]binding.String, 32),
	}

	w.Window.Resize(fyne.NewSize(800, 600))
	for i := 0; i < 32; i++ {
		w.RegBindings[i] = binding.NewString()
	}

	return w
}

func (mw *MainWindow) Build(onOpenTests func()) {
	// Cria grid de registradores usando o componente reutilizável
	regGroup := components.NewRegisterGrid(mw.RegBindings)

	pcLabel := widget.NewLabelWithData(mw.PcBinding)

	btnStep := widget.NewButton("Step", func() { mw.stepCPU() })

	btnRun := widget.NewButton("Auto Run", nil)
	btnRun.OnTapped = func() {
		if mw.IsRunning {
			mw.IsRunning = false
			btnRun.SetText("Auto Run")
			return
		}
		mw.IsRunning = true
		btnRun.SetText("Stop")
		go mw.runLoop()
	}

	btnTests := widget.NewButtonWithIcon("Test Suite", theme.ListIcon(), onOpenTests)

	topPanel := container.NewVBox(
		widget.NewCard("Status", "", pcLabel),
		container.NewHBox(btnStep, btnRun, btnTests),
	)

	content := container.NewBorder(topPanel, nil, nil, nil, regGroup)
	mw.Window.SetContent(content)
}

func (mw *MainWindow) Show() {
	mw.UpdateData()
	mw.Window.Show()
}

func (mw *MainWindow) UpdateData() {
	mw.PcBinding.Set(fmt.Sprintf("PC: 0x%08X", mw.Mb.CPU.PC))
	for i, val := range mw.Mb.CPU.Registers {
		mw.RegBindings[i].Set(fmt.Sprintf("x%02d: %d (0x%X)", i, val, val))
	}
}

func (mw *MainWindow) stepCPU() {
	pc := uint32(mw.Mb.CPU.PC)
	if mw.Mb.Bus.ReadWord(pc) == 0 {
		mw.IsRunning = false
		mw.PcBinding.Set(fmt.Sprintf("PC: 0x%08X (HALT)", pc))
		return
	}
	mw.Mb.CPU.Step()
	mw.UpdateData()
}

func (mw *MainWindow) runLoop() {
	for mw.IsRunning {
		mw.stepCPU()
		time.Sleep(10 * time.Millisecond)
	}
}