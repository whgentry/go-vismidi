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
	Input            chan IN
	Output           chan OUT
	ProcessorIndices map[string]int
	Processors       []ProcessInterface[IN, OUT]
	active           ProcessInterface[IN, OUT]
	ctx              context.Context
	cancel           func()
	wg               sync.WaitGroup
	mutex            sync.Mutex
	state            State
}

func NewIOBlock[IN any, OUT any](in chan IN, out chan OUT, procs []ProcessInterface[IN, OUT]) *IOBlock[IN, OUT] {
	io := &IOBlock[IN, OUT]{
		Input:            in,
		Output:           out,
		Processors:       procs,
		ProcessorIndices: map[string]int{},
		active:           procs[0],
		state:            Stopped,
	}

	// Populate process map
	for i, p := range procs {
		io.ProcessorIndices[p.Name()] = i
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

func (io *IOBlock[IN, OUT]) Next() {
	io.mutex.Lock()
	defer io.mutex.Unlock()

	currentIndex := io.ProcessorIndices[io.active.Name()]
	nextIndex := (currentIndex + 1) % len(io.Processors)

	io.stop()
	io.active = io.Processors[nextIndex]
	io.start()
}

func (io *IOBlock[IN, OUT]) Previous() {
	io.mutex.Lock()
	defer io.mutex.Unlock()

	currentIndex := io.ProcessorIndices[io.active.Name()]
	nextIndex := (currentIndex + len(io.Processors) - 1) % len(io.Processors)

	io.stop()
	io.active = io.Processors[nextIndex]
	io.start()
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

	procIndex, ok := io.ProcessorIndices[name]
	if !ok {
		return errors.New("process name doesn't exist")
	}

	proc := io.Processors[procIndex]

	// Stop current process if running
	if io.state == Running {
		io.stop()
	}

	io.active = proc
	io.start()
	return nil
}
