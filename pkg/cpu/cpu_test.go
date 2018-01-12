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
	for g := 0; g < 32; g++ {
		assert.Equal(int64(0x00000000), cpu.Vram[g])
	}

	// CPU should init with ShouldDraw = false
	assert.False(cpu.ShouldDraw)

	// Turn on some pixels.
	cpu.Vram[1] = 0x00000001
	cpu.Vram[2] = 0x00000001
	cpu.Vram[3] = 0x00000001
	cpu.Vram[4] = 0x00000001

	// Clear the Gfx buffer
	cpu.ClearVram()

	// Make sure everything is off again
	for g := 0; g < 32; g++ {
		assert.Equal(int64(0x00000000), cpu.Vram[g])
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

	// Turn on some pixels.
	cpu.Vram[1] = 0x00000001
	cpu.Vram[2] = 0x00000001
	cpu.Vram[3] = 0x00000001
	cpu.Vram[4] = 0x00000001

	// Load the opcode, and then process it.
	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)

	for g := 0; g < 32; g++ {
		assert.Equal(int64(0x00000000), cpu.Vram[g])
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

// Test 0x2NNN: Call subroutine at 0xNNN.
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

// Test 0x3XNN: Skip next instruction if VX == NN.
func Test3xnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x3A, 0x22}
	cpu := NewCpu(g, r, false)

	// Check that the program counter advances as normal if the register is not
	// set to the specified value.
	cpu.GetOp()
	err := cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x202), cpu.PC)

	// Set the register value and back up the program counter.
	cpu.Registers[0xA] = uint8(0x22)
	cpu.PC = uint16(0x200)

	// Check that the program counter advances by 4 bytes if the register matches
	// the specified value.
	cpu.GetOp()
	err = cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x204), cpu.PC)
}

// Test 0x4XNN: Skip next instruction if VX != NN.
func Test4xnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x4A, 0x22}
	cpu := NewCpu(g, r, false)

	// Check that the program counter advances by 4 bytes if the register is not
	// set to the specified value.
	cpu.GetOp()
	err := cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x204), cpu.PC)

	// Set the register value and back up the program counter.
	cpu.Registers[0xA] = uint8(0x22)
	cpu.PC = uint16(0x200)

	// Check that the program counter advances as normal if the register is
	// set to the specified value.
	cpu.GetOp()
	err = cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x202), cpu.PC)
}

// Test 0x5XY0: Skip next instruction if VX == VY.
func Test5xnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x5A, 0x10}
	cpu := NewCpu(g, r, false)

	// Check that the program counter advances by 4 bytes since the registers
	// match by default (0x00)
	cpu.GetOp()
	err := cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x204), cpu.PC)

	// Move the PC back and change one of the register values.
	cpu.PC = uint16(0x200)
	cpu.Registers[int(0xA)] = uint8(0xFF)

	// Check that the program counter advances by 2 bytes when the registers
	// don't match.
	cpu.GetOp()
	err = cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint16(0x202), cpu.PC)
}

// Test 0x6XNN: Set VX to NN.
func Test6xnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0x6A, 0xFF}
	cpu := NewCpu(g, r, false)

	cpu.GetOp()
	err := cpu.ProcessOpcode()
	assert.NoError(err)
	assert.Equal(uint8(0xFF), cpu.Registers[0xA])
}

// Test 0xANNN: Set index register to 0xNNN.
func TestAnnn(t *testing.T) {
	assert := asrt.New(t)

	g := &graphics.Noop{}
	r := []byte{0xA2, 0x34}
	cpu := NewCpu(g, r, false)

	// Make sure that the CPU state is good before processing the opcode.
	assert.Equal(uint16(0x0000), cpu.IndexRegister)

	cpu.GetOp()
	err := cpu.ProcessOpcode()

	assert.NoError(err)
	assert.Equal(uint16(0x234), cpu.IndexRegister)
}
