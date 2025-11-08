package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/seandisero/shog"
)

func handleInput(inputChan chan byte) {
	for {
		buf := make([]byte, 1)
		n, err := os.Stdin.Read(buf)
		if err != nil {
			slog.Error("error reding input", "error", err)
			return
		}
		if n > 0 {
			char := buf[0]
			if char == 3 { // this would be ctrl + c
				return
			}
		}
		inputChan <- buf[0]
	}
}

func main() {
	shoggoth, err := shog.SpawnShoggoth()
	if err != nil {
		slog.Error("could not spawn shoggoth", "error", err)
	}
	defer shoggoth.End()
	shoggoth.NameShoggoth("example app")

	background := shog.NewPannel(shog.WithSize(32, 2))
	wind := shog.NewPannel(
		shog.WithSize(16, 16),
		shog.WithOrigin(16, 2),
	)
	wind2 := shog.NewPannel(
		shog.WithSize(16, 8),
		shog.WithOrigin(0, 2),
	)
	wind3 := shog.NewPannel(
		shog.WithSize(16, 8),
		shog.WithOrigin(0, 10),
	)
	shoggoth.AddPannel("background", background)
	shoggoth.AddPannel("inputWind", wind)
	shoggoth.AddPannel("win2", wind2)
	shoggoth.AddPannel("win3", wind3)

	// background.AddImage(&shog.TEST_IMAGE)
	// background.AddImage(&shog.TEST_IMAGE2)
	// shog.TEST_IMAGE3.SetOrigin(shog.NewUV(1, 1))
	// background.AddImage(shog.TEST_IMAGE3)

	inputChan := make(chan byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wind.HandleInput(inputChan, ctx)
	go shoggoth.Listen(ctx)
	handleInput(inputChan)
}
