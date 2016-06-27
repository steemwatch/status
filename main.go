package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/steemwatch/status/checks"

	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Start the notifier.
	statusCh, err := startNotifier()
	if err != nil {
		return err
	}

	// Start the runner.
	runner := startRunner(checks.SetStatusChannel(statusCh))

	// Bind server listener.
	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		return errors.Wrap(err, "listen() failed")
	}
	defer listener.Close()

	// Start the server.
	go startServer(listener, runner)

	// Start catching signals.
	// Interrupt the check runner on signal.
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalCh
		signal.Stop(signalCh)
		runner.Interrupt()
	}()

	// Wait for the runner to exit.
	return runner.Wait()
}
