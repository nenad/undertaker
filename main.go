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
)

func main() {
	config, err := ParseArgs()
	if err != nil {
		fmt.Printf("could not load configuration: %s", err)
		os.Exit(1)
	}

	undertaker := Undertaker{
		FPMAddr:      config.FPMAddress,
		TombsAddress: config.TombsAddress,
		PreloadFile:  config.PreloadFile,
	}

	if config.Preload {
		if err := undertaker.Preload(); err != nil {
			log.Fatalf("could not run preload: %s", err)
		}
	}

	if config.Collect {
		funcs, err := undertaker.Collect()
		if err != nil {
			log.Fatalf("could not run collect: %s", err)
		}

		fmt.Println(strings.Join(funcs, "\n"))
	}

	if config.HTTPPort != "" {
		fmt.Println("Starting server on port: ", config.HTTPPort)
		srv := NewServer(undertaker, config.HTTPPort)
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
