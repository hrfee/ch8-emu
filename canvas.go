package main

import (
	"fmt"

	"github.com/tfriedel6/canvas/sdlcanvas"
	"github.com/veandco/go-sdl2/sdl"
)

var CANVAS = true

const (
	// - PAD - == - KB -
	// 1 2 3 C == 1 2 3 4
	// 4 5 6 D == Q W E R
	// 7 8 9 E == A S D F
	// A 0 B F == Z X C V
	k1 = sdl.SCANCODE_1
	k2 = sdl.SCANCODE_2
	k3 = sdl.SCANCODE_3
	kC = sdl.SCANCODE_4
	k4 = sdl.SCANCODE_Q
	k5 = sdl.SCANCODE_W
	k6 = sdl.SCANCODE_E
	kD = sdl.SCANCODE_R
	k7 = sdl.SCANCODE_A
	k8 = sdl.SCANCODE_S
	k9 = sdl.SCANCODE_D
	kE = sdl.SCANCODE_F
	kA = sdl.SCANCODE_Z
	k0 = sdl.SCANCODE_X
	kB = sdl.SCANCODE_C
	kF = sdl.SCANCODE_V
)

var codeToHex = map[int]byte{
	k1: 0x1,
	k2: 0x2,
	k3: 0x3,
	kC: 0xC,
	k4: 0x4,
	k5: 0x5,
	k6: 0x6,
	kD: 0xD,
	k7: 0x7,
	k8: 0x9,
	k9: 0x0,
	kE: 0xE,
	kA: 0xA,
	k0: 0x0,
	kB: 0xB,
	kF: 0xF,
}

// Keys stores the number of keys pressed, and which.
type Keys struct {
	pressedCount int
	pressed      map[uint8]bool
	// The 4 bits of the key hexcode + (0xF0 if pressed, 0x00 if not).
	presses chan byte
}

func newKeys() Keys {
	return Keys{
		pressedCount: 0,
		pressed: map[byte]bool{
			0x0: false,
			0x1: false,
			0x2: false,
			0x3: false,
			0x4: false,
			0x5: false,
			0x6: false,
			0x7: false,
			0x8: false,
			0x9: false,
			0xA: false,
			0xB: false,
			0xC: false,
			0xD: false,
			0xE: false,
			0xF: false,
		},
		presses: make(chan byte),
	}
}

