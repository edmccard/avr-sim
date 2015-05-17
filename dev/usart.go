package dev

import (
    "github.com/edmccard/avr-sim/core"
)

type USART struct {
    ucsraU2X  byte
    ucsrbRXEN byte
    ucsrbTXEN byte
    ubbrh     byte
    ucsrc     byte
    ucsrcRead int64
    timer     *core.Timer
    canRead   bool
    read      chan byte
    canWrite  bool
    write     chan byte
}

func NewUSART(read, write chan byte, timer *core.Timer) *USART {
    return &USART{timer: timer, read: read, write: write}
}

func (usart *USART) ReadUCSRA(addr core.Addr) (val byte) {
    usart.canRead = false
    if len(usart.read) > 0 {
        usart.canRead = true
        val |= 0x80
    }

    // TXC is always set when TXEN = 1
    if usart.ucsrbTXEN != 0 {
        val |= 0x40
    }

    usart.canWrite = false
    if len(usart.write) < cap(usart.write) {
        usart.canWrite = true
        val |= 0x20
    }

    val |= usart.ucsraU2X

    return val
}

func (usart *USART) WriteUCSRA(addr core.Addr, val byte) {
    // do not accept MPCM = 1 or TXC = 1
    if (val & 0x41) != 0 {
        panic("unsupported UCSRA configuration")
    }
    // store in case it is read back
    usart.ucsraU2X = val & 0x02
}

func (usart *USART) ReadUCSRB(addr core.Addr) byte {
    return usart.ucsrbRXEN | usart.ucsrbTXEN
}

func (usart *USART) WriteUCSRB(addr core.Addr, val byte) {
    // ignore read-only bit 1 RXB8
    val &^= 2
    // accept ones only in RXEN/TXEN
    if (val & 0xe7) != 0 {
        panic("unsupported UCSRB configuration")
    }
    usart.ucsrbRXEN = val & 0x10
    usart.ucsrbTXEN = val & 0x08
}

func (usart *USART) ReadUCSRC(addr core.Addr) byte {
    cyc := usart.timer.GetCount()
    if cyc == (usart.ucsrcRead + 1) {
        return usart.ucsrc
    } else {
        usart.ucsrcRead = cyc
        return usart.ubbrh
    }
}

func (usart *USART) WriteUCSRC(addr core.Addr, val byte) {
    if (val & 0x80) != 0 {
        if (val & 0x06) == 0 {
            panic("unsupported UCSRC configuration")
        }
        usart.ucsrc = val
    } else {
        usart.ubbrh = val
    }
}

func (usart *USART) ReadUDR(addr core.Addr) byte {
    if usart.ucsrbRXEN == 0 || !usart.canRead {
        return 0
    }
    var c byte
    select {
    case c  = <- usart.read:
        return c
    default:
        panic("read sync error")
    }
}

func (usart *USART) WriteUDR(addr core.Addr, val byte) {
    if usart.ucsrbTXEN == 0 || !usart.canWrite {
        return
    }
    select {
    case usart.write <- val:
    default:
        panic("write sync error")
    }
}
