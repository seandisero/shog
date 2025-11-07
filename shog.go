package shog

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"golang.org/x/term"
)

type Shoggoth struct {
	oldState *term.State
	newState *term.State
	Screen   Screen
}

type ShoggothConfig struct {
	// TODO: this should contain things like:
	//		- a new Theme struct to hold colors
	//		- the number and positions of pannels
}

func (s *Shoggoth) Listen(ctx context.Context) {
outer:
	for {
		select {
		case <-s.Screen.redraw:
			s.Screen.Draw()
		case <-ctx.Done():
			break outer
		}
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
	// SpawnShoggothWithPannels(number_of_pannels)
	// TODO: I should create a config for this

	screen := NewScreen(w, h)

	shoggoth := &Shoggoth{
		oldState: nil,
		newState: nil,
		Screen:   screen,
	}

	err = shoggoth.delve()
	if err != nil {
		return nil, fmt.Errorf("error delving into chaos%w", err)
	}
	return shoggoth, nil
}

func (s *Shoggoth) screenResized() {
	for {
		w, h, err := term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			slog.Error("couldn't get size of terminal")
			return
		}
		if w != s.Screen.ScreenSize.X || h != s.Screen.ScreenSize.Y {
			// TODO: I could put resizing logic here or call screen rezise
			//		function
		}
	}
}

func (s *Shoggoth) NameShoggoth(name string) {
	s.Screen.headerText = fmt.Sprintf("~~ %s ~~", name)
}

func (s *Shoggoth) delve() error {
	// TODO: I should at least find a way to reset the terminal to the old state
	//		on panic.
	var err error
	s.oldState, err = term.GetState(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	os.Stdout.Write([]byte("\033[?25l"))           // hide cursor
	os.Stdout.Write([]byte("\033[?1049h"))         // enable alternative screen buffer
	os.Stdout.Write([]byte("\033[48;2;24;26;26m")) // #1a1b26

	s.newState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	s.Screen.Draw()
	return nil
}

func (s *Shoggoth) AddPannel(p *Pannel) {
	ctx := context.Background()
	p.ctx = ctx
	s.Screen.AddPannel(p)
}

func (s *Shoggoth) End() {
	term.Restore(int(os.Stdin.Fd()), s.oldState)
	reset()
	for i := range s.Screen.Pannels {
		s.Screen.Pannels[i].ctx.Done()
	}
}

func reset() {
	os.Stdout.Write([]byte("\033[?25h"))
	os.Stdout.Write([]byte("\033[?1049l"))
	os.Stdout.Sync()
}
