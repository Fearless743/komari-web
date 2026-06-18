package main

import (
	"log"
	"log/slog"

	"github.com/Fearless743/komari/cmd"
	"github.com/Fearless743/komari/utils"
	logutil "github.com/Fearless743/komari/utils/log"
)

func main() {
	if utils.VersionHash == "unknown" {
		logutil.SetupGlobalLogger(slog.LevelDebug)
	} else {
		logutil.SetupGlobalLogger(slog.LevelInfo)
	}

	log.Printf("Komari Monitor %s (hash: %s)", utils.CurrentVersion, utils.VersionHash)

	cmd.Execute()
}
