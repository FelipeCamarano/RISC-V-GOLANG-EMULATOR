package windows

import (
	"os"
	"fmt"
	
	"path/filepath"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/components"
)

type TestData struct {
	Path string
	Row  *components.TestRow
}

const TestDirectory = "TESTES HEX RISCV"

func NewTestRunnerWindow(a fyne.App, onDebugLoad func(path string)) fyne.Window {
	win := a.NewWindow("Suite de Testes RISC-V")
	win.Resize(fyne.NewSize(600, 500))

	// 1. Carregar lista de arquivos
	testQueue, err := loadTestQueue(win, onDebugLoad)
	if err != nil {
		win.SetContent(widget.NewLabel(err.Error()))
		return win
	}

	// 2. Construir lista visual
	listContainer := container.NewVBox()
	for _, data := range testQueue {
		listContainer.Add(data.Row.Container)
	}

	// 3. Controles
	progressBar := widget.NewProgressBar()
	lblSummary := widget.NewLabel("Pronto para iniciar.")

	btnRunAll := widget.NewButtonWithIcon("Executar Todos", theme.MediaPlayIcon(), func() {
		RunAllTests(testQueue, progressBar, lblSummary)
	})

	// 4. Layout Final
	layout := container.NewBorder(
		container.NewVBox(
			widget.NewLabelWithStyle("Bateria de Testes", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			btnRunAll, progressBar, lblSummary, widget.NewSeparator(),
		),
		nil, nil, nil,
		container.NewVScroll(listContainer),
	)

	win.SetContent(layout)
	return win
}

func loadTestQueue(win fyne.Window, onDebugLoad func(string)) ([]*TestData, error) {
	files, err := os.ReadDir(TestDirectory)
	if err != nil {
		return nil, fmt.Errorf("Erro: Pasta '%s' não encontrada", TestDirectory)
	}

	var queue []*TestData
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".hex" {
			fullPath := filepath.Join(TestDirectory, f.Name())
			
			debugAction := func() {
				onDebugLoad(fullPath)
				win.Close()
			}

			row := components.NewTestRow(f.Name(), win, debugAction)
			row.SetPending()

			queue = append(queue, &TestData{Path: fullPath, Row: row})
		}
	}
	return queue, nil
}