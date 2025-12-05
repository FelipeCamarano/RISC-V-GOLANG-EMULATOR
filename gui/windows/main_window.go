package windows

import (
	"fmt"
	"math"
	"path/filepath"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/components"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

const (
	TargetFPS   = 60
	FrameTime   = time.Second / TargetFPS
	MinSpeedIPS = 1.0
	MaxSpeedIPS = 1000000.0
)

type MainWindow struct {
	Window    fyne.Window
	App       fyne.App
	Mb        *motherboard.Motherboard

	PcBinding     binding.String
	StatusBinding binding.String
	RegBindings   []binding.String

	Monitor   *components.MonitorWrapper
	SpeedCtrl *components.SpeedControl
	RamView   *components.RamViewer
	Terminal  *components.Terminal

	IsRunning bool
	StopChan  chan bool
}

func NewMainWindow(a fyne.App, mb *motherboard.Motherboard) *MainWindow {
	w := &MainWindow{
		App:           a,
		Window:        a.NewWindow("RISC-V Workstation"),
		Mb:            mb,
		PcBinding:     binding.NewString(),
		StatusBinding: binding.NewString(),
		RegBindings:   make([]binding.String, 32),
		StopChan:      make(chan bool),
	}

	w.Window.Resize(fyne.NewSize(1280, 850))
	w.StatusBinding.Set("System Ready. No cartridge loaded.")

	for i := 0; i < 32; i++ {
		w.RegBindings[i] = binding.NewString()
	}

	return w
}

func (mw *MainWindow) Build(onOpenTests func(), onLoadCartridge func(string)) {
	regGroup := components.NewRegisterGrid(mw.RegBindings)
	mw.Monitor = components.NewMonitor(mw.Mb)
	mw.SpeedCtrl = components.NewSpeedControl()
	mw.RamView = components.NewRamViewer(mw.Mb)
	mw.Terminal = components.NewTerminal(mw.Mb)

	toolbar := mw.createToolbar(onOpenTests, onLoadCartridge)

	statusBar := widget.NewLabelWithData(mw.StatusBinding)
	statusBar.TextStyle = fyne.TextStyle{Monospace: true}

	leftSplit := container.NewVSplit(mw.Monitor.Container, mw.Terminal.Container)
	leftSplit.SetOffset(0.70)

	rightSplit := container.NewVSplit(regGroup, mw.RamView.Container)
	rightSplit.SetOffset(0.35)

	mainSplit := container.NewHSplit(leftSplit, rightSplit)
	mainSplit.SetOffset(0.65)

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
				mw.SetStatus("Loaded: " + filepath.Base(path))
			}
		}, mw.Window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".bin", ".rom", ".img"}))
		fd.Show()
	})

	btnStep := widget.NewButtonWithIcon("Step", theme.MediaSkipNextIcon(), func() {
		mw.stepCPU()
		mw.UpdateUI(true)
	})

	btnRun := widget.NewButtonWithIcon("Run/Stop", theme.MediaPlayIcon(), nil)
	btnRun.OnTapped = func() {
		if mw.IsRunning {
			mw.IsRunning = false
			btnRun.SetIcon(theme.MediaPlayIcon())
			mw.SetStatus("Paused")
		} else {
			mw.IsRunning = true
			btnRun.SetIcon(theme.MediaStopIcon())
			mw.SetStatus("Running...")
			go mw.emulatorLoop()
		}
	}

	btnTests := widget.NewButtonWithIcon("Debug", theme.ListIcon(), onOpenTests)

	left := container.NewHBox(widget.NewIcon(theme.ComputerIcon()), pcLabel, layout.NewSpacer())
	right := container.NewHBox(btnLoadCart, btnStep, btnRun, btnTests)
	
	return container.NewVBox(
		container.NewBorder(nil, nil, left, right, container.NewPadded(mw.SpeedCtrl.Container)),
		widget.NewSeparator(),
	)
}

func (mw *MainWindow) Show() {
	mw.UpdateUI(true)
	mw.Window.Show()
}

func (mw *MainWindow) stepCPU() {
	pc := uint32(mw.Mb.CPU.PC)
	if mw.Mb.Bus.ReadWord(pc) == 0 {
		if mw.IsRunning {
			mw.IsRunning = false
			mw.SetStatus("Halted (Instruction 0x00000000)")
		}
		return
	}
	mw.Mb.CPU.Step()
}

func (mw *MainWindow) emulatorLoop() {
	drawTicker := time.NewTicker(FrameTime)
	defer drawTicker.Stop()

	cyclesToRun := 0.0

	for mw.IsRunning {
		speedPercent, _ := mw.SpeedCtrl.Value.Get()
		
		var targetIPS float64
		if speedPercent <= 1.0 {
			targetIPS = MinSpeedIPS
		} else {
			ratio := speedPercent / 100.0
			targetIPS = MinSpeedIPS * math.Pow(MaxSpeedIPS/MinSpeedIPS, ratio)
		}

		instructionsPerFrame := targetIPS / float64(TargetFPS)
		
		cyclesToRun += instructionsPerFrame
		
		count := int(cyclesToRun)
		cyclesToRun -= float64(count)

		for i := 0; i < count; i++ {
			if !mw.IsRunning { break }
			mw.stepCPU()
			
			if count > 1000 {
				if i%100 == 0 { mw.Terminal.CheckOutput() }
			} else {
				mw.Terminal.CheckOutput()
			}
		}
		
		<-drawTicker.C
		
		fyne.Do(func() {
			mw.UpdateUI(false)
		})
	}
	fyne.Do(func() { mw.UpdateUI(true) })
}

func (mw *MainWindow) UpdateUI(fullRefresh bool) {
	mw.PcBinding.Set(fmt.Sprintf("PC: 0x%08X", mw.Mb.CPU.PC))
	for i, val := range mw.Mb.CPU.Registers {
		mw.RegBindings[i].Set(fmt.Sprintf("%08X", val))
	}

	if mw.Mb.VRAM.Dirty {
		mw.Monitor.Refresh()
		mw.Mb.VRAM.Dirty = false
	}

	if fullRefresh {
		mw.RamView.Refresh()
	}
}