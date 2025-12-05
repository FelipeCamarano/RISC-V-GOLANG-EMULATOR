package cpu

import (

	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/bus"
)

type RegisterValue uint32
type CPU struct {
	Registers    [32]RegisterValue
	PC           RegisterValue
	ResetVector  RegisterValue
	AddressAdder int32
	rdInterface  bus.ReaderWriter
}

func NewCPU(b bus.ReaderWriter, resetVector uint32) *CPU {
	return &CPU{
		Registers:    [32]RegisterValue{},
		PC:           0,
		ResetVector:  RegisterValue(resetVector),
		AddressAdder: 0,
		rdInterface:  b,
	}
}

func (cpu *CPU) Step() {
	// TODO :: implementar busca da instrução armazenada no endereço de PC
	// Executar a instrução da memoria, e ir para o proximo pc.

	cpu.AddressAdder = 4
	var raw uint32 = cpu.Fetch()
	inst := cpu.Decode(raw)
	
	// fmt.Printf("\ninst: %b\n", inst)
	// fmt.Printf("Opcode:  0x%02X\n", inst.Opcode)
	// fmt.Printf("Type:    %v\n", inst.Type)
	// fmt.Printf("Rd:      %08b\n", inst.Rd)
	// fmt.Printf("Rs1:     %08b\n", inst.Rs1)
	// fmt.Printf("Rs2:     %08b\n", inst.Rs2)
	// fmt.Printf("Funct3:  %08b\n", inst.Funct3)
	// fmt.Printf("Funct7:  %08b\n", inst.Funct7)
	// fmt.Printf("Imm:     %032b\n", inst.Imm)

	cpu.Execute(inst)
	cpu.nextPC(cpu.AddressAdder)

	// fmt.Printf("CPU Rs1:   %d\n", cpu.Registers[inst.Rs1])
	// fmt.Printf("CPU Rs2:   %d\n", cpu.Registers[inst.Rs2])
	// fmt.Printf("CPU Rd:    %d\n", cpu.Registers[inst.Rd])
}

func (cpu *CPU) Execute(inst Instruction) {
	handler, existe := FuncTable[inst.Type]
	
	if !existe {
		return
	}

	handler(cpu, inst)
}

func (cpu *CPU) writeReg(rd RegisterAddress, v RegisterValue) {
	if rd != 0 {
		cpu.Registers[rd] = v
	}
}

func (cpu *CPU) Fetch() uint32 {
	pc := cpu.PC
	rdInterface := cpu.rdInterface

	result := uint32(rdInterface.ReadByte(uint32(pc)))
	result |= uint32(rdInterface.ReadByte(uint32(pc)+1)) << 8
	result |= uint32(rdInterface.ReadByte(uint32(pc)+2)) << 16
	result |= uint32(rdInterface.ReadByte(uint32(pc)+3)) << 24
	// fmt.Printf("\n %032b\n", result)
	return result
}

func (cpu *CPU) nextPC(adder int32) {
	cpu.PC = RegisterValue(int32(cpu.PC) + adder)
}

func (cpu *CPU) SetPC(addr uint32) {
	cpu.PC = RegisterValue(int32(addr))
}

func (cpu *CPU) Reset() {
	cpu.PC = cpu.ResetVector
	for i := range cpu.Registers {
		cpu.Registers[i] = 0
	}	
}