package motherboard

import (
	"fmt"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/bus"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/cpu"
	"github.com/DainSlash/RISC-V-GOLANG-EMULATOR/memory"
)

type Motherboard struct {
	CPU  *cpu.CPU
	Bus  *bus.Bus
	
	RAM       *memory.RAM
	VRAM      *memory.VRAM
	IO        *memory.RAM 
	BIOS      *memory.ROM
	Cartridge *memory.ROM
}

func NewMotherboard(biosPath string) (*Motherboard, error) {
	systemBus := bus.NewBus()

	bios, err := memory.NewROMFromFile(biosPath)
	if err != nil {
		return nil, err
	}

	mainRAM := memory.NewRAM(RAM_LIMIT - RAM_START + 1)
	videoRAM := memory.NewVRAM(VRAM_LIMIT - VRAM_START + 1)
	ioDevice := memory.NewRAM(IO_LIMIT - IO_START + 1)
	
	systemBus.MapDevice(RAM_START, mainRAM)
	systemBus.MapDevice(VRAM_START, videoRAM)
	systemBus.MapDevice(IO_START, ioDevice) 
	systemBus.MapDevice(BIOS_START, bios)

	core := cpu.NewCPU(systemBus, BIOS_START)
	core.PC = cpu.RegisterValue(BIOS_START) 

	mb := &Motherboard{
		CPU:  core,
		Bus:  systemBus,
		RAM:  mainRAM,
		VRAM: videoRAM,
		IO:   ioDevice, 
		BIOS: bios,
	}

	return mb, nil
}


func (mb *Motherboard) InsertCartridge(path string) error {
	cart, err := memory.NewROMFromFile(path)
	if err != nil { return err }
	
	mb.Cartridge = cart
	mb.Bus.MapDevice(CARTRIDGE_START, cart)
	fmt.Printf("Cartucho inserido: %s em 0x%X\n", path, CARTRIDGE_START)
	return nil
}

func (mb *Motherboard) Reset() {
	mb.CPU.Reset()
	mb.RAM.Reset()
	mb.VRAM.Reset()
	mb.IO.Reset()
}