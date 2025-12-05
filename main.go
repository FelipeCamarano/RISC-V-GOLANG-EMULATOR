package main

import (
	"fmt"
	"log"

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/gui"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/motherboard"
)

func main() {
	fmt.Println("Iniciando RISC-V Emulator...")

	// 1. Setup do Hardware
	mb, err := motherboard.NewMotherboard("bios.bin")
	if err != nil {
		log.Fatalf("Erro crítico de Boot: %v", err)
	}

	// Injeta valores para teste
	mb.CPU.Registers[1] = 10
	mb.CPU.Registers[2] = 3


	// 2. Interface
	app := gui.New(mb)
	app.Run()
}