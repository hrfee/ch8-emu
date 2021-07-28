package main

import (
	"fmt"

	"github.com/tfriedel6/canvas/sdlcanvas"
)

var CANVAS = true

func (m *Machine) draw() {
	wnd, cv, err := sdlcanvas.CreateWindow(w*m.scale, h*m.scale, "schoolasm")
	if err != nil {
		panic(fmt.Errorf("Failed to create window: %v", err))
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
