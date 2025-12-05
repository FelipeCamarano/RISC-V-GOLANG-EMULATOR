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
	a := app.NewWithID("com.dainslash.riscv.emulator")
	
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
	
	loadCartridgeCallback := func(path string) {
		e.LoadCartridge(path)
	}

	e.MainWindow.Build(openTestsCallback, loadCartridgeCallback)
	e.MainWindow.Show()
	e.FyneApp.Run()
}

func (e *EmulatorApp) LoadAndDebug(path string) {
	fmt.Println("Operador: Injetando Teste de Debug ->", path)
	e.MainWindow.IsRunning = false
	
	e.Mb.Reset()

	startPC, err := utils.LoadHexFile(path, e.Mb)
	
	if err != nil {
		dialog.ShowError(fmt.Errorf("Erro ao injetar HEX:\n%s", err.Error()), e.MainWindow.Window)
		return
	}
	
	e.Mb.CPU.SetPC(startPC)
	e.MainWindow.UpdateData()
	
	filename := filepath.Base(path)
	e.MainWindow.Window.SetTitle("RISC-V Monitor - Debugging: " + filename)
	
	fmt.Printf("Sucesso! Debugger conectado. PC: 0x%08X\n", startPC)
}

func (e *EmulatorApp) LoadCartridge(path string) {
	fmt.Println("Operador: Inserindo Cartucho ->", path)

	if err := e.Mb.InsertCartridge(path); err != nil {
		dialog.ShowError(err, e.MainWindow.Window)
		return
	}

	e.Mb.Reset()
	
	e.MainWindow.UpdateData()
	fmt.Println("Sistema Resetado. Bootloader (BIOS) iniciado.")
}