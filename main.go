package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/memory"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

func main() {

	mainboard := motherboard.NewMotherboard(motherboard.DefaultRAMSize, memory.BootProgram())
	mainboard.IntialBOOT()
	mainboard.CPU.Step()

	/*
	   dir := "TESTES HEX RISCV"

	   // Lê todos os arquivos da pasta
	   files, err := os.ReadDir(dir)

	   	if err != nil {
	   		fmt.Printf("Erro ao ler diretório: %v\n", err)
	   		return
	   	}

	   fmt.Println("=== INICIANDO BATERIA DE TESTES ===")

	   	for _, f := range files {
	   		// Filtra apenas arquivos .hex
	   		if filepath.Ext(f.Name()) == ".hex" {
	   			fullPath := filepath.Join(dir, f.Name())
	   			runSingleTest(fullPath, f.Name())
	   		}
	   	}
	*/
}

//Testes unitários inspirados no repositório

func runSingleTest(path string, name string) {
	// 1. Cria uma nova Motherboard limpa para cada teste
	mb := motherboard.NewMotherboard(0, nil)

	// 2. Carrega o arquivo
	entryPoint, err := LoadHexFile(path, mb)
	if err != nil {
		fmt.Printf("[%s] Erro ao carregar: %v\n", name, err)
		return
	}

	// 3. Configura PC inicial
	mb.CPU.SetPC(entryPoint)

	// 4. Loop de execução
	// Adicionei um limite de ciclos para evitar loops infinitos caso algo quebre
	maxCycles := 100000
	cycles := 0

	for !mb.CPU.Stopped && cycles < maxCycles {
		mb.CPU.Step()
		cycles++
	}

	// 5. Verifica o resultado
	// O registrador x10 (índice 10) contém o resultado do teste
	// O registrador x3 (índice 3) contém o número do caso de teste (útil em falhas)
	result := mb.CPU.Registers[10]
	testCase := mb.CPU.Registers[3]

	if cycles >= maxCycles {
		fmt.Printf("FAIL \t %-20s (Time Limit Exceeded)\n", name)
	} else if result == 0 {
		fmt.Printf("PASS \t %-20s\n", name)
	} else {
		fmt.Printf("FAIL \t %-20s (Error Code: %d, Test Case: %d)\n", name, result, testCase)
	}
}

// Mantive sua função LoadHexFile, apenas ajustada para não depender de variáveis globais
func LoadHexFile(filename string, mb *motherboard.Motherboard) (uint32, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentAddr := uint32(motherboard.DefaultRAMBase)
	startAddr := currentAddr
	firstInstructionFound := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 1. Ignora linhas vazias
		if line == "" {
			continue
		}

		// 2. NOVA VERIFICAÇÃO: Ignora rótulos de seção (ex: <tdat2>:)
		if strings.HasPrefix(line, "<") {
			continue
		}

		// 3. Trata marcadores de endereço (@80000000)
		if strings.HasPrefix(line, "@") {
			addrStr := strings.TrimSuffix(line[1:], ":")
			addrVal, err := strconv.ParseUint(addrStr, 16, 32)
			if err != nil {
				return 0, fmt.Errorf("invalid address marker: %s", line)
			}
			currentAddr = uint32(addrVal)
			continue
		}

		// 4. Trata instruções ou dados hexadecimais
		val, err := strconv.ParseUint(line, 16, 32)
		if err != nil {
			return 0, fmt.Errorf("invalid hex data: %s", line)
		}

		if !firstInstructionFound {
			startAddr = currentAddr
			firstInstructionFound = true
		}

		mb.Bus.WriteWord(currentAddr, uint32(val))
		currentAddr += 4
	}
	return startAddr, nil
}