func loadDefaultFont(m *Machine) {
	// Font from https://tobiasvl.github.io/blog/write-a-chip-8-emulator/#font
	m.Mem[0x050] = 0xF0
	m.Mem[0x051] = 0x90
	m.Mem[0x052] = 0x90
	m.Mem[0x053] = 0x90
	m.Mem[0x054] = 0xF0 // 0
	m.Mem[0x055] = 0x20
	m.Mem[0x056] = 0x60
	m.Mem[0x057] = 0x20
	m.Mem[0x058] = 0x20
	m.Mem[0x059] = 0x70 // 1
	m.Mem[0x05A] = 0xF0
	m.Mem[0x05B] = 0x10
	m.Mem[0x05C] = 0xF0
	m.Mem[0x05D] = 0x80
	m.Mem[0x05E] = 0xF0 // 2
	m.Mem[0x05F] = 0xF0
	m.Mem[0x060] = 0x10
	m.Mem[0x061] = 0xF0
	m.Mem[0x062] = 0x10
	m.Mem[0x063] = 0xF0 // 3
	m.Mem[0x064] = 0x90
	m.Mem[0x065] = 0x90
	m.Mem[0x066] = 0xF0
	m.Mem[0x067] = 0x10
	m.Mem[0x068] = 0x10 // 4
	m.Mem[0x069] = 0xF0
	m.Mem[0x06A] = 0x80
	m.Mem[0x06B] = 0xF0
	m.Mem[0x06C] = 0x10
	m.Mem[0x06D] = 0xF0 // 5
	m.Mem[0x06E] = 0xF0
	m.Mem[0x06F] = 0x80
	m.Mem[0x070] = 0xF0
	m.Mem[0x071] = 0x90
	m.Mem[0x072] = 0xF0 // 6
	m.Mem[0x073] = 0xF0
	m.Mem[0x074] = 0x10
	m.Mem[0x075] = 0x20
	m.Mem[0x076] = 0x40
	m.Mem[0x077] = 0x40 // 7
	m.Mem[0x078] = 0xF0
	m.Mem[0x079] = 0x90
	m.Mem[0x07A] = 0xF0
	m.Mem[0x07B] = 0x90
	m.Mem[0x07C] = 0xF0 // 8
	m.Mem[0x07D] = 0xF0
	m.Mem[0x07E] = 0x90
	m.Mem[0x07F] = 0xF0
	m.Mem[0x080] = 0x10
	m.Mem[0x081] = 0xF0 // 9
	m.Mem[0x082] = 0xF0
	m.Mem[0x083] = 0x90
	m.Mem[0x084] = 0xF0
	m.Mem[0x085] = 0x90
	m.Mem[0x086] = 0x90 // A
	m.Mem[0x087] = 0xE0
	m.Mem[0x088] = 0x90
	m.Mem[0x089] = 0xE0
	m.Mem[0x08A] = 0x90
	m.Mem[0x08B] = 0xE0 // B
	m.Mem[0x08C] = 0xF0
	m.Mem[0x08D] = 0x80
	m.Mem[0x08E] = 0x80
	m.Mem[0x08F] = 0x80
	m.Mem[0x090] = 0xF0 // C
	m.Mem[0x091] = 0xE0
	m.Mem[0x092] = 0x90
	m.Mem[0x093] = 0x90
	m.Mem[0x094] = 0x90
	m.Mem[0x095] = 0xE0 // D
	m.Mem[0x096] = 0xF0
	m.Mem[0x097] = 0x80
	m.Mem[0x098] = 0xF0
	m.Mem[0x099] = 0x80
	m.Mem[0x09A] = 0xF0 // E
	m.Mem[0x09B] = 0xF0
	m.Mem[0x09C] = 0x80
	m.Mem[0x09D] = 0xF0
	m.Mem[0x09E] = 0x80
	m.Mem[0x09F] = 0x80 // F
	m.fontMap = map[byte]addr{
		0x0: 0x050,
		0x1: 0x055,
		0x2: 0x05A,
		0x3: 0x05F,
		0x4: 0x064,
		0x5: 0x069,
		0x6: 0x06E,
		0x7: 0x073,
		0x8: 0x078,
		0x9: 0x07D,
		0xa: 0x082,
		0xb: 0x087,
		0xc: 0x08C,
		0xd: 0x091,
		0xe: 0x096,
		0xf: 0x09B,
	}
}

func (m *Machine) draw() {
	wnd, cv, err := sdlcanvas.CreateWindow(w*m.scale, h*m.scale, "ch8-emu")
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %v", err))
	}
	wnd.KeyDown = func(scancode int, rn rune, name string) {
		switch scancode {
		case k1, k2, k3, kC, k4, k5, k6, kD, k7, k8, k9, kE, kA, k0, kB, kF:
			m.keys.pressedCount++
			m.keys.pressed[codeToHex[scancode]] = true
			m.keys.presses <- (byte(codeToHex[scancode]) + 0xF0)
		}
	}
	wnd.KeyUp = func(scancode int, rn rune, name string) {
		switch scancode {
		case k1, k2, k3, kC, k4, k5, k6, kD, k7, k8, k9, kE, kA, k0, kB, kF:
			m.keys.pressedCount--
			m.keys.pressed[codeToHex[scancode]] = false
			m.keys.presses <- byte(codeToHex[scancode])
		}
		if m.keys.pressedCount < 0 {
			panic(fmt.Errorf("Number of keys pressed was negative"))
		}
	}

	defer wnd.Destroy()
	wnd.MainLoop(func() {
		for sx := 0; sx < w; sx++ {
			x := sx * m.scale
			for sy := 0; sy < h; sy++ {
				y := sy * m.scale
				white := m.Display[sx][sy]
				size := float64(m.scale)
				cv.SetFillStyle("#000000")
				cv.FillRect(float64(x), float64(y), size, size)
				if white {
					cv.SetFillStyle("#ffffff")
					size = float64(m.scale) * ((float64(m.scale) - float64(m.pad) - float64(m.pad)) / float64(m.scale))
					cv.FillRect(float64(x)+float64(m.pad), float64(y)+float64(m.pad), size, size)
				}
			}
		}
	})
}
