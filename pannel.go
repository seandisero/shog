package shog

import (
	"context"
)

type Canvas []rune

type Pannel struct {
	// TODO: should width and hight be in a UV struct?
	// dimensions
	Origin     UV // determines where the pannel is placed within the screen
	Size       UV // size of the pannel including border
	CanvasSize UV // size of the writable area inside the border

	// TODO: transition from Input to Canvas draw
	Input    []byte
	canvas   Canvas // TODO: needs implementation
	redrawCh chan struct{}
	inputCh  chan byte
	ctx      context.Context

	// style
	Border PannelBorder // border symbols
}

type PannelOption func(p *Pannel)

func (p *Pannel) NewCanvas() {
	p.CanvasSize = NewUV(p.Size.X-2, p.Size.Y-2)
	size := p.CanvasSize.X * p.CanvasSize.Y
	p.canvas = make(Canvas, size)
	for i := 0; i < len(p.canvas); i++ {
		p.canvas[i] = rune(160) // 160 is the ascii non breaking space
	}
}

func NewPannel(options ...PannelOption) *Pannel {
	pan := &Pannel{
		Border: NewPannelBorder(),
	}
	for _, option := range options {
		option(pan)
	}
	return pan
}

func WithBorderOptions(borderOptions ...PannelBorderOption) PannelOption {
	return func(p *Pannel) {
		for _, option := range borderOptions {
			option(&p.Border)
		}
	}
}

func WithSize(u, v int) PannelOption {
	return func(p *Pannel) {
		p.Size = NewUV(u, v)
		p.NewCanvas()
	}
}

func WithOrigin(u, v int) PannelOption {
	return func(p *Pannel) {
		p.Origin = NewUV(u, v)
	}
}

func (p *Pannel) SetSize(uv UV) {
	p.Size = uv
	p.NewCanvas()
}

func (w *Pannel) HandleInput(input chan byte, ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		case in := <-input:
			w.doInput(in)
		}
	}
}

func (w *Pannel) doInput(in byte) {
	// TODO: add delete and no-break space to Key.go
	// TODO: first I want to implement drawing images, so I can leave the inputs
	//		for now
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
		spaceNum := w.Size.X - (len(w.Input) % (w.Size.X - 1))
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
