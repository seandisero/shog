package shog

import (
	"fmt"
	"os"
)

type Screen struct {
	Pixels     []rune
	ScreenSize UV
	Pannels    []*Pannel
	headerText string
	redraw     chan struct{}
}

// TODO: if I'm moving to the functional options pattern for creating a
// Shoggoth, how would that look here? maybe that doesn't need to be
// answered yet
func NewScreen(w, h int) Screen {
	return Screen{
		Pixels:     make([]rune, w*h),
		ScreenSize: NewUV(w, h),
		Pannels:    make([]*Pannel, 0, 3),
		redraw:     make(chan struct{}),
	}
}

func (s *Screen) AddPannel(p *Pannel) {
	if len(s.Pannels) == 0 {
		if p.Size.Zero() {
			p.SetSize(NewUV(s.ScreenSize.X-1, s.ScreenSize.Y-2))
		}
	}
	p.redrawCh = s.redraw
	s.Pannels = append(s.Pannels, p)
	s.Draw()
}

func resetCursor() {
	fmt.Print("\033[H")
}

func (s *Screen) Draw() {
	resetCursor()
	s.initScreen()
	s.drawHeader()
	for i := range s.Pannels {
		s.drawBorder(s.Pannels[i])
		s.drawCanvas(s.Pannels[i])
	}
	os.Stdout.Write([]byte(string(s.Pixels)))
}

func (s *Screen) initScreen() {
	for i := range s.Pixels {
		s.Pixels[i] = ' '
	}
}

func (s *Screen) drawHeader() {
	message := s.headerText
	space := s.ScreenSize.X - len(message)
	lhSpace := space / 2
	rhSpace := space - lhSpace
	i := 0
	for i < lhSpace {
		s.Pixels[i] = ' '
		i++
	}
	for _, r := range message {
		s.Pixels[i] = r
		i++
	}
	offset := i
	for i < offset+rhSpace {
		s.Pixels[i] = ' '
		i++
	}
}

func (s *Screen) drawInput(w *Pannel) {
	j := w.Origin.X + (s.ScreenSize.X * (w.Origin.Y + 1 + 1)) + 1
	for i := 0; i < len(w.Input); i++ {
		if j%s.ScreenSize.X == w.Origin.X {
			j++
		}
		if j%s.ScreenSize.X == w.Origin.X+w.Size.X {
			j += s.ScreenSize.X - (w.Size.X - 1)
		}
		s.Pixels[j] = rune(w.Input[i])
		j++
	}
}

func (s *Screen) drawCanvas(p *Pannel) {
	// TODO: I would like to have pannels be full instead of appending and
	//		resizing, that way I can insert images or place text starting
	//		points wherever I want.
	// The new way of working would be to draw a buffer to the canvas instead
	//		of draw input
	j := p.Origin.X + (s.ScreenSize.X * (p.Origin.Y + 1 + 1)) + 1
	for i := 0; i < len(p.canvas); i++ {
		if j%s.ScreenSize.X == p.Origin.X { // should move past border
			j++
		}
		if j%s.ScreenSize.X == p.Origin.X+p.CanvasSize.X+1 {
			// should wrap to other side of canvas
			j += s.ScreenSize.X - (p.Size.X - 2)
		}
		s.Pixels[j] = p.canvas[i]
		j++
	}

}

func (s *Screen) drawBorder(p *Pannel) {
	i := p.Origin.X + (s.ScreenSize.X * (p.Origin.Y + 1))
	// top margin
	// TODO: these box symbols should be moved into a new symbols.go
	//		since they are not really keypresses
	s.Pixels[i] = rune(p.Border.TopLeft)
	i++
	offset := i + p.Size.X - 1
	for i < offset-1 {
		s.Pixels[i] = rune(p.Border.Horizontal)
		i++
	}
	s.Pixels[i] = rune(p.Border.TopRight)
	i++

	// sides
	i += s.ScreenSize.X - p.Size.X
	for j := 0; j < p.Size.Y-2; j++ {
		offset = i + p.Size.X
		s.Pixels[i] = rune(p.Border.Virtical)
		for i < offset-1 {
			i++
		}
		s.Pixels[i] = rune(p.Border.Virtical)
		i += (s.ScreenSize.X - p.Size.X) + 1
	}

	// bottom
	s.Pixels[i] = rune(p.Border.BottomLeft)
	i++
	offset = i + p.Size.X
	for i < offset-2 {
		s.Pixels[i] = rune(p.Border.Horizontal)
		i++
	}
	s.Pixels[i] = rune(p.Border.BottomRight)
	i++
}
