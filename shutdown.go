package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	log "github.com/sirupsen/logrus"
)

func Shutdown(ctx context.Context, cancel context.CancelFunc, extra func()) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.WithField("recover", err).Panic("Cannot gracefully close the app")
			}

			log.Debug("App closed gracefully")
			os.Exit(0)
		}()

		closeApp := func(ctx context.Context, withCancel bool) {
			go func() {
				<-time.After(10 * time.Second)
				log.WithContext(ctx).Fatalf("Forced shutdown")
			}()

			extra()
			if withCancel {
				cancel()
			}
		}

		select {
		case <-ctx.Done():
			closeApp(ctx, false)
		case <-osSignals:
			closeApp(ctx, true)
		}
	}()
}
