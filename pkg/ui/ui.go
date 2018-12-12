package ui

import "fmt"

// Graphics objects allow the emulator to draw to the screen in a backend
// agnostic way.
// type Graphics interface {
// 	Init()
// 	Draw([32]int64)
// 	Shutdown()
// }

// UI objects allow the emulator to draw to the screen, get input, etc. in a
// backend agnostic way.
type UI interface {
	Init()
	Draw([32]int64)
	GetInput() Input
	Shutdown()
}

// Input is a way to pass input state back to the CPU.
type Input struct {
	Key0   bool
	Key1   bool
	Key2   bool
	Key3   bool
	Key4   bool
	Key5   bool
	Key6   bool
	Key7   bool
	Key8   bool
	Key9   bool
	KeyA   bool
	KeyB   bool
	KeyC   bool
	KeyD   bool
	KeyE   bool
	KeyF   bool
	KeyEsc bool
}

// GetUI returns an initialized UI object.
func GetUI(UIType string) UI {
	switch UIType {
	case "termbox":
		u := &Termbox{}
		u.Init()
		return u
	case "noop":
		u := &Noop{}
		u.Init()
		return u
	case "sdl":
		fallthrough
	default:
		u := &Sdl{}
		u.Init()
		return u
	}
}

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
