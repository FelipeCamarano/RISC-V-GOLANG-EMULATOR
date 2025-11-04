package cpu

type InstrType uint8

const (
	TypeInvalid InstrType = iota
	R_TYPE
	I_TYPE
	S_TYPE
	B_TYPE
	J_TYPE
	U_TYPE
)

type Instruction struct {
	Opcode         uint8
	Type           InstrType
	Funct3, Funct7 uint8
	Rd, Rs1, Rs2   uint8
	Imm            uint32
}


func (inst *Instruction) decodeRType(raw uint32) {
	inst.Rd 			= uint8((raw >> 7)  & 0b0001_1111)
	inst.Funct3 		= uint8((raw >> 12) & 0b0000_0111)
	inst.Rs1 			= uint8((raw >> 15) & 0b0001_1111)
	inst.Rs2 			= uint8((raw >> 20) & 0b0001_1111)
	inst.Funct7 		= uint8((raw >> 25) & 0b0111_1111)
}

func (inst *Instruction) decodeIType(raw uint32) {
	inst.Rd     		= uint8((raw >> 7)  & 0b0001_1111)
	inst.Funct3 		= uint8((raw >> 12) & 0b0000_0111)
	inst.Rs1    		= uint8((raw >> 15) & 0b0001_1111)
	inst.Imm    		= uint32(int32(raw) >> 20)
}

func (inst *Instruction) decodeSType(raw uint32) {
	imm4_0  			:= (raw >> 7)  & 0b0001_1111
	imm11_5 			:= (raw >> 25) & 0b0111_1111

	inst.Funct3 		= uint8((raw >> 12) & 0b0000_0111)
	inst.Rs1    		= uint8((raw >> 15) & 0b0001_1111)
	inst.Rs2    		= uint8((raw >> 20) & 0b0001_1111)
	inst.Imm    		= uint32(int32((imm11_5<<5)|imm4_0) << 20 >> 20)
}

func (inst *Instruction) decodeBType(raw uint32) {
	imm11   			:= (raw >> 7)  & 0b0000_0001
	imm4_1  			:= (raw >> 8)  & 0b0000_1111
	imm10_5 			:= (raw >> 25) & 0b0011_1111
	imm12   			:= (raw >> 31) & 0b0000_0001

	inst.Funct3 		= uint8((raw >> 12) & 0b0000_0111)
	inst.Rs1    		= uint8((raw >> 15) & 0b0001_1111)
	inst.Rs2    		= uint8((raw >> 20) & 0b0001_1111)

	imm 				:= (imm12 << 12) | (imm11 << 11) | (imm10_5 << 5) | (imm4_1 << 1)
	inst.Imm 			= uint32(int32(imm << 19) >> 19)
}

func (inst *Instruction) decodeUType(raw uint32) {
	inst.Rd  			= uint8((raw >> 7) & 0b0001_1111)
	inst.Imm 			= raw & 0b1111_1111_1111_1111_1111_0000_0000_0000
}

func (inst *Instruction) decodeJType(raw uint32) {
	inst.Rd 			= uint8((raw >> 7) & 0b0001_1111)

	imm20    			:= (raw >> 31) & 0b0000_0001
	imm10_1  			:= (raw >> 21) & 0b0011_1111_1111
	imm11    			:= (raw >> 20) & 0b0000_0001
	imm19_12 			:= (raw >> 12) & 0b0000_1111_1111

	imm 				:= (imm20 << 20) | (imm19_12 << 12) | (imm11 << 11) | (imm10_1 << 1)
	inst.Imm 			= uint32(int32(imm << 11) >> 11)
}
