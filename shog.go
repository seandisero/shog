package shog

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Shoggoth struct {
	oldState *term.State
	newState *term.State
	Canvas   Screen
}

func (s *Shoggoth) Listen(input chan byte, ctx context.Context) {
	s.Canvas.Windows[0].inputCh = input
	s.Canvas.Windows[0].ctx = ctx
	go s.Canvas.Windows[0].handleInput()
	for range s.Canvas.redraw {
		s.Canvas.Draw()
	}
}

func SpawnShoggoth() (*Shoggoth, error) {
	w, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}

	// TODO: I should pull some functionality out of here maybe lean into the
	// functional options patern a bit which means creating a spawnShoggoth
	// function that takes in an Options struct so I can have
	// SpawnShoggothWithWindows(number_of_windows)

	screen := NewScreen(w, h)
	ctx1 := context.Background()
	ctx2 := context.Background()
	wind := NewWindowWithCoords(w/2-1, h-2, NewUV(0, 0), ctx1)
	wind2 := NewWindowWithCoords(w/2-1, h-2, NewUV(w/2, 0), ctx2)
	screen.AddWindow(wind)
	screen.AddWindow(wind2)

	shoggoth := &Shoggoth{
		oldState: nil,
		newState: nil,
		Canvas:   screen,
	}

	err = shoggoth.Delve()
	if err != nil {
		return nil, err
	}

	return shoggoth, nil
}

func (s *Shoggoth) Delve() error {
	var err error
	// NOTE: I wonder if there is a better pattern for this. I don't like that
	// my cleanup is a side effect of End(), but I can't think of anything
	// else at the moment
	s.oldState, err = term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	os.Stdout.Write([]byte("\033[?25l"))   // hide cursor
	os.Stdout.Write([]byte("\033[?1049h")) // enable alternative screen buffer
	fmt.Print("\033[48;2;24;26;26m")       // #1a1b26

	s.newState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	s.Canvas.Draw()
	return nil
}

func (s *Shoggoth) End() {
	reset()
	term.Restore(int(os.Stdin.Fd()), s.oldState)
	for i := range s.Canvas.Windows {
		s.Canvas.Windows[i].ctx.Done()
	}
}

func reset() {
	fmt.Printf("\033[?25h")   // show cursor
	fmt.Printf("\033[?1049l") // back to origional screen buffer
}
