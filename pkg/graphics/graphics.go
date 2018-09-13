package graphics

import (
	"fmt"
	"math"

	"github.com/gen2brain/raylib-go/raylib"
	termbox "github.com/nsf/termbox-go"
	"github.com/veandco/go-sdl2/sdl"
)

// Graphics objects allow the emulator to draw to the screen in a backend
// agnostic way.
type Graphics interface {
	Init()
	Draw([32]int64)
	Shutdown()
}

// GetGraphics() returns an initialized Graphics object.
func GetGraphics(graphicsType string) Graphics {
	switch graphicsType {
	case "termbox":
		g := &Termbox{}
		g.Init()
		return g
	case "noop":
		g := &Noop{}
		g.Init()
		return g
	case "sdl":
		g := &Sdl{}
		g.Init()
		return g
	case "raylib":
		fallthrough
	default:
		g := &Raylib{}
		g.Init()
		return g
	}
}

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

func (s Sdl) Shutdown() {
	sdl.Quit()
	s.Window.Destroy()
}

// Raylib will draw emulator output in a separate GUI window.
type Raylib struct{}

func (r Raylib) Init() {
	raylib.InitWindow(512, 256, "Chip8")
	raylib.BeginDrawing()
	raylib.ClearBackground(raylib.Black)
	raylib.EndDrawing()
}

func (r Raylib) Draw(buf [32]int64) {
	pixels := ConvertVramToBools(buf)

	raylib.BeginDrawing()
	raylib.ClearBackground(raylib.Black)
	for i := 0; i < 64*32; i++ {
		if pixels[i] {
			xrootpos := (i % 64)
			yrootpos := math.Floor(float64((i - (i % 64)) / 64))

			// Everything is 8x so that we will be able to see it on modern displays.
			yrootpos = yrootpos * 8
			xrootpos = xrootpos * 8

			raylib.DrawRectangle(int32(xrootpos), int32(yrootpos), 8, 8, raylib.White)
		}
	}
	raylib.EndDrawing()
}

func (r Raylib) Shutdown() {
	raylib.CloseWindow()
}

// Termbox will eventually use termbox-go to draw emulator output in a terminal window.
type Termbox struct{}

func (t Termbox) Init() {
	err := termbox.Init()
	if err != nil {
		panic(err.Error())
	}
}

func (t Termbox) Draw(buf [32]int64) {
	panic("Termbox is unimplemented.")
}

func (t Termbox) Shutdown() {
	termbox.Close()
}

// Noop will skip drawing altogether.
type Noop struct{}

func (n Noop) Init()              {}
func (n Noop) Draw(buf [32]int64) {}
func (n Noop) Shutdown()          {}

// Convert an array of int64 values to an array of booleans corresponding to the
// bits that make up the integer value. This is mostly to help with the raylib
// graphics, as it was implemented before Vram was switched to an array of int64
// values.
func ConvertVramToBools(vram [32]int64) [64 * 32]bool {

	var pixels [64 * 32]bool
	counter := 0
	for row := 0; row < 32; row++ {
		rowInt := vram[row]
		rowBinary := fmt.Sprintf("%064b", rowInt)

		for _, val := range []rune(rowBinary) {
			switch string(val) {
			case "0":
				pixels[counter] = false
				break
			case "1":
				pixels[counter] = true
				break
			}

			counter += 1
		}
	}

	return pixels
}
