package memory

type RAM struct {
	data []Byte
}

func NewRAM(size uint32) *RAM {
	return &RAM{
		data: make([]Byte, size),
	}
}

func (r *RAM) Size() uint32 {
	return uint32(len(r.data))
}

func (r *RAM) ReadByte(offset uint32) Byte {
	if offset >= uint32(len(r.data)) {
		return 0
	}
	return r.data[offset]
}

func (r *RAM) WriteByte(offset uint32, data Byte) {
	if offset < uint32(len(r.data)) {
		r.data[offset] = data
	}
}

func (r *RAM) Reset() {
	for i := range r.data {
		r.data[i] = 0
	}
}