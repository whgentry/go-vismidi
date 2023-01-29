package control

import (
	"context"
	"errors"
	"sync"
)

type State uint

const (
	Stopped State = iota
	Running
)

type RunFunc[IN any, OUT any] func(ctx context.Context, input chan IN, output chan OUT)

type ProcessInterface[IN any, OUT any] interface {
	Run(ctx context.Context, input chan IN, output chan OUT)
	Name() string
}

type IOBlock[IN any, OUT any] struct {
	Input      chan IN
	Output     chan OUT
	Processors map[string]ProcessInterface[IN, OUT]
	active     ProcessInterface[IN, OUT]
	ctx        context.Context
	cancel     func()
	wg         sync.WaitGroup
	mutex      sync.Mutex
	state      State
}

func NewIOBlock[IN any, OUT any](in chan IN, out chan OUT, procs []ProcessInterface[IN, OUT]) *IOBlock[IN, OUT] {
	io := &IOBlock[IN, OUT]{
		Input:      in,
		Output:     out,
		Processors: map[string]ProcessInterface[IN, OUT]{},
		active:     procs[0],
		state:      Stopped,
	}

	// Populate process map
	for _, p := range procs {
		io.Processors[p.Name()] = p
	}

	return io
}

func (io *IOBlock[IN, OUT]) start() {
	ctx, cancel := context.WithCancel(io.ctx)

	io.wg.Add(1)
	go func() {
		io.active.Run(ctx, io.Input, io.Output)
		io.wg.Done()
	}()
	io.state = Running
	io.cancel = cancel
}

// Starts the active processor
func (io *IOBlock[IN, OUT]) Start(ctx context.Context) {
	if io.state == Running {
		return
	}

	io.mutex.Lock()
	defer io.mutex.Unlock()

	io.ctx = ctx
	io.start()
}

func (io *IOBlock[IN, OUT]) stop() {
	io.cancel()
	io.wg.Wait()
	io.state = Stopped
}

// Stops the active processor
func (io *IOBlock[IN, OUT]) Stop() {
	if io.state == Stopped {
		return
	}

	io.mutex.Lock()
	defer io.mutex.Unlock()

	io.stop()
}

// Sets the active processor and starts if not started
func (io *IOBlock[IN, OUT]) SetActive(name string) error {
	io.mutex.Lock()
	defer io.mutex.Unlock()

	// Check if already active
	if name == io.active.Name() {
		switch io.state {
		case Running:
			return nil
		case Stopped:
			io.start()
		}
	}

	proc, ok := io.Processors[name]
	if !ok {
		return errors.New("process name doesn't exist")
	}

	// Stop current process if running
	if io.state == Running {
		io.stop()
	}

	io.active = proc
	io.start()
	return nil
}
