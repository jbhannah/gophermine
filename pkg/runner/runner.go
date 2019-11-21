package runner

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type Runnable interface {
	Cleanup()
	Run()
	Setup()
	WithRunner(*Runner) Runnable
}

type Runner struct {
	Runnable
	context.Context
	stopped chan struct{}
	cancel  context.CancelFunc
}

func NewRunner(ctx context.Context, runnable Runnable) Runnable {
	ctx, cancel := context.WithCancel(ctx)

	return runnable.WithRunner(&Runner{
		Runnable: runnable,
		Context:  ctx,
		stopped:  make(chan struct{}),
		cancel:   cancel,
	})
}

func (runner *Runner) Setup() {}

func (runner *Runner) Start() {
	runner.Runnable.Setup()
	log.Debug("Starting")
	go runner.run()
}

func (runner *Runner) Run() {}

func (runner *Runner) Stop() <-chan struct{} {
	defer runner.cancel()
	log.Debug("Stop received")
	return runner.stopped
}

func (runner *Runner) Cleanup() {}

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
