package cpu

import (
	"testing"

	"github.com/cweagans/chip8/pkg/graphics"
	asrt "github.com/stretchr/testify/assert"
)

// Test that a minimal CPU initialization works.
//   - ROM is loaded into RAM at 0x200
//   - PC is set to 0x200
func TestNewCpu(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0xff, 0xff}
	cpu := NewCpu(g, r, false)

	assert.Equal(uint8(0xff), cpu.Memory[0x200])
	assert.Equal(uint8(0xff), cpu.Memory[0x201])
	assert.Equal(uint16(0x200), cpu.PC)
}

// Test that ClearGfx clears the graphics buffer + sets ShouldDraw.
func TestClearGfx(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{}
	cpu := NewCpu(g, r, false)

	// CPU should init with an empty Gfx buffer
	for g := 0; g < (64 * 32); g++ {
		assert.False(cpu.Vram[g])
	}

	// CPU should init with ShouldDraw = false
	assert.False(cpu.ShouldDraw)

	// Set some pixels to true.
	cpu.Vram[1] = true
	cpu.Vram[2] = true
	cpu.Vram[3] = true
	cpu.Vram[4] = true

	// Clear the Gfx buffer
	cpu.ClearVram()

	// Make sure everything is off again
	for g := 0; g < (64 * 32); g++ {
		assert.False(cpu.Vram[g])
	}

	// After clearing the Gfx buffer, the CPU should know to draw to the screen.
	assert.True(cpu.ShouldDraw)
}

// Test that GetOp() finds and constructs opcodes properly.
func TestGetOp(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x00, 0xE0}
	cpu := NewCpu(g, r, false)

	// The CPU shouldn't have an opcode loaded before running.
	assert.Equal(uint16(0x0000), cpu.Op)

	// We've loaded two bytes into memory at 0x200, and cpu.PC starts at 0x200,
	// so the end result should be 0x00E0.
	cpu.GetOp()
	assert.Equal(uint16(0x00E0), cpu.Op)

	// Advancing the program counter by two bytes should be pointing at empty memory,
	// so calling GetOp() again should yield 0x0000 and c.ShouldHalt should be true.
	cpu.PC += 2
	cpu.GetOp()
	assert.Equal(uint16(0x0000), cpu.Op)
	assert.True(cpu.ShouldHalt)

	// Let's try with an opcode that doesn't start with 0x00, as that may hide
	// some errors.
	cpu.LoadRom([]byte{0x12, 0x34})
	cpu.PC = uint16(0x0200)
	cpu.GetOp()
	assert.Equal(uint16(0x1234), cpu.Op)
}

// Test 00e0: Clears the vram and sets ShouldDraw to true.
func Test00e0(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x00, 0xE0}
	cpu := NewCpu(g, r, false)

	// CPU should init with ShouldDraw = false
	assert.False(cpu.ShouldDraw)

	// Set some pixels to true.
	cpu.Vram[1] = true
	cpu.Vram[2] = true
	cpu.Vram[3] = true
	cpu.Vram[4] = true

	// Load the opcode, and then process it.
	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)

	for g := 0; g < (64 * 32); g++ {
		assert.False(cpu.Vram[g])
	}
	assert.True(cpu.ShouldDraw)
}

// Test 00EE: Return from a subroutine.
func Test00ee(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x00, 0xE0, 0x00, 0xEE}
	cpu := NewCpu(g, r, false)

	cpu.PC = 0x202
	cpu.Stack[0] = 0x200
	cpu.StackPointer = 1

	// Load the opcode, and then process it.
	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)

	assert.Equal(uint16(0x200), cpu.PC)
	assert.Equal(0, cpu.StackPointer)
	assert.Equal(uint16(0), cpu.Stack[0])
}

// Test 0x1NNN: Jump to 0xNNN.
func Test1nnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x12, 0x34}
	cpu := NewCpu(g, r, false)

	// Make sure that the CPU state is good before processing the opcode.
	assert.Equal(uint16(0x200), cpu.PC)

	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)
	assert.Equal(uint16(0x234), cpu.PC)
}

// Test 0x2NNN: Call subrouting at 0xNNN.
func Test2nnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x22, 0x34}
	cpu := NewCpu(g, r, false)

	// Make sure that the CPU state is good before processing the opcode.
	assert.Equal(0, cpu.StackPointer)
	assert.Equal(uint16(0x200), cpu.PC)

	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)
	assert.Equal(1, cpu.StackPointer)
	assert.Equal(uint16(0x234), cpu.PC)
	assert.Equal(uint16(0x200), cpu.Stack[0])
}
