package cpu

type InstrType uint8
type RegisterAddress uint8

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
	Rd, Rs1, Rs2   RegisterAddress
	Imm            int32 
}


func signExtend(u uint32, bits uint) int32 {
	shift := 32 - bits
	return int32(u<<shift) >> shift
}

func (inst *Instruction) decodeRType(raw uint32) {
	inst.Rd     = RegisterAddress((raw >> 7)  & 0x1F)
	inst.Funct3 = uint8((raw >> 12) & 0x07)
	inst.Rs1    = RegisterAddress((raw >> 15) & 0x1F)
	inst.Rs2    = RegisterAddress((raw >> 20) & 0x1F)
	inst.Funct7 = uint8((raw >> 25) & 0x7F)
}

// --------------------------------------
// I-TYPE  imm[11:0]  -> 12b sign-extend
// --------------------------------------
func (inst *Instruction) decodeIType(raw uint32) {
	inst.Rd     = RegisterAddress((raw >> 7)  & 0x1F)
	inst.Funct3 = uint8((raw >> 12) & 0x07)
	inst.Rs1    = RegisterAddress((raw >> 15) & 0x1F)

	imm12 := (raw >> 20) & 0xFFF // bits [31:20]
	inst.Imm = signExtend(imm12, 12)
}

// --------------------------------------
// S-TYPE  imm = { [31:25]=imm[11:5], [11:7]=imm[4:0] } -> 12b sign-extend
// --------------------------------------
func (inst *Instruction) decodeSType(raw uint32) {
	inst.Funct3 = uint8((raw >> 12) & 0x07)
	inst.Rs1    = RegisterAddress((raw >> 15) & 0x1F)
	inst.Rs2    = RegisterAddress((raw >> 20) & 0x1F)

	imm := ((raw >> 25) << 5) | ((raw >> 7) & 0x1F) // 12 bits brutos
	inst.Imm = signExtend(imm, 12)
}


// --------------------------------------
// B-TYPE  imm = { [31]=imm[12], [7]=imm[11], [30:25]=imm[10:5], [11:8]=imm[4:1] } << 1
// total 13b -> sign-extend
// --------------------------------------
func (inst *Instruction) decodeBType(raw uint32) {
	inst.Funct3 = uint8((raw >> 12) & 0x07)
	inst.Rs1    = RegisterAddress((raw >> 15) & 0x1F)
	inst.Rs2    = RegisterAddress((raw >> 20) & 0x1F)

	imm :=
		(((raw >> 31) & 0x1) << 12) | // imm[12]
		(((raw >> 7)  & 0x1) << 11) | // imm[11]
		(((raw >> 25) & 0x3F) << 5) | // imm[10:5]
		(((raw >> 8)  & 0x0F) << 1)   // imm[4:1]
	inst.Imm = signExtend(imm, 13)
}

// --------------------------------------
// U-TYPE  imm = raw[31:12] << 12 
// --------------------------------------
func (inst *Instruction) decodeUType(raw uint32) {
	inst.Rd  = RegisterAddress((raw >> 7) & 0x1F)
	inst.Imm = int32(raw & 0xFFFFF000) // sem sign-extend
}

// --------------------------------------
// J-TYPE  imm = { [31]=imm[20], [19:12]=imm[19:12], [20]=imm[11], [30:21]=imm[10:1] } << 1
// total 21b -> sign-extend
// --------------------------------------
func (inst *Instruction) decodeJType(raw uint32) {
	inst.Rd = RegisterAddress((raw >> 7) & 0x1F)

	imm :=
		(((raw >> 31) & 0x1) << 20) |      // imm[20]
		(((raw >> 12) & 0xFF) << 12) |     // imm[19:12]
		(((raw >> 20) & 0x1) << 11) |      // imm[11]
		(((raw >> 21) & 0x3FF) << 1)       // imm[10:1]
	inst.Imm = signExtend(imm, 21)
}
