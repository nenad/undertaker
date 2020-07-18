package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nenad/undertaker/internal/config"
	"github.com/nenad/undertaker/internal/database"
	"github.com/nenad/undertaker/internal/loader"
	"github.com/nenad/undertaker/internal/server"
)

func main() {
	cfg, err := config.ParseArgs()
	if err != nil {
		fmt.Printf("could not load configuration: %s", err)
		os.Exit(1)
	}

	undertaker := loader.Undertaker{
		FPMAddr:      cfg.FPMAddress,
		TombsAddress: cfg.TombsAddress,
		PreloadFile:  cfg.PreloadFile,
	}

	if cfg.Preload {
		if err := undertaker.Preload(); err != nil {
			log.Fatalf("could not run preload: %s", err)
		}
	}

	if cfg.Collect {
		funcs, err := undertaker.Collect()
		if err != nil {
			log.Fatalf("could not run collect: %s", err)
		}

		fmt.Println(strings.Join(funcs, "\n"))
	}

	if cfg.Store != "" {
		c, err := database.NewPostgres(cfg.Store, "__undertaker")
		if err != nil {
			log.Fatalf("could not init postgres: %s", err)
		}

		if err := c.Store([]string{
			"test3",
			"test1",
		}); err != nil {
			log.Fatalf("could not store dump: %s", err)
		}
	}

	if cfg.HTTPPort != "" {
		fmt.Println("Starting server on port: ", cfg.HTTPPort)
		srv := server.NewServer(undertaker, cfg.HTTPPort)
		sigChan := make(chan os.Signal, 5)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGQUIT)
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Fatalf("failed to start server: %s", err)
			}
		}()
		<-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
}
