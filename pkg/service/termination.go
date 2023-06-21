package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moov-io/base/log"
)

func NewTerminationListener() chan error {
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	return errs
}

func AwaitTermination(logger log.Logger, terminationListener chan error) {
	if err := <-terminationListener; err != nil {
		_ = logger.Fatal().LogErrorf("Terminated: %v", err)
	}
}
