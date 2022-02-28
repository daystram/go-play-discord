package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/daystram/go-play-discord/internal/command"
	"github.com/daystram/go-play-discord/internal/config"
	"github.com/daystram/go-play-discord/internal/server"
)

func init() {
	log.SetOutput(os.Stderr)
}

func main() {
	err := Main()
	if err != nil {
		log.Println("init:", err)
		os.Exit(exitErr)
	}

	os.Exit(exitOk)
}

func Main() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	srv, err := server.Start(cfg)
	if err != nil {
		return err
	}

	defer func() {
		err := srv.Stop()
		if err != nil {
			log.Println("exit:", err)
			os.Exit(exitErr)
		}
	}()

	log.Println("init: server started")

	err = command.RegisterAll(srv)
	if err != nil {
		return err
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-stop

	return nil
}
