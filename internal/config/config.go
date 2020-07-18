package config

import (
	"flag"
	"fmt"
	"os"
)

type (
	Config struct {
		TombsAddress string
		FPMAddress   string
		PreloadFile  string

		HTTPPort string

		Preload bool
		Collect bool
	}
)

var (
	fTombs    = flag.String("tombs", "", "Points to the tombs tcp address")
	fFPM      = flag.String("fpm", "", "Points to the FPM tcp address")
	fPreload  = flag.String("file", "", "Points to the undertaker.php preload file")
	fHTTPPort = flag.String("port", "", "Port to listen to incoming connections. Empty port will disable the HTTP server.")

	fPreloadAction = flag.Bool("preload", true, "Runs preload action if specified")
	fCollectAction = flag.Bool("collect", false, "Runs collect action if specified")
)

func ParseArgs() (*Config, error) {
	flag.Parse()
	tombsAddr, err := merge(fTombs, "TOMBS_ADDRESS")
	if err != nil {
		return nil, fmt.Errorf("tombs address not provided")
	}
	fpmAddr, err := merge(fFPM, "FPM_ADDRESS")
	if err != nil {
		return nil, fmt.Errorf("fpm address not provided")
	}
	preloadFile, err := merge(fPreload, "PRELOAD_FILE")
	if err != nil {
		return nil, fmt.Errorf("preload file not provided")
	}

	port, _ := merge(fHTTPPort, "HTTP_PORT")

	c := &Config{
		TombsAddress: tombsAddr,
		FPMAddress:   fpmAddr,
		PreloadFile:  preloadFile,
		HTTPPort:     port,
		Collect:      *fCollectAction,
		Preload:      *fPreloadAction,
	}

	return c, nil
}

func merge(flagVal *string, env string) (string, error) {
	if flagVal != nil && *flagVal != "" {
		return *flagVal, nil
	}
	if os.Getenv(env) == "" {
		return "", fmt.Errorf("could not find value for env %s", env)
	}

	return os.Getenv(env), nil
}
