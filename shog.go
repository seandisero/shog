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

func (s *Shoggoth) Listen() {
	err := s.Delve()
	if err != nil {
		slog.Error("could not setup terminal", "error", err)
		s.End()
		return
	}
	for range s.Screen.redraw {
		s.Screen.Draw()
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
	// ctx1 := context.Background()
	// ctx2 := context.Background()
	// wind := NewPannelWithCoords(w/2-1, h-2, NewUV(0, 0), ctx1)
	// wind2 := NewPannelWithCoords(w/2-1, h-2, NewUV(w/2, 0), ctx2)
	// screen.AddPannel(wind)
	// screen.AddPannel(wind2)

	shoggoth := &Shoggoth{
		oldState: nil,
		newState: nil,
		Screen:   screen,
	}

	// err = shoggoth.Delve()
	// if err != nil {
	// 	return nil, err
	// }

	return shoggoth, nil
}

func (s *Shoggoth) NameShoggoth(name string) {
	s.Screen.headerText = fmt.Sprintf("~~ %s ~~", name)
}

func (s *Shoggoth) Delve() error {
	// NOTE: at first I thought changing the terminal state would be better when
	//		spawning a new shoggoth, but I think doing that on delve would be a
	//		better way to go. don't change the state until the program runs.
	// NOTE: I wonder if there is a better pattern for this. I don't like that
	// my cleanup is a side effect of End(), but I can't think of anything
	// else at the moment
	// TODO: I should at least find a way to reset the terminal to the old state
	//		on panic.
	var err error
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

	s.Screen.Draw()
	return nil
}

func (s *Shoggoth) AddPannel(p *Pannel) {
	ctx := context.Background()
	p.ctx = ctx
	s.Screen.AddPannel(p)
}

func (s *Shoggoth) End() {
	reset()
	if s.oldState != nil {
		term.Restore(int(os.Stdin.Fd()), s.oldState)
	}
	for i := range s.Screen.Pannels {
		s.Screen.Pannels[i].ctx.Done()
	}
}

func reset() {
	fmt.Printf("\033[?25h")   // show cursor
	fmt.Printf("\033[?1049l") // back to origional screen buffer
}
