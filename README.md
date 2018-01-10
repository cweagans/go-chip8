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

## Opcodes:

| Implemented | Opcode | Description |
| --- | --- | --- |
| ❌ | `0x0NNN` |  |
| ✅ | `0x00E0` | Clear the screen |
| ✅ | `0x00EE` | Return from subroutine |
| ✅ | `0x1NNN` | Jump to 0xNNN |
| ✅ | `0x2NNN` | Call subroutine at 0xNNN  |
| ❌ | `0x3XNN` |  |
| ❌ | `0x4XNN` |  |
| ❌ | `0x5XY0` |  |
| ❌ | `0x6XNN` |  |
| ❌ | `0x7XNN` |  |
| ❌ | `0x8XY0` |  |
| ❌ | `0x8XY1` |  |
| ❌ | `0x8XY2` |  |
| ❌ | `0x8XY3` |  |
| ❌ | `0x8XY4` |  |
| ❌ | `0x8XY5` |  |
| ❌ | `0x8XY6` |  |
| ❌ | `0x8XY7` |  |
| ❌ | `0x8XYE` |  |
| ❌ | `0x9XY0` |  |
| ✅ | `0xANNN` | Set index register to 0xNNN |
| ❌ | `0xBNNN` |  |
| ❌ | `0xCXNN` |  |
| ❌ | `0xDXYN` |  |
| ❌ | `0xEX9E` |  |
| ❌ | `0xEXA1` |  |
| ❌ | `0xFX07` |  |
| ❌ | `0xFX0A` |  |
| ❌ | `0xFX15` |  |
| ❌ | `0xFX18` |  |
| ❌ | `0xFX1E` |  |
| ❌ | `0xFX29` |  |
| ❌ | `0xFX33` |  |
| ❌ | `0xFX55` |  |
| ❌ | `0xFX65` |  |
