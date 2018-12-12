package ui

import (
	"math"

	termbox "github.com/nsf/termbox-go"
)

// Termbox will eventually use termbox-go to draw emulator output in a terminal window.
type Termbox struct {
	events chan termbox.Event
}

func (t Termbox) Init() {
	err := termbox.Init()
	if err != nil {
		panic(err.Error())
	}
	termbox.SetInputMode(termbox.InputEsc)

	t.events = make(chan termbox.Event)
	go func() {
		for {
			t.events <- termbox.PollEvent()
		}
	}()
}

func (t Termbox) Draw(buf [32]int64) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	pixels := ConvertVramToBools(buf)

	for i := 0; i < 64*32; i++ {
		if pixels[i] {
			xrootpos := (i % 64)
			yrootpos := math.Floor(float64((i - (i % 64)) / 64))

			xrootpos = xrootpos * 2

			termbox.SetCell(int(xrootpos), int(yrootpos), ' ', termbox.ColorDefault, termbox.ColorWhite)
			termbox.SetCell(int(xrootpos)+1, int(yrootpos), ' ', termbox.ColorDefault, termbox.ColorWhite)
		}
	}
}

func (t Termbox) GetInput() Input {
	var curEvent termbox.Event

	select {
	case e, ok := <-t.events:
		curEvent = e
		if !ok {
			return Input{}
		}
	}

	i := Input{}

	// If we've gotten to this point, we have an event that's ready to process.
	switch curEvent.Type {
	case termbox.EventKey:
		switch curEvent.Key {

		case termbox.KeyEsc:
			i.KeyEsc = true

		case termbox.KeyCtrlC:
			i.KeyEsc = true
		}
	}

	return i
}

func (t Termbox) Shutdown() {
	termbox.Close()
}
