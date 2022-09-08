package control

import (
	"context"
)

type RunFunc[IN any, OUT any] func(ctx context.Context, input chan IN, output chan OUT)

type ProcessInterface[IN any, OUT any] interface {
	Run(ctx context.Context, input chan IN, output chan OUT)
}

type IOBlock[IN any, OUT any] struct {
	Input      chan IN
	Output     chan OUT
	Processors []ProcessInterface[IN, OUT]
	active     ProcessInterface[IN, OUT]
	ctx        context.Context
	cancel     func()
}

// Starts the active processor
func (io *IOBlock[IN, OUT]) Start(ctx context.Context) {
	io.ctx, io.cancel = context.WithCancel(ctx)
	io.active = io.Processors[0]
	go io.active.Run(io.ctx, io.Input, io.Output)
}

// Stops the active processor
func (io *IOBlock[IN, OUT]) Stop() {
	io.cancel()
}
