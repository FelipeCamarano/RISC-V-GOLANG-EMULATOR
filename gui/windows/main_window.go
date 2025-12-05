package windows

import (
	"fmt"
	"time"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
	StatusBinding binding.String
	RegBindings []binding.String
	
	Monitor     *components.MonitorWrapper 
	SpeedCtrl   *components.SpeedControl
	RamView     *components.RamViewer
	Terminal    *components.Terminal

	IsRunning bool
}

func NewMainWindow(a fyne.App, mb *motherboard.Motherboard) *MainWindow {
	w := &MainWindow{
		App:         a,
		Window:      a.NewWindow("RISC-V Workstation"),
		Mb:          mb,
		PcBinding:   binding.NewString(),
		StatusBinding: binding.NewString(),
		RegBindings: make([]binding.String, 32),
	}

	w.Window.Resize(fyne.NewSize(1280, 850))
	w.StatusBinding.Set("System Ready. No cartridge loaded.")
	for i := 0; i < 32; i++ {
		w.RegBindings[i] = binding.NewString()
	}

	return w
}


func (mw *MainWindow) Build(onOpenTests func(), onLoadCartridge func(string)) {
	// 1. Componentes
	regGroup := components.NewRegisterGrid(mw.RegBindings)
	mw.Monitor = components.NewMonitor(mw.Mb)
	mw.SpeedCtrl = components.NewSpeedControl()
	mw.RamView = components.NewRamViewer(mw.Mb)
	mw.Terminal = components.NewTerminal(mw.Mb)

	// 2. Toolbar (Passamos os callbacks para ela)
	toolbar := mw.createToolbar(onOpenTests, onLoadCartridge)

	// 3. Layout (Split View)
	statusBar := widget.NewLabelWithData(mw.StatusBinding)
	statusBar.TextStyle = fyne.TextStyle{Monospace: true}

	leftSplit := container.NewVSplit(mw.Monitor.Container, mw.Terminal.Container)
	leftSplit.SetOffset(0.65)

	rightSplit := container.NewVSplit(regGroup, mw.RamView.Container)
	rightSplit.SetOffset(0.4)

	mainSplit := container.NewHSplit(leftSplit, rightSplit)
	mainSplit.SetOffset(0.6)

	centralPanel := container.NewBorder(nil, statusBar, nil, nil, mainSplit)

	content := container.NewBorder(toolbar, nil, nil, nil, centralPanel)
	mw.Window.SetContent(content)
}

func (mw *MainWindow) SetStatus(msg string) {
	mw.StatusBinding.Set(msg)
}

func (mw *MainWindow) createToolbar(onOpenTests func(), onLoadCartridge func(string)) fyne.CanvasObject {
	pcLabel := widget.NewLabelWithData(mw.PcBinding)
	pcLabel.TextStyle = fyne.TextStyle{Monospace: true, Bold: true}

	btnLoadCart := widget.NewButtonWithIcon("Insert Cartridge", theme.FolderOpenIcon(), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				path := reader.URI().Path()
				onLoadCartridge(path)
				mw.SetStatus("Loaded Cartridge: " + filepath.Base(path))
			}
		}, mw.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bin", ".rom", ".img"}))
		fd.Show()
	})

	btnStep := widget.NewButtonWithIcon("Step", theme.MediaSkipNextIcon(), func() { 
		mw.stepCPU() 
		mw.UpdateData() 
	})
	
	btnRun := widget.NewButtonWithIcon("Run", theme.MediaPlayIcon(), nil)
	btnRun.OnTapped = func() {
		if mw.IsRunning {
			mw.IsRunning = false
			btnRun.SetText("Run")
			btnRun.SetIcon(theme.MediaPlayIcon())
			mw.SetStatus("System Paused")
			return
		}
		mw.IsRunning = true
		btnRun.SetText("Stop")
		btnRun.SetIcon(theme.MediaStopIcon())
		mw.SetStatus("System Running...")
		go mw.runLoop()
	}

	btnTests := widget.NewButtonWithIcon("Debug", theme.ListIcon(), onOpenTests)

	toolbar := container.NewBorder(
		nil, nil, 
		container.NewHBox(widget.NewIcon(theme.ComputerIcon()), pcLabel), 
		container.NewHBox(btnLoadCart, btnStep, btnRun, btnTests),
		container.NewPadded(mw.SpeedCtrl.Container),
	)

	return container.NewVBox(toolbar, widget.NewSeparator())
}

func (mw *MainWindow) Show() {
	mw.UpdateData()
	mw.Window.Show()
}

func (mw *MainWindow) stepCPU() {
	pc := uint32(mw.Mb.CPU.PC)
	if pc == 0 && mw.Mb.Bus.ReadWord(0) == 0 {
		mw.IsRunning = false
		return
	}
	mw.Mb.CPU.Step()
	mw.UpdateData()
}

func (mw *MainWindow) runLoop() {
	for mw.IsRunning {
		delay := mw.SpeedCtrl.GetDelay()
		if delay > 0 {
			mw.stepCPU()
			mw.Terminal.CheckOutput()
			time.Sleep(delay)
		} else {
			for i:=0; i<1000; i++ {
				mw.Mb.CPU.Step()
				if i%100 == 0 { mw.Terminal.CheckOutput() }
			}
			mw.UpdateData()
			time.Sleep(1 * time.Microsecond)
		}
	}
}

func (mw *MainWindow) UpdateData() {
	mw.PcBinding.Set(fmt.Sprintf("PC: 0x%08X", mw.Mb.CPU.PC))
	for i, val := range mw.Mb.CPU.Registers {
		mw.RegBindings[i].Set(fmt.Sprintf("%08X", val))
	}
	if mw.Mb.VRAM.Dirty {
		mw.Monitor.Refresh()
		mw.Mb.VRAM.Dirty = false
	}
	mw.RamView.Refresh()
}