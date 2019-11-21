package runner

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// Runnable defines the interface for a controllable looping Goroutine.
type Runnable interface {
	Cleanup()
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
	runner.Runnable.Setup()
	log.Debug("Starting")
	go runner.run()
}

// Setup is an empty no-op for Runnables with no setup steps.
func (runner *Runner) Setup() {}

// Run is an empty no-op for Runnables with no run steps.
func (runner *Runner) Run() {}

// Stop stops the looping goroutine and returns a channel that closes when the
// Runner has come to a complete stop.
func (runner *Runner) Stop() <-chan struct{} {
	defer runner.cancel()
	log.Debug("Stop received")
	return runner.stopped
}

// Cleanup is an empty no-op for Runnables with no cleanup steps.
func (runner *Runner) Cleanup() {}

// Stopped returns a channel that closes when the Runner has come to a complete
// stop, to allow waiting for a Runnable to stop in a separate goroutine from
// the call to Runner.Stop().
func (runner *Runner) Stopped() <-chan struct{} {
	return runner.stopped
}

func (runner *Runner) run() {
	defer runner.cleanup()
	log.Debug("Running")

	for {
		select {
		case <-runner.Done():
			log.Debug("Stopping")
			return
		default:
			runner.Run()
		}
	}
}

func (runner *Runner) cleanup() {
	runner.Cleanup()
	log.Debug("Stopped")
	close(runner.stopped)
}
