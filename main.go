package main

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/server"
	"github.com/sethvargo/go-password/password"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	baseLogger "github.com/rs/zerolog/log"
)

var (
	version = "dev"
	commit  = "none"    //nolint:gochecknoglobals
	date    = "unknown" //nolint:gochecknoglobals
	builtBy = "unknown" //nolint:gochecknoglobals
)

// Func main should be as small as possible and do as little as possible by convention.
func main() {
	// Generate our config based on the config supplied
	// by the user in the flags
	cfgPath, err := config.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logger zerolog.Logger

	switch cfg.LogSettings.Format {
	case config.LogSettingsFormatPretty:
		logger = zerolog.New(zerolog.NewConsoleWriter()).With().Timestamp().Logger()
	case config.LogSettingsFormatJSON:
		fallthrough
	case config.LogSettingsFormatNone:
		fallthrough
	default:
		logger = baseLogger.With().Logger()
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		revision := info.Main.Version
		if revision == "" {
			revision = "dev"
		}
		logger = logger.With().Str("revision", revision).Logger()
		logger.Info().
			Str("version", version).
			Str("commit", commit).
			Str("date", date).
			Str("builtBy", builtBy).
			Msg("Welcome to Splitsmart!")
	}

	ctx := logger.WithContext(context.Background())

	cfg.PasswordGenerator, err = password.NewGenerator(nil)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	uuid.EnableRandPool()

	// Run the server
	if err := server.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
