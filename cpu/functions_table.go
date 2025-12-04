package cpu

import "fmt"

type ExecFn func(cpu *CPU, inst Instruction)

var FuncTable = map[InstrType]ExecFn{
	R_TYPE: execRType,
	I_TYPE: execIType,
	S_TYPE: execSType,
	B_TYPE: execBType,
	U_TYPE: execUType,
	J_TYPE: execJType,
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
	0x13: execIArith,
	0x03: execILoad,
	0x67: execIJALR,
	0x73: execISystem,
}

var iArithTable = map[uint8]ExecFn{
	0x0: execADDI,
	0x2: execSLTI,
	0x3: execSLTIU,
	0x4: execXORI,
	0x6: execORI,
	0x7: execANDI,
	0x1: execSLLI,
	0x5: execShiftRightImm,
}

var iLoadTable = map[uint8]ExecFn{
	0x0: execLB,
	0x1: execLH,
	0x2: execLW,
	0x4: execLBU,
	0x5: execLHU,
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

func execILoad(cpu *CPU, inst Instruction) {
	if fn, ok := iLoadTable[inst.Funct3]; ok {
		fn(cpu, inst)
	}
}

func execIJALR(cpu *CPU, inst Instruction) {
	execJALR(cpu, inst)
}

func execISystem(cpu *CPU, inst Instruction) {
	// TODO: Implementar chamadas de sistema (ECALL)
	// 0x73 is the opcode for SYSTEM instructions.
	// In these tests, an ECALL with imm=0 indicates the end of the test.

	// Only handle ECALL (Environment Call) which usually has Imm=0
	// (You can add stricter checks for Imm==0 or Imm==1 for EBREAK if needed)

	// 1. Get the result code from Register x10 (a0)
	//    0 = Success
	//    Anything else = The test number that failed
	result := cpu.Registers[10]

	// 2. Get the current test case number from Register x3 (gp)
	//    Useful for debugging if result != 0
	testNum := cpu.Registers[3]

	if result == 0 {
		fmt.Printf("\n[SUCCESS] Test Passed! All logic checks correct.\n")
	} else {
		fmt.Printf("\n[FAILURE] Test Failed at Case #%d (Error Code: %d)\n", testNum, result)
	}

	// 3. Stop the CPU
	cpu.Stopped = true
}

// ======================================
// S-TYPE
// ======================================

var sTypeTable = map[uint8]ExecFn{
	0x0: execSB,
	0x1: execSH,
	0x2: execSW,
}

func execSType(cpu *CPU, inst Instruction) {
	if fn, ok := sTypeTable[inst.Funct3]; ok {
		fn(cpu, inst)
	}
}

// ======================================
// B-TYPE
// ======================================

var bTypeTable = map[uint8]ExecFn{
	0x0: execBEQ,
	0x1: execBNE,
	0x4: execBLT,
	0x5: execBGE,
	0x6: execBLTU,
	0x7: execBGEU,
}

func execBType(cpu *CPU, inst Instruction) {
	if fn, ok := bTypeTable[inst.Funct3]; ok {
		fn(cpu, inst)
	}
}

// ======================================
// U-TYPE
// ======================================

var uTypeTable = map[uint8]ExecFn{
	0x37: execLUI,
	0x17: execAUIPC,
}

func execUType(cpu *CPU, inst Instruction) {
	if fn, ok := uTypeTable[inst.Opcode]; ok {
		fn(cpu, inst)
	}
}

// ======================================
// J-TYPE
// ======================================

func execJType(cpu *CPU, inst Instruction) {
	execJAL(cpu, inst)
}
