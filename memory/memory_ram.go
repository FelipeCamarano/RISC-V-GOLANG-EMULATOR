package memory

type Byte int8

const RAM_SIZE = 0xFFFFFF // 16 MB

type RAM struct {
	memory [RAM_SIZE]Byte
}

func (ram *RAM) ReadByte(address uint32) Byte{
	return Byte(ram.memory[address])
}