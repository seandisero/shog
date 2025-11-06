package shog

import (
	"fmt"
	"os"
)

type Screen struct {
	Pixels     []rune
	ScreenSize UV
	Windows    []Window
	redraw     chan struct{}
}

// TODO: if I'm moving to the functional options pattern for creating a
// Shoggoth, how would that look here? maybe that doesn't need to be
// answered yet
func NewScreen(w, h int) Screen {
	return Screen{
		Pixels:     make([]rune, w*h),
		ScreenSize: NewUV(w, h),
		Windows:    make([]Window, 0, 3),
		redraw:     make(chan struct{}),
	}
}

func (s *Screen) AddWindow(w Window) {
	// NOTE: maybe I don't want to support the adding of windows, I could
	// just create a new Shoggoth instead.
	w.redrawCh = s.redraw
	s.Windows = append(s.Windows, w)
}

func (s *Screen) Draw() {
	fmt.Print("\033[H") // move cursor back to top left
	s.initScreen()
	s.drawHeader()
	for i := range s.Windows {
		s.drawBorder(s.Windows[i])
		s.drawInput(s.Windows[i])
	}
	os.Stdout.Write([]byte(string(s.Pixels)))
}

func (s *Screen) initScreen() {
	for i := range s.Pixels {
		s.Pixels[i] = ' '
	}
}

func (s *Screen) drawHeader() {
	message := " -> new shoggoth program <-"
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

func (s *Screen) drawInput(w Window) {
	j := w.Origin.X + (s.ScreenSize.X * (w.Origin.Y + 1 + 1)) + 1
	for i := 0; i < len(w.Input); i++ {
		if j%s.ScreenSize.X == w.Origin.X {
			j++
		}
		if j%s.ScreenSize.X == w.Origin.X+w.W {
			j += s.ScreenSize.X - (w.W - 1)
		}
		s.Pixels[j] = rune(w.Input[i])
		j++
	}
}

func (s *Screen) drawBorder(w Window) {
	i := w.Origin.X + (s.ScreenSize.X * (w.Origin.Y + 1))
	// top margin
	// TODO: these box symbols should be moved to the key.go file,
	// or maybe into a new symbols.go since they are not really keypresses
	s.Pixels[i] = '\u250c' // ┌
	i++
	offset := i + w.W
	for i < offset-1 {
		s.Pixels[i] = rune('\u2500') // ─
		i++
	}
	s.Pixels[i] = '\u2510' // ┐
	i++

	// sides
	i += s.ScreenSize.X - w.W - 1
	for j := 0; j < w.H-1; j++ {
		offset = i + w.W
		s.Pixels[i] = '\u2502' // │
		for i < offset {
			i++
		}
		s.Pixels[i] = '\u2502' // │
		i += s.ScreenSize.X - w.W
	}

	// bottom
	s.Pixels[i] = '\u2514' // └
	i++
	offset = i + w.W
	for i < offset-1 {
		s.Pixels[i] = '\u2500' // ─
		i++
	}
	s.Pixels[i] = '\u2518' // ┘
	i++
}
