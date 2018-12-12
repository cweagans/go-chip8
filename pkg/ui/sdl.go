package ui

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

// Sdl will draw emulator output in a separate GUI window with SDL.
type Sdl struct {
	Window *sdl.Window
}

func (s *Sdl) Init() {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	// Create an SDL window.
	window, err := sdl.CreateWindow("Chip8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 512, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}

	// Save the window handle for later.
	s.Window = window

	// Delay a bit to make sure that the window is ready for drawing before we
	// start drawing.
	sdl.Delay(100)

	// Get the window surface,
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// fill it black, and update the surface.
	surface.FillRect(nil, 0)
	window.UpdateSurface()
}

func (s Sdl) Draw(buf [32]int64) {
	surface, err := s.Window.GetSurface()
	if err != nil {
		// @TODO: Maybe there's something better than skipping a frame?
		return
	}

	pixels := ConvertVramToBools(buf)

	surface.FillRect(nil, 0)

	for i := 0; i < 64*32; i++ {
		if pixels[i] {
			xrootpos := (i % 64)
			yrootpos := math.Floor(float64((i - (i % 64)) / 64))

			// Everything is 8x so that we will be able to see it on modern displays.
			yrootpos = yrootpos * 8
			xrootpos = xrootpos * 8

			rect := sdl.Rect{
				X: int32(xrootpos),
				Y: int32(yrootpos),
				H: 8,
				W: 8,
			}
			surface.FillRect(&rect, 0xffffffff)
		}
	}

	s.Window.UpdateSurface()
}

func (s Sdl) GetInput() Input {
	return Input{}
}

func (s Sdl) Shutdown() {
	sdl.Quit()
	s.Window.Destroy()
}
