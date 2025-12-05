package motherboard

const (
    // (RAM)
    RAM_START       = 0x00000000
    RAM_LIMIT       = 0x0007FFFF 

    // VRAM
    VRAM_START      = 0x00080000
    VRAM_LIMIT      = 0x0008FFFF

    // CARTUCHOS
    CARTRIDGE_START = 0x00090000
    
    // IO
    IO_START        = 0x0009FC00
    IO_LIMIT        = 0x0009FFFF
    
    BIOS_START      = 0x80000000 
)
