package cpu

import (
	"fmt"
	"time"

	"github.com/cweagans/chip8/pkg/graphics"
)

// Cpu is the core model of the system.
type Cpu struct {
	Graphics      graphics.Graphics
	Vram          [64 * 32]bool
	ShouldDraw    bool
	ClockSpeed    int
	Memory        [4096]byte
	PC            uint16
	Op            uint16
	ShouldHalt    bool
	Stack         [16]uint16
	StackPointer  int
	Debug         bool
	Registers     [16]uint8
	IndexRegister uint16
	DelayTimer    uint8
	SoundTimer    uint8
	Keys          [16]uint8
}

// UnknownOpcodeError is returned when the CPU encounters an opcode that it does
// not know how to process.
type UnknownOpcodeError struct {
	Opcode  uint16
	Address uint16
}

func (uoe *UnknownOpcodeError) Error() string {
	return fmt.Sprintf("Unknown opcode 0x%X at address 0x%X", uoe.Opcode, uoe.Address)
}

// InitCpu() sets up a new CPU and loads the rom into memory.
func NewCpu(g graphics.Graphics, r []byte, debug bool) *Cpu {
	cpu := &Cpu{}
	cpu.Graphics = g
	cpu.PC = 0x200
	cpu.ShouldDraw = false
	cpu.ShouldHalt = false
	cpu.Debug = debug
	cpu.ClockSpeed = 60
	cpu.IndexRegister = 0x0000

	cpu.LoadRom(r)

	return cpu
}

// Set the CPU clock speed.
func (c *Cpu) SetClockSpeed(s int) {
	c.ClockSpeed = s
}

// Loads the supplied ROM bytes into memory starting at 0x200.
func (c *Cpu) LoadRom(r []byte) {
	// Clear memory.
	for m := 0; m < 4096; m++ {
		c.Memory[m] = 0x00
	}

	// Copy program into memory starting at 0x200.
	for index, b := range r {
		c.Memory[index+0x200] = b
	}
}

// Clear Vram.
func (c *Cpu) ClearVram() {
	for g := 0; g < (64 * 32); g++ {
		c.Vram[g] = false
	}

	c.ShouldDraw = true
}

// Runs the CPU until halted.
func (c *Cpu) Run() {

	// c.ClockSpeed defaults to 60 Hz, but this can be adjusted as needed for debugging.
	for range time.Tick(time.Duration(1000/c.ClockSpeed) * time.Millisecond) {
		// Get the next opcode.
		c.GetOp()

		// If GetOp() couldn't find another opcode, then it will set the ShouldHalt flag.
		if c.ShouldHalt {
			break
		}

		// Process the current opcode.
		err := c.ProcessOpcode()
		if err != nil {
			// @TODO: Is there a better way to handle this? It's not really something
			// that can be recovered from gracefully. Does it need to bring down the
			// entire emulator though?
			panic(err.Error())
		}

		// If ShouldDraw has been set, we need to update the screen.
		if c.ShouldDraw {
			c.Graphics.Draw(c.Vram)
		}

		// @TODO: Get input state.
	}
}

func (c *Cpu) DumpMemory() {
	fmt.Println("Address\tValue")
	for m := 0; m < 4096; m++ {
		fmt.Printf("0x%X\t0x%X\n", m, c.Memory[m])
	}
}

// Get next opcode.
func (c *Cpu) GetOp() {
	oldOp := c.Op

	// An opcode is two bytes, starting at c.PC. The first byte is bitshift-ed to the left,
	// and then ORed with the second byte. The end result is a 16 bit opcode.
	c.Op = (uint16(c.Memory[c.PC]) << 8) | uint16(c.Memory[c.PC+1])

	if c.Debug && (oldOp != c.Op) {
		fmt.Printf("New opcode: 0x%X\n", c.Op)
		fmt.Printf("First byte: 0x%X\n", c.Memory[c.PC])
		fmt.Printf("Second byte: 0x%X\n", c.Memory[c.PC+1])
	}

	if c.Op == 0x0000 {
		c.ShouldHalt = true
	}
}

// Process the current opcode.
func (c *Cpu) ProcessOpcode() error {

	opcodeFound := false

	// Start by reading the first four bits of the opcode.
	switch c.Op & 0xF000 {

	case 0x0000:
		switch c.Op & 0x000F {
		case 0x0000:
			// 0x00E0: Clear the screen and advance to the next opcode.
			opcodeFound = true
			c.ClearVram()
			c.PC += 2
			break

		case 0x000E:
			// 0x00EE: Returns from a subroutine.
			opcodeFound = true
			c.StackPointer -= 1
			c.PC = c.Stack[c.StackPointer]
			c.Stack[c.StackPointer] = 0
			break
		}

	case 0x1000:
		// 0x1NNN: Jump to 0xNNN
		opcodeFound = true
		c.PC = c.Op & 0x0FFF
		break

	case 0x2000:
		// 0x2NNN: Call subroutine at 0xNNN
		opcodeFound = true
		c.Stack[c.StackPointer] = c.PC
		c.StackPointer += 1
		c.PC = c.Op & 0x0FFF
		break

	case 0xA000:
		// 0xANNN: Set index register to 0xNNN.
		opcodeFound = true
		c.IndexRegister = c.Op & 0x0FFF
		c.PC += 2
	}

	// If we didn't find a way to process the opcode, return an error.
	if !opcodeFound {
		return &UnknownOpcodeError{
			Opcode:  c.Op,
			Address: c.PC,
		}
	}

	return nil
}
