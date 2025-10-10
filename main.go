package main

import (
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/term"
)

type UV struct {
	X int
	Y int
}

func NewUV(x, y int) UV {
	return UV{
		X: x,
		Y: y,
	}
}

type Window struct {
	W      int
	H      int
	Origin UV
	Input  []byte
}

func (w *Window) SetInput(input byte) {
	w.Input = append(w.Input, input)
}

type Screen struct {
	Pixels     []rune
	ScreenSize UV
	Windows    []Window
}

func (s *Screen) AddWindow(w Window) {
	s.Windows = append(s.Windows, w)
}

func (s *Screen) Draw() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[H")
	s.initScreen()
	for _, w := range s.Windows {
		s.drawBorder(w)
		s.drawInput(w)
	}
	fmt.Printf("%s", string(s.Pixels))
}

func (s *Screen) initScreen() {
	for i := range s.Pixels {
		s.Pixels[i] = ' '
	}
}

func (s *Screen) drawInput(w Window) {
	j := w.Origin.X + (s.ScreenSize.X * (w.Origin.Y + 1)) + 1
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
	// screen size is 56, 18
	i := w.Origin.X + (s.ScreenSize.X * w.Origin.Y)
	// top margin
	s.Pixels[i] = '|'
	i++
	offset := i + w.W
	for i < offset-1 {
		s.Pixels[i] = rune('\u203E')
		i++
	}
	s.Pixels[i] = '|'
	i++

	// sides
	i += s.ScreenSize.X - w.W - 1
	for j := 0; j < w.H-1; j++ {
		offset = i + w.W
		s.Pixels[i] = '|'
		for i < offset {
			i++
		}
		s.Pixels[i] = '|'
		i += s.ScreenSize.X - w.W
	}

	// bottom
	s.Pixels[i] = '|'
	i++
	offset = i + w.W
	for i < offset-1 {
		s.Pixels[i] = '_'
		i++
	}
	s.Pixels[i] = '|'
	i++
}

func NewWindow(w, h int) Window {
	return Window{
		W:      w,
		H:      h,
		Origin: NewUV(0, 0),
		Input:  make([]byte, 0, 1024),
	}
}

func NewWindowWithCoords(w, h int, uv UV) Window {
	return Window{
		W:      w,
		H:      h,
		Origin: uv,
		Input:  make([]byte, 0, 1024),
	}
}

func NewScreen(w, h int) Screen {
	return Screen{
		Pixels:     make([]rune, w*h),
		ScreenSize: NewUV(w, h),
		Windows:    make([]Window, 0, 3),
	}
}

func reset() {
	fmt.Printf("\033[?25h")
	fmt.Printf("\033[?1049l")
}

func main() {
	oldState, err := term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		slog.Error("could not get term state", "error", err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	fmt.Printf("\033[?25l")
	fmt.Printf("\033[?1049h")
	defer reset()

	term.MakeRaw(int(os.Stdin.Fd()))

	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		slog.Error("problem getting size of window", "error", err)
	}

	screen := NewScreen(w, h)
	wind := NewWindowWithCoords(w/2, h-1, NewUV(0, 0))
	wind2 := NewWindowWithCoords(w/2-1, h-1, NewUV(w/2, 0))
	screen.AddWindow(wind)
	screen.AddWindow(wind2)
	screen.Draw()

	for {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			slog.Error("error reding input", "error", err)
			return
		}
		if n > 0 {
			char := buf[0]
			if char == 3 {
				return
			}
		}

		// cant iterate over a slice of structs, if I were to use _, w :=
		// that w would be a coppy of what is in the slice.
		for i, _ := range screen.Windows {
			screen.Windows[i].SetInput(buf[0])
		}
		screen.Draw()
	}
}
