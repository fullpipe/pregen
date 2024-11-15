package pregen

import (
	"context"
	"time"
)

type Generator[T any] struct {
	gc            chan T
	genFunc       func() (T, error)
	startDelay    time.Duration
	errorCooldown time.Duration
}

const DefaultPregenSize = 19
const DefaultErrorCooldown = 100 * time.Millisecond

func NewGenerator[T any](
	genFunc func() (T, error),
	options ...Option[T],
) (*Generator[T], func()) {
	return NewGeneratorContext(genFunc, context.Background(), options...)
}

func NewGeneratorContext[T any](
	genFunc func() (T, error),
	ctx context.Context,
	options ...Option[T],
) (*Generator[T], func()) {
	generator := &Generator[T]{
		genFunc:       genFunc,
		gc:            nil,
		startDelay:    time.Duration(0),
		errorCooldown: DefaultErrorCooldown,
	}

	for _, option := range options {
		option(generator)
	}

	if generator.gc == nil {
		generator.gc = make(chan T, DefaultPregenSize)
	}

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		defer close(generator.gc)

		select {
		case <-time.After(generator.startDelay):
		case <-ctx.Done():
		}

		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := generator.genFunc()
				if err != nil {
					time.Sleep(generator.errorCooldown)

					continue
				}

				generator.gc <- data
			}
		}
	}()

	return generator, func() { cancel() }
}

func (g *Generator[T]) Gen() (T, error) {
	select {
	case data, ok := <-g.gc:
		if !ok {
			return g.genFunc()
		}

		return data, nil
	default:
		return g.genFunc()
	}
}
