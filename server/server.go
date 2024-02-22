package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/jakobmoellerdev/splitsmart/api"
	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/logger"
	"github.com/uptrace/bun"
)

// Run will run the HTTP Server.
func Run(ctx context.Context, db *bun.DB, cfg *config.Config) error {
	log := logger.FromContext(ctx)
	startUpContext, cancelStartUpContext := context.WithCancel(ctx)
	defer cancelStartUpContext()

	srv := createServer(startUpContext, cfg)

	idleConsClosed := make(chan struct{})
	closeServer := func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint // We received an interrupt signal, shut down.
		// Set up a context to allow for graceful server shutdowns in the event
		// of an OS interrupt (defers the cancel just in case)
		ctx, cancel := context.WithTimeout(
			startUpContext,
			cfg.Server.Timeout.Server,
		)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// Error from closing listeners, or context timeout:
			log.Warn().Msg("server shutdown error: " + err.Error())
		}
		if err := db.Close(); err != nil {
			log.Warn().Msg("database close error: " + err.Error())
		}
		close(idleConsClosed)
	}

	go closeServer()

	// service connections
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Msg("listen: %s" + err.Error())
	}

	<-idleConsClosed
	log.Info().Msg("server shut down finished")

	return nil
}

func createServer(startUpContext context.Context, cfg *config.Config) *http.Server {
	// Define server options
	srv := &http.Server{
		Addr:              cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:           api.New(startUpContext, cfg),
		ReadTimeout:       cfg.Server.Timeout.Read,
		ReadHeaderTimeout: cfg.Server.Timeout.Read,
		WriteTimeout:      cfg.Server.Timeout.Write,
		IdleTimeout:       cfg.Server.Timeout.Idle,
	}

	return srv
}
