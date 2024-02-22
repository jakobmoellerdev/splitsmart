package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/jakobmoellerdev/splitsmart/config"
	"github.com/jakobmoellerdev/splitsmart/server"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	log := zerolog.New(zerolog.NewTestWriter(t))

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context, t *testing.T) {
		assert.New(t).NoError(
			server.Run(
				log.WithContext(ctx), &config.Config{},
			),
		)
	}(ctx, t)
	time.Sleep(100 * time.Millisecond)
	cancel()
}
