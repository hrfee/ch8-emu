package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	memSize = 4096
	w       = 64
	h       = 32
)

type Machine struct {
	PC                 addr
	Mem                [4096]byte
	stack              Stack
	delay              uint8
	sound              uint8
	I                  uint16
	V                  [16]byte
	Display            [w][h]bool
	Canvas             chan [2]int
	keys               Keys
	scale              int
	pad                int
	step               time.Duration
	fontMap            map[byte]addr
	IncrementLoadStore bool
}

func newMachine(step time.Duration, scale, pad int) *Machine {
	m := &Machine{
		PC:    0x200,
		step:  step,
		scale: scale,
		pad:   pad,
		keys:  newKeys(),
	}
	return m
}

type addr uint16
type opcode uint16

func toOpcode(b [2]byte) opcode {
	return opcode(b[1]) + (opcode(b[0]) << 8)
}

func (m *Machine) Load(f *os.File) {
	off := int64(0)
	ops := make([]byte, 1)
	for {
		_, err := f.ReadAt(ops, off)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		m.Mem[off+0x200] = ops[0]
		off++
	}
}

func (m *Machine) Run() {
	go m.draw()
	go func() {
		for range time.Tick((1000000000 / 60) * time.Nanosecond) {
			if m.delay != 0 {
				m.delay--
			}
			if m.sound != 0 {
				m.sound--
			}
		}
	}()
	for {
		opc := [2]byte{m.Mem[m.PC], m.Mem[m.PC+1]}
		m.PC += 2
		m.DecodeExecute(toOpcode(opc))
		time.Sleep(m.step)
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
		case 0x0EE:
			Println("RETSUB")
			m.RetSub()
		default:
			fmt.Printf("Unknown instruction: %x", o)
		}
	case 0x1:
		Println("JUMP", NNN)
		m.Jump(addr(NNN))
	case 0x2:
		Println("JUMPSUB", NNN)
		m.JumpSub(addr(NNN))
	case 0x3:
		Println("SKIPIFREG", X, NN)
		m.SkipIfReg(int(X), byte(NN))
	case 0x4:
		Println("SKIPIFNOTREG", X, NN)
		m.SkipIfNotReg(int(X), byte(NN))
	case 0x5:
		Println("SKIPIFREGS", X, Y)
		m.SkipIfRegs(int(X), int(Y))
	case 0x6:
		Println("SETREG", X, NN)
		m.SetReg(int(X), byte(NN))
	case 0x7:
		Println("ADDREG", X, NN)
		m.AddReg(int(X), byte(NN))
	case 0x8:
		switch N {
		case 0x0:
			Println("SETREGS", X, Y)
			m.SetRegs(int(X), int(Y))
		case 0x1:
			Println("ORREGS", X, Y)
			m.OrRegs(int(X), int(Y))
		case 0x2:
			Println("ANDREGS", X, Y)
			m.AndRegs(int(X), int(Y))
		case 0x3:
			Println("XORREGS", X, Y)
			m.XorRegs(int(X), int(Y))
		case 0x4:
			Println("ADDREGS", X, Y)
			m.Add(int(X), int(Y))
		case 0x5:
			Println("SUBFROMX", X, Y)
			m.SubFromX(int(X), int(Y))
		case 0x6:
			Println("SHIFTRIGHT", X, Y)
			m.ShiftRight(int(X), int(Y))
		case 0x7:
			Println("SUBFROMY", X, Y)
			m.SubFromY(int(X), int(Y))
		case 0xE:
			Println("SHIFTLEFT", X, Y)
			m.ShiftLeft(int(X), int(Y))
		default:
			fmt.Printf("Unknown instruction: %x", o)
		}
	case 0x9:
		Println("SKIPIFNOTREGS", X, Y)
		m.SkipIfNotRegs(int(X), int(Y))
	case 0xA:
		Println("SETI", NNN)
		m.SetI(uint16(NNN))
	case 0xB:
		Println("JUMPOFFSET", NNN)
		m.JumpOffset(addr(NNN))
	case 0xC:
		Println("RANDAND", X, NN)
		m.RandAnd(int(X), byte(NN))
	case 0xD:
		Println("DRAW", X, Y, N)
		m.Draw(int(X), int(Y), int(N))
	case 0xE:
		switch NN {
		case 0x9E:
			Println("SKIPIFPRESSED", X)
			m.SkipIfPressed(int(X))
		case 0xA1:
			Println("SKIPIFNOTPRESSED", X)
			m.SkipIfNotPressed(int(X))
		default:
			fmt.Printf("Unknown instruction: %x", o)
		}
	case 0xF:
		switch NN {
		case 0x07:
			Println("READDELAY", X)
			m.ReadDelay(int(X))
		case 0x15:
			Println("SETDELAY", X)
			m.SetDelay(int(X))
		case 0x18:
			Println("SETSOUNDTIMER", X)
			m.SetSoundTimer(int(X))
		case 0x1E:
			Println("ADDI", X)
			m.AddI(int(X))
		case 0x0A:
			Println("WAITFORKEY", X)
			m.WaitForKey(int(X))
		case 0x29:
			Println("GETFONT", X)
			m.GetFont(int(X))
		case 0x33:
			Println("GETDIGITS", X)
			m.GetDigits(int(X))
		case 0x55:
			Println("STOREREGS", X)
			m.StoreRegs(int(X))
		case 0x65:
			Println("LOADREGS", X)
			m.LoadRegs(int(X))
		default:
			fmt.Printf("Unknown instruction: %x", o)
		}
	default:
		fmt.Printf("Unknown instruction: %x", o)
	}
}

