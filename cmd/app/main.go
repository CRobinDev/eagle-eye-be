package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/CRobinDev/karsa/config/server"
	"github.com/sirupsen/logrus"
)

func main() {
	ctx := context.Background()
	app := server.Init()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func () {
		sig := <- signalChan
		logrus.WithFields(logrus.Fields{
			"signal": sig.String(),
			"pid" : os.Getpid(),
		}).Info("Received shutdown signal, shutting down gracefully")

		app.Shutdown(ctx)
		os.Exit(0)
	}()
	
	app.Start()

}
