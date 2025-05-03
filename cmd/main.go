package main

import (
	"fmt"
	"os"

	"github.com/Magic-Kot/Store-notification-service/internal/application"
	"github.com/Magic-Kot/Store-notification-service/internal/config"
)

var (
	appName    = "store-notification-service" //nolint:gochecknoglobals
	appVersion = "v0.0.0"                     //nolint:gochecknoglobals
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = application.New(appName, appVersion, cfg).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
