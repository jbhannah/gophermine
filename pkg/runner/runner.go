package runner

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

type ContextKey string

const (
	UnknownKey      ContextKey = "unknown"
	RunnableStarted ContextKey = "started"
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
	started chan struct{}
	stopped chan struct{}
	cancel  context.CancelFunc
}

// NewRunner creates a new Runner for the given Runnable.
func NewRunner(ctx context.Context, runnable Runnable) *Runner {
	started := make(chan struct{})
	ctx, cancel := context.WithCancel(context.WithValue(ctx, RunnableStarted, started))

	return &Runner{
		Runnable: runnable,
		Context:  ctx,
		started:  started,
		stopped:  make(chan struct{}),
		cancel:   cancel,
	}
}

// Start runs the setup steps for the Runnable and starts the looping goroutine,
// and returns a channel that closes when the runner has started.
func (runner *Runner) Start() <-chan struct{} {
	log.Debugf("Starting loop for %s", runner.Name())
	startTime := time.Now()

	runner.Setup()

	go func() {
		<-runner.started
		log.Debugf("Started loop for %s in %s", runner.Name(), time.Since(startTime))
	}()

	go runner.run()
	return runner.started
}

// Started returns a channel that closes when the Runner has started, to allow
// waiting for a Runnable to start in a separate goroutine from the call to
// Runner.Start().
func (runner *Runner) Started() <-chan struct{} {
	return runner.started
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

	close(runner.started)
	runner.Run()
}

func (runner *Runner) cleanup() {
	defer close(runner.stopped)

	log.Debugf("Stopping loop for %s", runner.Name())
	stopTime := time.Now()

	runner.Cleanup()
	log.Debugf("Stopped loop for %s in %s", runner.Name(), time.Since(stopTime))
}
