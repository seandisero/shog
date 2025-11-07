package shog

import (
	"context"
	"fmt"
)

type Canvas []rune

type Pannel struct {
	Origin     UV // determines where the pannel is placed within the screen
	Size       UV // size of the pannel including border
	CanvasSize UV // size of the writable area inside the border

	// TODO: transition from Input to Canvas draw
	input    []byte
	canvas   Canvas
	redrawCh chan struct{}
	inputCh  chan byte
	ctx      context.Context

	// style
	Border PannelBorder // border symbols

	// images
	Images []*Image
}

type PannelOption func(p *Pannel)

func (p *Pannel) NewCanvas() {
	p.CanvasSize = NewUV(p.Size.X-2, p.Size.Y-2)
	size := p.CanvasSize.X * p.CanvasSize.Y
	p.canvas = make(Canvas, size)
	for i := 0; i < len(p.canvas); i++ {
		p.canvas[i] = rune(NonBreakingSpace)
	}
}

func NewPannel(options ...PannelOption) *Pannel {
	pan := &Pannel{
		Border: NewPannelBorder(),
		Images: make([]*Image, 0),
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

func (p *Pannel) AddImage(image *Image) error {
	p.Images = append(p.Images, image)
	return nil
}

func (p *Pannel) DrawImages() {
	for i := range p.Images {
		p.DrawImage(p.Images[i])
	}
}

func (p *Pannel) DrawImage(image *Image) error {
	if p.CanvasSize.X < image.origin.X+image.size.X {
		return fmt.Errorf("could not draw image")
	}
	j := image.origin.X + (p.CanvasSize.X * (image.origin.Y))
	for i := 0; i < len(image.Data); i++ {
		if j%p.CanvasSize.X == image.origin.X { // should move past border
			j++
		}
		if j%p.CanvasSize.X == image.origin.X+image.size.X+1 {
			// should wrap to other side of canvas
			j += p.CanvasSize.X - (image.size.X)
		}
		p.canvas[j] = image.Data[i]
		j++
	}
	return nil
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

func (p *Pannel) doInput(in byte) {
	// TODO: add delete and no-break space to Key.go
	// TODO: first I want to implement drawing images, so I can leave the inputs
	//		for now
	switch in {
	case byte(127):
		if len(p.input) == 0 {
			break
		}
		i := 1
		for p.input[len(p.input)-i] == byte(160) {
			i++
		}
		p.input = p.input[:len(p.input)-i]
	case byte(CarriageReturn):
		spaceNum := p.CanvasSize.X - (len(p.input) % (p.CanvasSize.X)) + 1
		spaces := make([]byte, spaceNum-1)
		for i := range spaces {
			spaces[i] = byte(160)
		}
		p.input = append(p.input, spaces...)
	default:
		p.input = append(p.input, in)
	}
	p.redrawCh <- struct{}{}
}
