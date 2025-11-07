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

	background := shog.NewPannel()
	wind := shog.NewPannel(
		shog.WithSize(16, 16),
		shog.WithOrigin(4, 4),
		shog.WithBorderOptions(shog.WithCustomPannelBorder(
			shog.Symbol(160),
			shog.Symbol(160),
			shog.Symbol(160),
			shog.Symbol(160),
			shog.Symbol(160),
			shog.Symbol(160),
		)),
	)
	shoggoth.AddPannel(background)
	shoggoth.AddPannel(wind)

	background.AddImage(&shog.TEST_IMAGE)
	background.AddImage(&shog.TEST_IMAGE2)
	shog.TEST_IMAGE3.SetOrigin(shog.NewUV(64, 16))
	background.AddImage(shog.TEST_IMAGE3)

	inputChan := make(chan byte)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go wind.HandleInput(inputChan, ctx)
	go shoggoth.Listen(ctx)
	handleInput(inputChan)
}
