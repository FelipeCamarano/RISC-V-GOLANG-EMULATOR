package main

import (
    "fmt"
    "github.com/DainSlash/RISC-V-GOLANG-EMULATOR/cpu"
)

func main() {
    fmt.Println("RISC-V Emulator iniciado!")

    raw := uint32(0b0000000_01011_00010_000_00000_0110011)
    inst := cpu.Decode(raw)

	var cpu cpu.CPU
	cpu.Registers[inst.Rs1] = 10
	cpu.Registers[inst.Rs2] = 14
	cpu.Registers[inst.Rd] = 696969

    fmt.Printf("Opcode:  0x%02X\n", inst.Opcode)
    fmt.Printf("Type:    %v\n", inst.Type)
    fmt.Printf("Rd:      %08b\n", inst.Rd)
    fmt.Printf("Rs1:     %08b\n", inst.Rs1)
    fmt.Printf("Rs2:     %08b\n", inst.Rs2)
    fmt.Printf("Funct3:  %08b\n", inst.Funct3)
    fmt.Printf("Funct7:  %08b\n", inst.Funct7)
    fmt.Printf("Imm:     %032b\n", inst.Imm)

	cpu.Execute(inst)

	fmt.Printf("CPU Rs1:   %d\n", cpu.Registers[inst.Rs1])
	fmt.Printf("CPU Rs2:   %d\n", cpu.Registers[inst.Rs2])
	fmt.Printf("CPU Rd:    %d\n", cpu.Registers[inst.Rd])
}
