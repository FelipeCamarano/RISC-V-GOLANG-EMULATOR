package cpu

type CPU struct {
	Registers [32]int32
}

func (cpu *CPU) Execute(inst Instruction) {
	handler := FuncTable[inst.Type]
	handler(cpu, inst)
}