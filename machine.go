package main

import (
	"fmt"
	"time"
)

const (
	memSize = 4096
	w       = 64
	h       = 32
)

type Machine struct {
	PC      addr
	Mem     [4096]byte
	stack   Stack
	delay   uint8
	sound   uint8
	I       uint16
	V       [16]byte
	Display [w][h]bool
	scale   int
	pad     int
	step    time.Duration
}

func newMachine(step, scale, pad int) *Machine {
	return &Machine{
		PC:    0x200,
		step:  time.Duration(step),
		scale: scale,
		pad:   pad,
	}
}

type addr uint16
type opcode uint16

func toOpcode(b [2]byte) opcode {
	return opcode(b[1]) + (opcode(b[0]) << 8)
}

func (m *Machine) Run() {
	go m.draw()
	for {
		opc := [2]byte{m.Mem[m.PC], m.Mem[m.PC+1]}
		m.PC += 2
		m.DecodeExecute(toOpcode(opc))
		time.Sleep(m.step * time.Millisecond)
	}
}

func (m *Machine) DecodeExecute(o opcode) {
	op := (o & 0xF000) >> 12
	X := (o & (0x0F00)) >> 8
	Y := (o & (0x00F0)) >> 4
	N := (o & (0x000F))
	NN := (o & (0x00FF))
	NNN := (o & (0x0FFF))
	switch op {
	case 0x0:
		switch NNN {
		case 0x0E0:
			Println("CLEAR")
			m.ClearScreen()
		}
	case 0x1:
		Println("JUMP", NNN)
		m.Jump(addr(NNN))
	case 0x6:
		Println("SETREG", X, NN)
		m.SetReg(int(X), byte(NN))
	case 0x7:
		Println("ADDREG", X, NN)
		m.AddReg(int(X), byte(NN))
	case 0xA:
		Println("SETI", NNN)
		m.SetI(uint16(NNN))
	case 0xD:
		Println("DRAW", X, Y, N)
		m.Draw(int(X), int(Y), int(N))
	default:
		fmt.Printf("Unknown instruction: %x", o)
	}
}

func (m *Machine) ClearScreen() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			m.Display[i][j] = false
		}
	}
}

func (m *Machine) Jump(a addr) {
	m.PC = a
}

func (m *Machine) SetReg(reg int, val byte) {
	m.V[reg] = val
}

func (m *Machine) AddReg(reg int, val byte) {
	m.V[reg] += val
}

func (m *Machine) SetI(val uint16) {
	m.I = val
}

func (m *Machine) Draw(xa, ya, rows int) {
	x := int(m.V[xa]) % w
	y := int(m.V[ya]) % h
	m.V[15] = 0
	for i := 0; i < rows; i++ {
		sprite := m.Mem[m.I+uint16(i)]
		shift := 7
		and := byte(255)
		for xo := 0; xo < 8; xo++ {
			val := (sprite & and) >> shift
			shift--
			and = and >> 1
			if val == 1 {
				// Println("SET", x+xo, y+i)
				c := m.Display[x+xo][y+i]
				m.Display[x+xo][y+i] = !c
				if c {
					m.V[15] = 1
				}
			}
		}
	}
}
