package server

import (
	"context"
	"errors"
	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/service/sql"
	"net/http"
	"os"
	"os/signal"

	"github.com/jakobmoellerdev/splitsmart/api"
)

// Run will run the HTTP Server.
func Run(ctx context.Context, cfg *config.Config) error {
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
			cfg.Logger.Warn().Msg("server shutdown error: " + err.Error())
		}

		close(idleConsClosed)
	}

	go closeServer()

	// service connections
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		cfg.Logger.Fatal().Msg("listen: %s" + err.Error())
	}

	<-idleConsClosed
	cfg.Logger.Info().Msg("server shut down finished")

	return nil
}

func createServer(startUpContext context.Context, cfg *config.Config) *http.Server {

	cfg.Services.Accounts = new(sql.Accounts)

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
