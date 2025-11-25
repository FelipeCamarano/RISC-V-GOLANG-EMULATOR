package cpu

type ExecFn func(cpu *CPU, inst Instruction)


var FuncTable = map[InstrType]ExecFn{
	R_TYPE: execRType,
	I_TYPE: execIType,
// 	S_TYPE: execSType,
// 	B_TYPE: execBType,
// 	U_TYPE: execUType,
// 	J_TYPE: execJType,
}

// ======================================
// R-TYPE
// ======================================

var rTypeTable = map[uint8]map[uint8]ExecFn{
	0x0: {
		0x00: execADD,
		0x20: execSUB,
	},
	0x1: {
		0x00: execSLL,
	},
	0x2: {
		0x00: execSLT,
	},
	0x3: {
		0x00: execSLTU,
	},
	0x4: {
		0x00: execXOR,
	},
	0x5: {
		0x00: execSRL,
		0x20: execSRA,
	},
	0x6: {
		0x00: execOR,
	},
	0x7: {
		0x00: execAND,
	},
}

func execRType(cpu *CPU, inst Instruction) {
	if inner, ok := rTypeTable[inst.Funct3]; ok {
		if fn, ok := inner[inst.Funct7]; ok {
			fn(cpu, inst)
		}
	}
}

// ======================================
// I-TYPE
// ======================================

var iTypeTable = map[uint8]ExecFn{
	0x13: execIArith,  // OP-IMM: ADDI, SLTI, ...
	0x03: execILoad,   // LOAD: LB, LH, LW, LBU, LHU
	0x67: execIJALR,   // JALR
	0x73: execISystem, // ECALL/EBREAK
}

var iArithTable = map[uint8]ExecFn{
	0x0: execADDI,          // ADDI
	0x2: execSLTI,          // SLTI
	0x3: execSLTIU,         // SLTIU
	0x4: execXORI,          // XORI
	0x6: execORI,           // ORI
	0x7: execANDI,          // ANDI
	0x1: execSLLI,          // SLLI
	0x5: execShiftRightImm, // SRLI/SRAI (decidido por imm[11:5])
}

func execIType(cpu *CPU, inst Instruction) {
	if fn, ok := iTypeTable[inst.Opcode]; ok {
		fn(cpu, inst)
	}
}

func execIArith(cpu *CPU, inst Instruction) {
	if fn, ok := iArithTable[inst.Funct3]; ok {
		fn(cpu, inst)
	}
}

func execILoad(cpu *CPU, inst Instruction)   {}
func execIJALR(cpu *CPU, inst Instruction)   {}
func execISystem(cpu *CPU, inst Instruction) {}

