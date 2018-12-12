
bin:
	go build -i -o ./chip8 github.com/cweagans/chip8/cmd/chip8

vendor:
	dep ensure

test:
	go test -v github.com/cweagans/chip8/pkg/cpu github.com/cweagans/chip8/pkg/ui

lint:
	go vet github.com/cweagans/chip8/pkg/cpu github.com/cweagans/chip8/pkg/ui github.com/cweagans/chip8/cmd/chip8
