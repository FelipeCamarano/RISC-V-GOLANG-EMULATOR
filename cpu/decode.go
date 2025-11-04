package cpu

const opcodeMask = 0x7F // 0b0111_1111
var opcodeTypeMap = map[uint8]InstrType{
	0x33: R_TYPE, // 0b 0011_0011
	0x13: I_TYPE, // 0b 0001_0011
	0x03: I_TYPE, // 0b 0000_0011
	0x67: I_TYPE, // 0b 0110_0111
	0x73: I_TYPE, // 0b 0111_0011
	0x23: S_TYPE, // 0b 0010_0011
	0x63: B_TYPE, // 0b 0110_0011
	0x6F: J_TYPE, // 0b 0110_1111
	0x37: U_TYPE, // 0b 0011_0111
	0x17: U_TYPE, // 0b 0001_0111
}

func Decode(raw uint32) Instruction {
	inst := Instruction{}

	inst.Opcode = uint8(raw & opcodeMask)
	inst.Type = opcodeTypeMap[inst.Opcode]

	switch inst.Type {
	case R_TYPE:
		inst.decodeRType(raw)
	case I_TYPE:
		inst.decodeIType(raw)
	case S_TYPE:
		inst.decodeSType(raw)
	case B_TYPE:
		inst.decodeBType(raw)
	case U_TYPE:
		inst.decodeUType(raw)
	case J_TYPE:
		inst.decodeJType(raw)
	default:
		inst.Type = TypeInvalid
	}

	return inst
}