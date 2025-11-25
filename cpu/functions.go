package cpu


// ======================================
// R-TYPE
// ======================================

func execADD(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(a+b)

	cpu.writeReg(inst.Rd, result)
}

func execSUB(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(a-b)

	cpu.writeReg(inst.Rd, result)
}

func execSLL(cpu *CPU, inst Instruction) {
	shamt 			:= uint32(cpu.Registers[inst.Rs2]) & 0x1F // 0b11111
	value 			:= uint32(cpu.Registers[inst.Rs1])
	result 			:= RegisterValue(value << shamt)
	
	cpu.writeReg(inst.Rd, result)
}

func execSLT(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(0)
	
	if (a<b) {
		result 		= RegisterValue(1)
	}

	cpu.writeReg(inst.Rd, result)
}

func execSLTU(cpu *CPU, inst Instruction) {
	a 				:= uint32(cpu.Registers[inst.Rs1])
	b 				:= uint32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(0)
	
	if a < b {
		result 		= RegisterValue(1)
	}

	cpu.writeReg(inst.Rd, result)
}

func execXOR(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(a^b)

	cpu.writeReg(inst.Rd, result)
}

func execSRL(cpu *CPU, inst Instruction) {
	shamt 			:= uint32(cpu.Registers[inst.Rs2]) & 0x1F
	value 			:= uint32(cpu.Registers[inst.Rs1])
	result 			:= RegisterValue(value >> shamt)

	cpu.writeReg(inst.Rd, result)
}

func execSRA(cpu *CPU, inst Instruction) {
	shamt 			:= uint32(cpu.Registers[inst.Rs2]) & 0x1F
	value 			:= int32(cpu.Registers[inst.Rs1])
	result 			:= RegisterValue(value >> shamt)

	cpu.writeReg(inst.Rd, result)
}

func execOR(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(a|b)

	cpu.writeReg(inst.Rd, result)
}

func execAND(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	b 				:= int32(cpu.Registers[inst.Rs2])
	result 			:= RegisterValue(a&b)

	cpu.writeReg(inst.Rd, result)
}


// ======================================
// I-TYPE
// ======================================


// ADDI  rd = rs1 + imm
func execADDI(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	imm 			:= int32(inst.Imm)
	result 			:= RegisterValue(a+imm)

	cpu.writeReg(inst.Rd, result)
}

func execSLTI(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	imm 			:= int32(inst.Imm)
	result 			:= RegisterValue(0)

	if a < imm {
		result 		= RegisterValue(1)
	}

	cpu.writeReg(inst.Rd, result)
}

func execSLTIU(cpu *CPU, inst Instruction) {
	a 				:= uint32(cpu.Registers[inst.Rs1])
	imm 			:= uint32(inst.Imm)
	result 			:= RegisterValue(0)

	if a < imm {
		result 		= RegisterValue(1)
	}

	cpu.writeReg(inst.Rd, result)
}


func execXORI(cpu *CPU, inst Instruction) {
	a 				:= int32(cpu.Registers[inst.Rs1])
	imm 			:= int32(inst.Imm)
	result 			:= RegisterValue(a^imm)

	cpu.writeReg(inst.Rd, result)
}

// ORI   rd = rs1 | imm
func execORI(cpu *CPU, inst Instruction) {
	a 				:= RegisterValue(cpu.Registers[inst.Rs1])
	imm 			:= RegisterValue(inst.Imm)

	cpu.writeReg(inst.Rd, a|imm)
}

// ANDI  rd = rs1 & imm
func execANDI(cpu *CPU, inst Instruction) {
	a 				:= RegisterValue(cpu.Registers[inst.Rs1])
	imm 			:= RegisterValue(inst.Imm)

	cpu.writeReg(inst.Rd, a&imm)
}


// SLLI  rd = rs1 << shamt
func execSLLI(cpu *CPU, inst Instruction) {
	shamt 			:= uint32(inst.Imm) & 0x1F
	val 			:= uint32(cpu.Registers[inst.Rs1])
	cpu.writeReg(inst.Rd, RegisterValue(val<<shamt))
}

// SRLI/SRAI
func execShiftRightImm(cpu *CPU, inst Instruction) {
	imm   			:= uint32(inst.Imm)
	shamt 			:= imm & 0x1F
	top7  			:= (imm >> 5) & 0x7F

	switch top7 {
	case 0x00:
		val := uint32(cpu.Registers[inst.Rs1])
		cpu.writeReg(inst.Rd, RegisterValue(val>>shamt))
	case 0x20:
		val := int32(cpu.Registers[inst.Rs1])
		cpu.writeReg(inst.Rd, RegisterValue(val>>shamt))
	default:
		// valores reservados/ilegais: no-op por enquanto
	}
}
