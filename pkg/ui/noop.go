package ui

// Noop will skip drawing and input entirely.
type Noop struct{}

func (n Noop) Init()              {}
func (n Noop) Draw(buf [32]int64) {}
func (n Noop) GetInput() Input    { return Input{} }
func (n Noop) Shutdown()          {}
