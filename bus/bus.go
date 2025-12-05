package bus

import (
    "github.com/DainSlash/RISC-V-GOLANG-EMULATOR/memory"
)

type ReaderWriter interface {
    ReadByte(addr uint32) memory.Byte
    ReadHalf(addr uint32) uint16
    ReadWord(addr uint32) uint32
    
    WriteByte(addr uint32, data memory.Byte)
    WriteHalf(addr uint32, val uint16)
    WriteWord(addr uint32, val uint32)
}

type deviceMapping struct {
    start uint32
    end   uint32
    dev   memory.Device
}

type Bus struct {
    devices []deviceMapping
}

func NewBus() *Bus {
    return &Bus{
        devices: make([]deviceMapping, 0),
    }
}

func (b *Bus) MapDevice(start uint32, dev memory.Device) {
    size := dev.Size()
    if size == 0 {
        return
    }
    end := start + size - 1

    for i := range b.devices {
        if b.devices[i].start == start {
            b.devices[i].dev = dev
            b.devices[i].end = end
            return
        }
    }

    m := deviceMapping{
        start: start,
        end:   end,
        dev:   dev,
    }
    b.devices = append(b.devices, m)
}

func (b *Bus) findDevice(addr uint32) *deviceMapping {
    for i := range b.devices {
        deviceMapping := &b.devices[i]
        if addr >= deviceMapping.start && addr <= deviceMapping.end {
            return deviceMapping
        }
    }
    return nil
}

// ==========================================
//  READ / WRITE
// ==========================================

func (b *Bus) ReadByte(addr uint32) memory.Byte {
    deviceMapping := b.findDevice(addr)
    if deviceMapping == nil {
        return 0
    }

    offset := addr - deviceMapping.start
    return deviceMapping.dev.ReadByte(offset)
}

func (b *Bus) ReadHalf(addr uint32) uint16 {
    low := uint16(b.ReadByte(addr))
    high := uint16(b.ReadByte(addr + 1))

    return (high << 8) | low
}

func (b *Bus) ReadWord(addr uint32) uint32 {
    low := uint32(b.ReadHalf(addr))
    high := uint32(b.ReadHalf(addr + 2))
    
    return (high << 16) | low
}



func (b *Bus) WriteByte(addr uint32, data memory.Byte) {
    deviceMapping := b.findDevice(addr)
    if deviceMapping == nil {
        return
    }

    offset := addr - deviceMapping.start
    deviceMapping.dev.WriteByte(offset, data)
}

func (b *Bus) WriteHalf(addr uint32, val uint16) {
    b.WriteByte(addr, memory.Byte(val&0xFF))
    b.WriteByte(addr+1, memory.Byte((val>>8)&0xFF))
}

func (b *Bus) WriteWord(addr uint32, val uint32) {
    b.WriteHalf(addr, uint16(val&0xFFFF))
    b.WriteHalf(addr+2, uint16((val>>16)&0xFFFF))
}