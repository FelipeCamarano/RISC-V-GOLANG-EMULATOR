package memory

import "fmt"

type VRAM struct {
	data []Byte
	Dirty bool
}

func NewVRAM(size uint32) *VRAM {
	return &VRAM{
		data: make([]Byte, size),
		Dirty: false,
	}
}

func (v *VRAM) Size() uint32 {
	return uint32(len(v.data))
}

func (v *VRAM) ReadByte(offset uint32) Byte {
	if offset >= uint32(len(v.data)) {
		return 0 
	}
	return v.data[offset]
}

func (v *VRAM) WriteByte(offset uint32, data Byte) {
	if offset < uint32(len(v.data)) {
		v.data[offset] = data
		v.Dirty = true
	}
}


func (v *VRAM) DumpToTerminal() {
    for _, b := range v.data {
        if b >= 32 && b <= 126 {
            fmt.Printf("%c", b)
        } else if b == 0 {
			
        } else {
            fmt.Print(".")
        }
    }
}

func (v *VRAM) Reset() {
	for i := range v.data {
		v.data[i] = 0
	}
	v.Dirty = false
}