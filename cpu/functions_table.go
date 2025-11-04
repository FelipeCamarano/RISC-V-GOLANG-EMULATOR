package cpu

type ExecFn func(cpu *CPU, inst Instruction)


var FuncTable = map[InstrType]ExecFn{
    R_TYPE: execRType, // ADD, SUB, AND...
    // I_TYPE: execIType, 
}


var rTypeTable = map[uint8]map[uint8]ExecFn{
    0x0: {
        0x00: execADD,
    },
}


func execRType(cpu *CPU, inst Instruction) {
	if inner, ok := rTypeTable[inst.Funct3]; ok {
        if fn, ok := inner[inst.Funct7]; ok {
            fn(cpu, inst)
        }
    }
}
