package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cweagans/chip8/pkg/cpu"
	"github.com/cweagans/chip8/pkg/graphics"
)

var (
	RomFile      string
	GraphicsMode string
	Debug        bool
)

func init() {
	flag.StringVar(&GraphicsMode, "ui", "raylib", "Which UI should the emulator use? Options: raylib (default), termbox.")
	flag.BoolVar(&Debug, "debug", false, "Set debug to true if you want to log CPU internals")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("ROM file argument is required.")
		os.Exit(1)
	}
	RomFile = os.Args[1]
}

func main() {
	// Get a graphics object.
	g := graphics.GetGraphics(GraphicsMode)

	// Load ROM to pass to CPU.
	rom, err := loadRom(RomFile)
	if err != nil {
		fmt.Println("Could not open specified ROM file: " + err.Error())
	}

	// create cpu package
	_ = cpu.NewCpu(g, rom, Debug)

	// add termbox graphics implementation + tests for output
}

func loadRom(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}
