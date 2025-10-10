package shog

import (
	"context"
)

type Window struct {
	W        int
	H        int
	Origin   UV
	Input    []byte
	redrawCh chan struct{}
	inputCh  chan byte
	ctx      context.Context
}

func NewWindowWithCoords(w, h int, uv UV, ctx context.Context) Window {
	win := Window{
		W:       w,
		H:       h,
		Origin:  uv,
		Input:   make([]byte, 0, 1024),
		inputCh: make(chan byte),
		ctx:     ctx,
	}
	go win.handleInput()
	return win
}

func (w *Window) handleInput() {
	for {
		select {
		case <-w.ctx.Done():
			return
		case input := <-w.inputCh:
			w.doInput(input)
		}
	}
}

func (w *Window) doInput(in byte) {
	switch in {
	case byte(127):
		if len(w.Input) == 0 {
			break
		}
		i := 1
		for w.Input[len(w.Input)-i] == byte(160) {
			i++
		}
		w.Input = w.Input[:len(w.Input)-i]
	case byte(CarriageReturn):
		spaceNum := w.W - (len(w.Input) % (w.W - 1))
		spaces := make([]byte, spaceNum-1)
		for i := range spaces {
			spaces[i] = byte(160)
		}
		w.Input = append(w.Input, spaces...)
	default:
		w.Input = append(w.Input, in)
	}
	w.redrawCh <- struct{}{}
}
