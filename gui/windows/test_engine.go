package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui/utils"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

func RunAllTests(queue []*TestData, progress *widget.ProgressBar, summary *widget.Label) {
	progress.SetValue(0)
	summary.SetText("Executando...")

	go func() {
		passed := 0
		total := len(queue)

		for i, data := range queue {
			fyne.Do(func() {
				data.Row.SetRunning()
				summary.SetText(fmt.Sprintf("Rodando: %s", data.Row.Name))
			})

			cycles, result, testCase, err := executeSingleTest(data.Path)

			fyne.Do(func() {
				if err != nil {
					data.Row.SetFail(err.Error(), cycles, result, testCase)
				} else if result == 0 {
					data.Row.SetPass(cycles, result)
					passed++
				} else {
					data.Row.SetFail("Falha Lógica", cycles, result, testCase)
				}
				progress.SetValue(float64(i+1) / float64(total))
			})
		}

		fyne.Do(func() {
			summary.SetText(fmt.Sprintf("Fim! Sucesso: %d | Falhas: %d", passed, total-passed))
		})
	}()
}

func executeSingleTest(path string) (cycles int, result int32, testCase int32, err error) {
	mb, err := motherboard.NewMotherboard("bios.bin")
	if err != nil {
		fmt.Printf("Aviso no TestEngine: Não foi possível carregar bios.bin: %v\n", err)
	}

	startPC, loadErr := utils.LoadHexFile(path, mb)
	if loadErr != nil {
		return 0, 0, 0, loadErr
	}

	mb.CPU.SetPC(startPC)
	maxCycles := 200000
	cycles = 0

	for cycles < maxCycles {
		if mb.Bus.ReadWord(uint32(mb.CPU.PC)) == 0 {
			break
		}
		mb.CPU.Step()
		cycles++
	}

	if cycles >= maxCycles {
		return cycles, 0, 0, fmt.Errorf("Timeout")
	}

	return cycles, int32(mb.CPU.Registers[10]), int32(mb.CPU.Registers[3]), nil
}