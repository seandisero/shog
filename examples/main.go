package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/seandisero/shog"
)

func handleInput(in chan byte) {
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
		in <- buf[0]
	}
}

func main() {
	// NOTE: I wonder if this can be simplified a bit more, how do I want
	// the user to interact with the Shoggoth? all I know is I want it
	// to be very simple interaction and have it 'just work'
	shoggoth, err := shog.SpawnShoggoth()
	if err != nil {
		slog.Error("could not spawn shoggoth", "error", err)
	}
	defer shoggoth.End()

	// NOTE: how would I do this if creating a chat application? I would want
	// my input channel to be listening for text responces from the server.
	// so what would that look like here, how should I interact and set
	// channels
	inputChan := make(chan byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go shoggoth.Listen(inputChan, ctx)
	handleInput(inputChan)
}
