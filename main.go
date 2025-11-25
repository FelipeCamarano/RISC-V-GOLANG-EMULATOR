package main

import (
    "fmt"
    "github.com/DainSlash/RISC-V-GOLANG-EMULATOR/cpu"
)

func main() {
    fmt.Println("RISC-V Emulator iniciado!")
    
    var cpu cpu.CPU

    cpu.Registers[1] = 10
    cpu.Registers[2] = 3
    cpu.Registers[3] = 6

    // programa := [...]uint32{
    //     0b0000000_00010_00001_000_00011_011_0011, // 10+3
    //     0b0000000_00001_00011_000_00011_011_0011, // 13 + 10
    //     0b0000000_00001_00011_000_00011_011_0011, // 23 + 10
    //     0b0000000_00001_00011_000_00011_011_0011, // 33 + 10
    //     0b0000000_00001_00011_000_00011_011_0011, // 43 + 10
    //     0b0000000_00001_00011_000_00011_011_0011, // 53 + 10
    // }

    programa := [...]uint32{
        0b0000000_00010_00001_001_00011_011_0011, // 10+3
    }

    
    
    for _, element := range programa {
        cpu.Step(element)
    }
}
