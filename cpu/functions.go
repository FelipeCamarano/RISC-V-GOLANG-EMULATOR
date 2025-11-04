package cpu

func execADD(cpu *CPU, inst Instruction) {
    cpu.Registers[inst.Rd] = cpu.Registers[inst.Rs1] + cpu.Registers[inst.Rs2]
}
