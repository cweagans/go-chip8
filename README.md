# CHIP-8 emulator

[![Build Status](https://travis-ci.org/cweagans/chip8.svg?branch=master)](https://travis-ci.org/cweagans/chip8)

As a learning project, I decided to write an emulator in Go. It seems like CHIP-8
is generally agreed to be the easiest starting point, so I started there.

## Usage

Right now, the emulator is not functional (unless your ROM just clears the screen
repeatedly).

When it's functional, you might be interested in the CHIP-8 program pack, which
can be found here: https://web.archive.org/web/20130903155600/http://chip8.com/?page=109

## Reference material

* [How to write an emulator (CHIP-8 interpreter)](http://www.multigesture.net/articles/how-to-write-an-emulator-chip-8-interpreter/)
* [Wikipedia: CHIP-8](https://en.wikipedia.org/wiki/CHIP-8)
* [Cowgod's CHIP-8 tech reference](http://devernay.free.fr/hacks/chip8/C8TECH10.HTM)

## Opcodes:

| Implemented | Opcode | Description |
| --- | --- | --- |
| ❌ | `0x0NNN` | Calls RCA 1802 program at address 0xNNN. |
| ✅ | `0x00E0` | Clear the screen |
| ✅ | `0x00EE` | Return from subroutine |
| ✅ | `0x1NNN` | Jump to 0xNNN |
| ✅ | `0x2NNN` | Call subroutine at 0xNNN  |
| ✅ | `0x3XNN` | Skip next instruction if `VX == NN` |
| ✅ | `0x4XNN` | Skip next instruction if `VX != NN` |
| ✅ | `0x5XY0` | Skip next instruction if `VX == VY` |
| ✅ | `0x6XNN` | Set `VX` to `NN` |
| ✅ | `0x7XNN` | Add `NN` to `VX` (carry flag is not changed) |
| ❌ | `0x8XY0` | Set `VX` to the value of `VY` |
| ❌ | `0x8XY1` | Set `VX` to `VX \| VY` (bitwise OR) |
| ❌ | `0x8XY2` | Set `VX` to `VX & VY` (bitwise AND)|
| ❌ | `0x8XY3` | Set `VX` to `VX xor VY` |
| ❌ | `0x8XY4` | Add `VY` to `VX`. `VF` is set to 1 when there's a carry, and 0 when there isn't. |
| ❌ | `0x8XY5` | Subtract `VY` from `VX`. `VF` is set to 0 when there's a borrow and 1 when there isn't. |
| ❌ | `0x8XY6` | Shift `VY` right by one and copy the result to `VX`. `VF` is set to the value of the least significant bit of `VY` before the shift. |
| ❌ | `0x8XY7` | Set `VX` to `VY - VX`. `VF` is set to 0 when there's a borrow and 1 when there isn't. |
| ❌ | `0x8XYE` | Shift `VY` left by one and copy the result to `VX`. `VF` is set to the value of the least significant bit of `VY` before the shift. |
| ❌ | `0x9XY0` | Skip the next instruction if `VX` doesn't equal `VY`. |
| ✅ | `0xANNN` | Set index register to 0xNNN |
| ❌ | `0xBNNN` | Jump to the address `NNN` plus `V0` |
| ❌ | `0xCXNN` | Set `VX` to the result of a bitwise and operation on a random number and `NN` |
| ❌ | `0xDXYN` | Draw a sprite at coordinate (`VX`, `VY`) that has a width of 8 pixels and a height of `N` pixels. Each row is read as bit-coded starting from the index register, I. I doesn't change after the execution of this instruction. `VF` is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that doesn't happen. |
| ❌ | `0xEX9E` | Skip the next instruction if the key stored in `VX` is pressed. |
| ❌ | `0xEXA1` | Skip the next instruction if the key stored in `VX` is not pressed. |
| ❌ | `0xFX07` | Set `VX` to the value of the delay timer. |
| ❌ | `0xFX0A` | A key press is awaited and then stored in `VX` (blocking operation - all instructions are halted until the next key event) |
| ❌ | `0xFX15` | Set the delay timer to the value of `VX` |
| ❌ | `0xFX18` | Set the sound timer to the value of `VX` |
| ❌ | `0xFX1E` | Add the value of `VX` to the index register |
| ❌ | `0xFX29` | Set `I` to the location of the sprite for the character in `VX`. Characters 0-F (in hex) are represented by a 4x5 font. |
| ❌ | `0xFX33` | Stores the binary-coded decimal representation of `VX`, with the most significant of three digits at the address in `I`, the middle digit at `I` plus 1, and the least significant digit at `I` plus 2. (In other words, take the decimal representation of VX, place the hundreds digit in memory at location in `I`, the tens digit at location `I+1`, and the ones digit at location `I+2`.) |
| ❌ | `0xFX55` | Stores `V0` to `VX` (including `VX`) in memory starting at address `I`. `I` is increased by 1 for each value written. |
| ❌ | `0xFX65` | Fills `V0` to `VX` (including `VX`) with values from memory starting at address `I`. `I` is increased by 1 for each value written. |
