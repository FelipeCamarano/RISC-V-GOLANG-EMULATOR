package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

func LoadHexFile(filename string, mb *motherboard.Motherboard) (uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentAddr := uint32(motherboard.RAM_START)
	startAddr := currentAddr
	firstInstructionFound := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "<") {
			continue
		}

		if strings.HasPrefix(line, "@") {
			addrStr := strings.TrimSuffix(line[1:], ":")
			addrVal, err := strconv.ParseUint(addrStr, 16, 32)
			if err != nil {
				return 0, fmt.Errorf("marcador de endereço inválido: %s", line)
			}
			currentAddr = uint32(addrVal)
			continue
		}

		parts := strings.Fields(line)
		for _, part := range parts {
			val, err := strconv.ParseUint(part, 16, 32)
			if err == nil {
				if !firstInstructionFound {
					startAddr = currentAddr
					firstInstructionFound = true
				}
				mb.Bus.WriteWord(currentAddr, uint32(val))
				currentAddr += 4
			}
		}
	}
	return startAddr, nil
}