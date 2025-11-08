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
	FixedSize  bool
	MaxSize    UV
	Offscreen  bool
	ZOrder     int

	// TODO: transition from Input to Canvas draw
	input    []byte
	canvas   Canvas
	redrawCh chan struct{}
	inputCh  chan byte
	ctx      context.Context
	Dirty    bool

	// style
	Border PannelBorder // border symbols

	// images
	Images []*Image
}

type PannelOption func(p *Pannel)

func (p *Pannel) NewCanvas() {
	p.CanvasSize = NewUV(p.Size.X-2, p.Size.Y-2)
	size := p.CanvasSize.X * p.CanvasSize.Y
	if size < 0 {
		size = 0
	}
	p.canvas = make(Canvas, size)
	for i := 0; i < len(p.canvas); i++ {
		p.canvas[i] = rune(NonBreakingSpace)
	}
}

func NewPannel(options ...PannelOption) *Pannel {
	pan := &Pannel{
		Border:    NewPannelBorder(),
		Images:    make([]*Image, 0),
		Dirty:     true,
		FixedSize: false,
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
		p.MaxSize = NewUV(u, v)
		p.FixedSize = true
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
	// p.adjustInput()
}

func (p *Pannel) AdjustSize(uv UV) {
	if p.Origin.X < uv.X || p.Origin.Y < uv.Y {
		p.Offscreen = true
	}
	p.Offscreen = false
	p.Size.X = min(p.MaxSize.X, uv.X-p.Origin.X)
	p.Size.Y = min(p.MaxSize.Y, uv.Y-p.Origin.Y-1)
	p.NewCanvas()
	// p.adjustInput()
}

func (p *Pannel) AddImage(image *Image) error {
	p.Images = append(p.Images, image)
	p.Dirty = true
	return nil
}

func (p *Pannel) DrawImages() {
	for i := range p.Images {
		p.DrawImage(p.Images[i])
	}
}

func (p *Pannel) DrawImage(image *Image) error {
	if image.origin.Y >= p.CanvasSize.Y || image.origin.X >= p.CanvasSize.X {
		return fmt.Errorf("image outside the bounds of pannel")
	}
	j := image.origin.X + (p.CanvasSize.X * (image.origin.Y))
	if j > len(p.canvas) {
		return nil
	}
	PX := p.CanvasSize.X - image.origin.X
	PPX := image.size.X - PX
	for i := 0; i < len(image.Data); i++ {
		jmod := j % p.CanvasSize.X // u index in pannel
		if j > len(p.canvas) {
			break
		}
		if jmod == p.CanvasSize.X-1 {
			p.canvas[j] = image.Data[i]
			i += PPX
			j += p.CanvasSize.X - (image.size.X - PPX - 1)
			continue
		}
		// }
		if jmod == image.origin.X+image.size.X-1 {
			p.canvas[j] = image.Data[i]
			// should wrap to other side of canvas
			j += p.CanvasSize.X - (image.size.X)
			if j > p.CanvasSize.Square {
				break
			}
			j++
			continue
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
	case byte(NonBreakingSpace):
		{
			return
		}
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
	p.Dirty = true
	p.redrawCh <- struct{}{}
}

// func (p *Pannel) adjustInput() {
// 	newInput := string(p.input)
// 	lines := strings.Split(newInput, string(NonBreakingSpace))
// 	nonEmptyLines := make([]string, 0)
// 	for _, line := range lines {
// 		if line != "" {
// 			nonEmptyLines = append(nonEmptyLines, line)
// 		}
// 	}
//
// 	joined := strings.Join(nonEmptyLines, "\n")
//
// 	i := 0
// 	for _, v := range joined {
// 		switch v {
// 		case '\n':
// 			spaceNum := p.CanvasSize.X - (len(p.input) % (p.CanvasSize.X)) + 1
// 			current := i
// 			for i < current+spaceNum {
// 				p.input[i] = byte(NonBreakingSpace)
// 				i++
// 			}
// 		default:
// 			p.input[i] = byte(v)
// 		}
// 	}
// }
