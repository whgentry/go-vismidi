package main

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type ProcessInterface[IN any, OUT any] interface {
	Run(ctx context.Context, input chan IN, output chan OUT)
}

type IOBlock[IN any, OUT any] struct {
	Input           chan IN
	Output          chan OUT
	ActiveProcessor ProcessInterface[IN, OUT]
	Processors      []ProcessInterface[IN, OUT]
	ctx             context.Context
	cancelProcessor func()
}

// Starts the active processor
func (io *IOBlock[IN, OUT]) Start(ctx context.Context) {
	io.ctx, io.cancelProcessor = context.WithCancel(ctx)
	go io.ActiveProcessor.Run(io.ctx, io.Input, io.Output)
}

// Stops the active processor
func (io *IOBlock[IN, OUT]) Stop() {
	io.cancelProcessor()
}

type IntToString struct{}

func (p IntToString) Run(ctx context.Context, input chan int, output chan string) {
	for {
		select {
		case packet := <-input:
			fmt.Printf("int %d\n", packet)
			output <- strconv.Itoa(packet + 1)
		case <-ctx.Done():
			return
		}
	}
}

type StringToInt struct{}

func (p StringToInt) Run(ctx context.Context, input chan string, output chan int) {
	for {
		select {
		case packet := <-input:
			fmt.Printf("string %s\n", packet)
			val, _ := strconv.Atoi(packet)
			output <- val + 1
		case <-ctx.Done():
			return
		}
	}
}

type Generator struct{}

func (p Generator) Run(ctx context.Context, _ chan any, output chan int) {
	generateTicker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-generateTicker.C:
			val := rand.Int()
			fmt.Printf("Generated %d\n", val)
			output <- val
		case <-ctx.Done():
			return
		}
	}
}

type Printer struct{}

func (p Printer) Run(ctx context.Context, input chan int, _ chan any) {
	for {
		select {
		case val := <-input:
			fmt.Printf("Final Output: %d\n", val)
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	chan1 := make(chan int)
	chan2 := make(chan string)
	chan3 := make(chan int)

	generatorCB := IOBlock[any, int]{Input: nil, Output: chan1, ActiveProcessor: Generator{}}
	intToStringCB := IOBlock[int, string]{Input: chan1, Output: chan2, ActiveProcessor: IntToString{}}
	stringToIntCB := IOBlock[string, int]{Input: chan2, Output: chan3, ActiveProcessor: StringToInt{}}
	printerCB := IOBlock[int, any]{Input: chan3, Output: nil, ActiveProcessor: Printer{}}

	ctx, _ := context.WithCancel(context.Background())

	generatorCB.Start(ctx)
	intToStringCB.Start(ctx)
	stringToIntCB.Start(ctx)
	printerCB.Start(ctx)

	time.Sleep(5 * time.Second)

	intToStringCB.Stop()

	time.Sleep(5 * time.Second)

	intToStringCB.Start(ctx)

	for {
	}
}