// ClearScreen clears the screen to black.
// 00E0
func (m *Machine) ClearScreen() {
	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			m.Display[i][j] = false
		}
	}
}

// Jump jumps to the instruction at NNN
// 1NNN
func (m *Machine) Jump(a addr) {
	m.PC = a
}

// JumpSub Jumps to the subroutine at NNN
// 2NNN
func (m *Machine) JumpSub(a addr) {
	m.stack.Push(m.PC)
	m.PC = a
}

// JumpOffset Jumps to the address (NNN + offset stored in V0).
// Original interpreter version.
// BNNN
func (m *Machine) JumpOffset(a addr) {
	m.PC = a + addr(m.V[0])
}

// RetSub returns from the current subroutine.
// 00EE
func (m *Machine) RetSub() {
	m.PC = m.stack.Pop()
}

// SkipIfReg will skip the next instruction if the value in V<X> equals NN.
// 3XNN
func (m *Machine) SkipIfReg(reg int, val byte) {
	if m.V[reg] == val {
		m.PC += 2
	}
}

// SkipIfNotReg will skip the next instruction if the value in V<X> doesn't equal NN.
// 4XNN
func (m *Machine) SkipIfNotReg(reg int, val byte) {
	if m.V[reg] != val {
		m.PC += 2
	}
}

// SkipIfRegs will skip the next instruction if the values in V<X> and V<Y> are equal.
// 5XY0
func (m *Machine) SkipIfRegs(reg1, reg2 int) {
	if m.V[reg1] == m.V[reg2] {
		m.PC += 2
	}
}

// SkipIfNotRegs will skip the next instruction if the values in V<X> and V<Y> are equal.
// 9XY0
func (m *Machine) SkipIfNotRegs(reg1, reg2 int) {
	if m.V[reg1] != m.V[reg2] {
		m.PC += 2
	}
}

// SkipIfPressed will skip the next instruction if the key corresponding to the hex value in V<X> if pressed.
// EX9E
func (m *Machine) SkipIfPressed(reg int) {
	if m.keys.pressed[m.V[reg]] {
		m.PC += 2
	}
}

// SkipIfNotPressed will skip the next instruction if the key corresponding to the hex value in V<X> if pressed.
// EXA1
func (m *Machine) SkipIfNotPressed(reg int) {
	if !m.keys.pressed[m.V[reg]] {
		m.PC += 2
	}
}

// SetReg Sets the V<X> register to NN
// 6XNN
func (m *Machine) SetReg(reg int, val byte) {
	m.V[reg] = val
}

// SetRegs sets register V<X> to the value in V<Y>
// 8XY0
func (m *Machine) SetRegs(toReg, fromReg int) {
	m.V[toReg] = m.V[fromReg]
}

// OrRegs stores the results of an OR between the values of V<X> and V<Y> in V<X>.
// 8XY1
func (m *Machine) OrRegs(toReg, fromReg int) {
	m.V[toReg] = m.V[toReg] | m.V[fromReg]
}

// AndRegs stores the results of an AND between the values of V<X> and V<Y> in V<X>.
// 8XY2
func (m *Machine) AndRegs(toReg, fromReg int) {
	m.V[toReg] = m.V[toReg] & m.V[fromReg]
}

// XorRegs stores the results of an XOR between the values of V<X> and V<Y> in V<X>.
// 8XY3
func (m *Machine) XorRegs(toReg, fromReg int) {
	m.V[toReg] = m.V[toReg] ^ m.V[fromReg]
}

// Add stores the addition of the values in V<X> and V<Y> in V<X>. VF is set to 1 if overflow occurs.
// 8XY4
func (m *Machine) Add(toReg, fromReg int) {
	m.V[15] = 0
	if uint8(toReg) > 255-uint8(fromReg) {
		m.V[15] = 1
	}
	m.V[toReg] += m.V[fromReg]
}

// SubFromX stores the subtraction of V<Y> from V<X> in V<X>. VF is set to 1 if V<X> is greater than V<Y>.
// 8XY5
func (m *Machine) SubFromX(fromReg, toReg int) {
	m.V[15] = 1
	if m.V[toReg] > m.V[fromReg] {
		m.V[15] = 0
	}
	m.V[fromReg] -= m.V[toReg]
}

// SubFromY stores the subtraction of V<X> from V<Y> in V<X>. VF is set to 1 if V<Y> is greater than V<X>.
// 8XY7
func (m *Machine) SubFromY(toReg, fromReg int) {
	m.V[15] = 1
	if m.V[toReg] > m.V[fromReg] {
		m.V[15] = 0
	}
	m.V[toReg] = m.V[fromReg] - m.V[toReg]
}

