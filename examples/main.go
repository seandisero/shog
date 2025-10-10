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
	shoggoth, err := shog.SpawnShoggoth()
	if err != nil {
		slog.Error("could not spawn shoggoth", "error", err)
	}
	defer shoggoth.End()

	inputChan := make(chan byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go shoggoth.Listen(inputChan, ctx)
	handleInput(inputChan)
}
