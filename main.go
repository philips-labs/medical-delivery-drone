package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/philips-labs/medical-delivery-drone/drone"
	"github.com/philips-labs/medical-delivery-drone/video"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go gracefulShutdown(cancel)

	converter, err := video.NewConverter()
	if err != nil {
		return
	}
	defer converter.Close()

	videoChan, err := drone.Connect(ctx, converter)
	err = video.Display(ctx, videoChan)

	log.Println("Shutdown, completed")
}

func gracefulShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	sig := <-quit

	log.Println("Shutting down, reason:", sig.String())
	cancel()
}