// ShiftRight stores the value of V<Y> shifted one bit right in V<X>. VF is set to the bit that was shifted out.
// 8XY6
func (m *Machine) ShiftRight(toReg, shiftReg int) {
	m.V[15] = m.V[shiftReg] & 1
	m.V[toReg] = m.V[shiftReg] >> 1
}

// ShiftLeft stores the value of V<Y> shifted one bit left in V<X>. VF is set to the bit that was shifted out.
// 8XYE
func (m *Machine) ShiftLeft(toReg, shiftReg int) {
	m.V[15] = (m.V[shiftReg] & 0b10000000) >> 7
	m.V[toReg] = m.V[shiftReg] << 1
}

// AddReg Adds NN to the the V<X> register.
// 7XNN
func (m *Machine) AddReg(reg int, val byte) {
	m.V[reg] += val
}

// SetI sets the I (index) register to NNN.
// ANNN
func (m *Machine) SetI(val uint16) {
	m.I = val
}

// AddI adds the value stored in V<X> to the I register. VF is set to 1 if I overflows above 0xFFF (4096).
// FX1E
func (m *Machine) AddI(reg int) {
	m.V[15] = 0
	m.I += uint16(m.V[reg])
	if m.I > 0xFFF {
		m.V[15] = 1
	}
}

// StoreRegs Stores each V register's value up to V<X> in memory starting from the address stored in I.
// If m.IncrementLoadStore is true, I will be incremented as each register is stored.
// FX55
func (m *Machine) StoreRegs(upTo int) {
	I := int(m.I)
	for i := 0; i <= upTo; i++ {
		m.Mem[I+i] = m.V[i]
		if m.IncrementLoadStore {
			m.I++
		}
	}
}

// LoadRegs Loads each V register's value up to V<X> from memory starting from the address stored in I.
// If m.IncrementLoadStore is true, I will be incremented as each register is loaded.
// FX65
func (m *Machine) LoadRegs(upTo int) {
	I := int(m.I)
	for i := 0; i <= upTo; i++ {
		m.V[i] = m.Mem[I+i]
		if m.IncrementLoadStore {
			m.I++
		}
	}
}

// RandAnd stores the result of an AND between NN and a random number in V<X>
// CXNN
func (m *Machine) RandAnd(reg int, val byte) {
	m.V[reg] = byte(rand.Intn(255)) & val
}

// Draw draws a 8-pixel wide, N-pixel high sprite from the address stored in I at the x-y coords stored in X & Y.
// DXYN
func (m *Machine) Draw(xa, ya, rows int) {
	x := int(m.V[xa]) % w
	y := int(m.V[ya]) % h
	m.V[15] = 0
	for i := 0; i < rows; i++ {
		if m.I+uint16(i) >= 4096 {
			continue
		}
		sprite := m.Mem[m.I+uint16(i)]
		shift := 7
		and := byte(255)
		for xo := 0; xo < 8; xo++ {
			val := (sprite & and) >> shift
			shift--
			and = and >> 1
			if val == 1 {
				// Println("SET", x+xo, y+i)
				if x+xo >= w || y+i >= h {
					continue
				}
				c := m.Display[x+xo][y+i]
				m.WriteDisplay(x+xo, y+i)
				if c {
					m.V[15] = 1
				}
			}
		}
	}
}

// ReadDelay stores the current value of the delay timer in V<X>.
// FX07
func (m *Machine) ReadDelay(reg int) {
	m.V[reg] = m.delay
}

// SetDelay sets the delay timer to the value stored in V<X>.
// FX15
func (m *Machine) SetDelay(reg int) {
	m.delay = m.V[reg]
}

// SetSoundTimer sets the sound timer to the value in V<X>.
// FX18
func (m *Machine) SetSoundTimer(reg int) {
	m.sound = m.V[reg]
}

// WaitForKey blocks until a key is pressed, then stores the hex value in V<X>.
// FX0A
func (m *Machine) WaitForKey(reg int) {
	m.keys.useChan = true
	for v := range m.keys.presses {
		if v&0xF0 == 0xF0 {
			m.V[reg] = v - 0xF0
			break
		}
	}
	m.keys.useChan = false
}

// GetFont stores in I the address of the character sprite corresponding to to the hex value stored in V<X>.
// FX29
func (m *Machine) GetFont(reg int) {
	m.I = uint16(m.fontMap[m.V[reg]])
}

// GetDigits converts the number stored in V<X> and converts it to a three digit decimal representation, stored separately at (mem[I], mem[I+1], mem[I+2]).
// FX33
func (m *Machine) GetDigits(reg int) {
	n := uint8(m.V[reg])
	h := n / 100
	t := (n - (h * 100)) / 10
	d := n - (h * 100) - (t * 10)
	m.Mem[m.I] = h
	m.Mem[m.I+1] = t
	m.Mem[m.I+2] = d
}
