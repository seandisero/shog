package shog

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/term"
)

type Screen struct {
	Pixels     []rune
	ScreenSize UV
	Pannels    map[string]*Pannel
	pnls       []*Pannel
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
		Pannels:    make(map[string]*Pannel),
		pnls:       make([]*Pannel, 0),
		redraw:     make(chan struct{}),
	}
}

func (s *Screen) AddPannel(name string, p *Pannel) {
	if len(s.Pannels) == 0 {
		if p.Size.Zero() {
			p.SetSize(NewUV(s.ScreenSize.X-1, s.ScreenSize.Y-2))
		}
	}
	p.redrawCh = s.redraw
	s.Pannels[name] = p
	s.pnls = append(s.pnls, p)
	s.Draw()
}

func (s *Screen) MarkAllPannelsDirty() {
	for key := range s.Pannels {
		s.Pannels[key].Dirty = true
	}
}

func resetCursor() {
	fmt.Print("\033[H")
}

func (s *Screen) checkSize() {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		slog.Error("couldn't get size of terminal")
		return
	}
	if w != s.ScreenSize.X || h != s.ScreenSize.Y {
		s.changeScreenSize(w, h)
	}
}

func (s *Screen) changeScreenSize(w, h int) {
	s.Pixels = make([]rune, w*h)
	s.ScreenSize = NewUV(w, h)
	for key := range s.Pannels {
		if !s.Pannels[key].FixedSize {
			s.Pannels[key].SetSize(NewUV(w, h-1))
			continue
		}
		s.Pannels[key].AdjustSize(NewUV(w, h))
	}
}

func (s *Screen) OutOfBounds(key string) bool {
	pannel, ok := s.Pannels[key]
	if !ok {
		return false
	}
	return pannel.Size.X+pannel.Origin.X > s.ScreenSize.X ||
		pannel.Size.Y+pannel.Origin.Y > s.ScreenSize.Y
}

func (s *Screen) Draw() {
	s.checkSize()
	resetCursor()
	s.initScreen()
	s.drawHeader()
	for i := range s.pnls {
		if !s.pnls[i].Dirty {
			continue
		}
		if s.pnls[i].Offscreen {
			continue
		}
		s.drawBorder(s.pnls[i])
		s.drawCanvas(s.pnls[i])
	}
	os.Stdout.Write([]byte(string(s.Pixels)))
}

func (s *Screen) initScreen() {
	for i := range s.Pixels {
		s.Pixels[i] = rune(160)
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
	for i := 0; i < len(w.input); i++ {
		if j%s.ScreenSize.X == w.Origin.X {
			j++
		}
		if j%s.ScreenSize.X == w.Origin.X+w.Size.X {
			j += s.ScreenSize.X - (w.Size.X - 1)
		}
		s.Pixels[j] = rune(w.input[i])
		j++
	}
}

func (s *Screen) drawCanvas(p *Pannel) {
	if p.Offscreen {
		return
	}
	for i := 0; i < len(p.canvas); i++ {
		p.canvas[i] = rune(160)
	}
	if len(p.input) > 0 {
		for i := range p.input {
			if i >= len(p.canvas) {
				break
			}
			p.canvas[i] = rune(p.input[i])
		}
	}
	p.DrawImages()
	j := p.Origin.X + (s.ScreenSize.X * (p.Origin.Y + 1 + 1)) + 1
	for i := 0; i < len(p.canvas); i++ {
		if j%s.ScreenSize.X == p.Origin.X { // should move past border
			j++
		}
		if j%s.ScreenSize.X == p.Origin.X+p.CanvasSize.X+1 {
			// should wrap to other side of canvas
			j += s.ScreenSize.X - (p.Size.X - 2)
		}
		if i > p.CanvasSize.Square {
			break
		}
		s.setPixel(j, p.canvas[i])
		// s.Pixels[j] = p.canvas[i]
		j++
	}
}

func (s *Screen) drawBorder(p *Pannel) {
	if p.Offscreen {
		return
	}
	i := p.Origin.X + (s.ScreenSize.X * (p.Origin.Y + 1))
	// top margin
	if i > s.ScreenSize.Square {
		return
	}
	s.setPixel(i, rune(p.Border.TopLeft))
	i++
	offset := i + p.Size.X - 1
	for i < offset-1 {
		s.setPixel(i, rune(p.Border.Horizontal))
		// s.Pixels[i] = rune(p.Border.Horizontal)
		i++
	}
	s.setPixel(i, rune(p.Border.TopRight))
	// s.Pixels[i] = rune(p.Border.TopRight)
	i++

	// sides
	i += s.ScreenSize.X - p.Size.X
	for j := 0; j < p.Size.Y-2; j++ {
		if p.Origin.Y+j > s.ScreenSize.Y {
			return
		}
		offset = i + p.Size.X
		s.setPixel(i, rune(p.Border.Virtical))
		for i < offset-1 {
			i++
		}
		s.setPixel(i, rune(p.Border.Virtical))
		// s.Pixels[i] = rune(p.Border.Virtical)
		i += (s.ScreenSize.X - p.Size.X) + 1
	}

	// bottom
	if i > s.ScreenSize.Square {
		return
	}
	s.setPixel(i, rune(p.Border.BottomLeft))
	// s.Pixels[i] = rune(p.Border.BottomLeft)
	i++
	offset = i + p.Size.X
	for i < offset-2 {
		s.setPixel(i, rune(p.Border.Horizontal))
		// s.Pixels[i] = rune(p.Border.Horizontal)
		i++
	}
	s.setPixel(i, rune(p.Border.BottomRight))
	// s.Pixels[i] = rune(p.Border.BottomRight)
	i++
}

func (s *Screen) setPixel(index int, r rune) {
	if index > s.ScreenSize.Square-1 {
		return
	}
	s.Pixels[index] = r
}
