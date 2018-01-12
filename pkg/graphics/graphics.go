package graphics

import (
	"github.com/gen2brain/raylib-go/raylib"
	termbox "github.com/nsf/termbox-go"
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
	case "raylib":
	default:
		g := &Raylib{}
		g.Init()
		return g
	}

	return nil
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
	raylib.BeginDrawing()
	raylib.ClearBackground(raylib.Black)
	for i := 0; i < 64*32; i++ {
		// if buf[i] {
		// 	xrootpos := (i % 64)
		// 	yrootpos := math.Floor(float64((i - (i % 64)) / 64))

		// 	// Everything is 16x so that we will be able to see it on modern displays.
		// 	yrootpos = yrootpos * 8
		// 	xrootpos = xrootpos * 8

		// 	raylib.DrawRectangle(int32(xrootpos), int32(yrootpos), 8, 8, raylib.White)
		// }
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
