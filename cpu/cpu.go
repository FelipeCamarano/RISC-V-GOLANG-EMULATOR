package cpu

import (
	"fmt"
    "github.com/DainSlash/RISC-V-GOLANG-EMULATOR/memory"
)

type RegisterValue int32;
type CPU struct {
	Registers [32]RegisterValue
	PC RegisterValue
	AddressAdder int32
	ram memory.RAM // Vai ser substituido por uma interface de barramento dedicado que acessa memoria
}

func (cpu *CPU) Step() {
	// TODO :: implementar busca da instrução armazenada no endereço de PC
	// Executar a instrução da memoria, e ir para o proximo pc.

	cpu.AddressAdder = 4
	var raw uint32 = cpu.Fetch(cpu.ram, cpu.PC)
    inst := cpu.Decode(raw)
	
	fmt.Printf("Opcode:  0x%02X\n", inst.Opcode)
	fmt.Printf("Type:    %v\n", inst.Type)
	fmt.Printf("Rd:      %08b\n", inst.Rd)
	fmt.Printf("Rs1:     %08b\n", inst.Rs1)
	fmt.Printf("Rs2:     %08b\n", inst.Rs2)
	fmt.Printf("Funct3:  %08b\n", inst.Funct3)
	fmt.Printf("Funct7:  %08b\n", inst.Funct7)
	fmt.Printf("Imm:     %032b\n", inst.Imm)
    
	
	cpu.Execute(inst)
	cpu.nextPC(cpu.AddressAdder)

    fmt.Printf("CPU Rs1:   %d\n", cpu.Registers[inst.Rs1])
    fmt.Printf("CPU Rs2:   %d\n", cpu.Registers[inst.Rs2])
    fmt.Printf("CPU Rd:    %d\n", cpu.Registers[inst.Rd])
}

func (cpu *CPU) Execute(inst Instruction) {
	handler := FuncTable[inst.Type]
	handler(cpu, inst)
}

func (cpu *CPU) writeReg(rd RegisterAddress, v RegisterValue) {
	if rd != 0 {
		cpu.Registers[rd] = v
	}
}

func (cpu *CPU) Fetch(ram memory.RAM, pc RegisterValue) uint32 {
	result := uint32(ram.ReadByte(uint32(pc)))
	result |= uint32(ram.ReadByte(uint32(pc)+1)) << 8
	result |= uint32(ram.ReadByte(uint32(pc)+2)) << 16
	result |= uint32(ram.ReadByte(uint32(pc)+3)) << 24
	return result
}

func (cpu *CPU) nextPC(adder int32) {
	cpu.PC = RegisterValue(int32(cpu.PC) + adder)
}