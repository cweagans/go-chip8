
bin:
	go build -i -o ./chip8 github.com/cweagans/chip8/cmd/chip8

test:
	go test -v github.com/cweagans/chip8/pkg/cpu github.com/cweagans/chip8/pkg/graphics

lint:
	go vet github.com/cweagans/chip8/pkg/cpu github.com/cweagans/chip8/pkg/graphics github.com/cweagans/chip8/cmd/chip8
