package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cweagans/chip8/pkg/cpu"
	"github.com/cweagans/chip8/pkg/ui"
)

var (
	RomFile    string
	UIMode     string
	Debug      bool
	ClockSpeed int
)

func init() {
	flag.StringVar(&UIMode, "ui", "sdl", "Which UI should the emulator use? Options: sdl (default), termbox.")
	flag.BoolVar(&Debug, "debug", false, "Set debug to true if you want to log CPU internals")
	flag.IntVar(&ClockSpeed, "clock-speed", 60, "Set the CPU clock speed (in Hertz).")
	flag.StringVar(&RomFile, "rom", "", "Set the ROM filename that the emulator will load.")
	flag.Parse()

	if RomFile == "" {
		fmt.Println("-rom flag is required.")
		os.Exit(1)
	}
}

func main() {
	// Get a UI object.
	u := ui.GetUI(UIMode)

	// Load ROM to pass to CPU.
	rom, err := loadRom(RomFile)
	if err != nil {
		fmt.Println("Could not open specified ROM file: " + err.Error())
	}

	// Create a new CPU.
	c := cpu.NewCpu(u, rom, Debug)

	// Set the clock speed based on input.
	c.SetClockSpeed(ClockSpeed)

	// Run the CPU.
	c.Run()
}

func loadRom(filename string) ([]byte, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return content, nil
}
