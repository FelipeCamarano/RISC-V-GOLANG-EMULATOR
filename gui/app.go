package gui

import (
	"fmt"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/utils"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/windows"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

type EmulatorApp struct {
	FyneApp    fyne.App
	Mb         *motherboard.Motherboard
	MainWindow *windows.MainWindow
}

func New(mb *motherboard.Motherboard) *EmulatorApp {
	a := app.New()
	emu := &EmulatorApp{
		FyneApp: a,
		Mb:      mb,
	}
	emu.MainWindow = windows.NewMainWindow(a, mb)
	return emu
}

func (e *EmulatorApp) Run() {
	openTestsCallback := func() {
		tWindow := windows.NewTestRunnerWindow(e.FyneApp, e.LoadAndDebug)
		tWindow.Show()
	}

	e.MainWindow.Build(openTestsCallback)
	e.MainWindow.Show()
	e.FyneApp.Run()
}

func (e *EmulatorApp) LoadAndDebug(path string) {
	fmt.Println("Tentando carregar debug para:", path)

	e.MainWindow.IsRunning = false
	e.resetCPU()

	startPC, err := utils.LoadHexFile(path, e.Mb)
	
	if err != nil {
		dialog.ShowError(fmt.Errorf("Falha ao carregar HEX:\n%s", err.Error()), e.MainWindow.Window)
		fmt.Println("Erro no Load:", err)
		return
	}

	e.Mb.CPU.SetPC(startPC)
	
	e.MainWindow.UpdateData()
	
	filename := filepath.Base(path)
	e.MainWindow.Window.SetTitle("RISC-V Monitor - Debugging: " + filename)
	
	fmt.Printf("Sucesso! PC inicial: 0x%08X\n", startPC)
}

func (e *EmulatorApp) resetCPU() {
	e.Mb.CPU.PC = 0
	for i := range e.Mb.CPU.Registers {
		e.Mb.CPU.Registers[i] = 0
	}
}