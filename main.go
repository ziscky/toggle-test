package main

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	log := logger.WithFields(logrus.Fields{
		"go":  runtime.Version(),
		"app": "toggle-test",
	})

	srv, err := NewServer(log, loadOptionsFromEnv())
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.Start(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Errorf("failed to start server: %s", err)
				done <- syscall.SIGTERM
			}
		}
	}()

	<-done

	if err := srv.Stop(); err != nil {
		log.Fatalf("failed to shutdown correctly: %s", err)
	}

	log.Info("server shutdown gracefully")
}
