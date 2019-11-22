package runner

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Runnable defines the interface for a controllable looping Goroutine.
type Runnable interface {
	Cleanup()
	Name() string
	Run()
	Setup()
}

// Runner is the lifecycle and loop controller for a Runnable.
type Runner struct {
	Runnable
	context.Context
	stopped chan struct{}
	cancel  context.CancelFunc
}

// NewRunner creates a new Runner for the given Runnable.
func NewRunner(ctx context.Context, runnable Runnable) *Runner {
	ctx, cancel := context.WithCancel(ctx)

	return &Runner{
		Runnable: runnable,
		Context:  ctx,
		stopped:  make(chan struct{}),
		cancel:   cancel,
	}
}

// Start runs the setup steps for the Runnable and starts the looping
// goroutine.
func (runner *Runner) Start() {
	log.Debugf("Starting loop for %s", runner.Name())
	runner.Setup()
	go runner.run()
	log.Debugf("Started loop for %s", runner.Name())
}

// Stop stops the looping goroutine and returns a channel that closes when the
// Runner has come to a complete stop.
func (runner *Runner) Stop() <-chan struct{} {
	defer runner.cancel()
	log.Debugf("Stop requested for %s", runner.Name())
	return runner.stopped
}

// Stopped returns a channel that closes when the Runner has come to a complete
// stop, to allow waiting for a Runnable to stop in a separate goroutine from
// the call to Runner.Stop().
func (runner *Runner) Stopped() <-chan struct{} {
	return runner.stopped
}

func (runner *Runner) run() {
	defer runner.cleanup()
	runner.Run()
	log.Debugf("Stopping loop for %s", runner.Name())
}

func (runner *Runner) cleanup() {
	runner.Cleanup()
	log.Debugf("Stopped loop for %s", runner.Name())
	close(runner.stopped)
}
