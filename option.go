package pregen

import (
	"time"
)

type Option[T any] func(g *Generator[T])

func PregenSize[T any](size int) func(g *Generator[T]) {
	return func(g *Generator[T]) {
		g.gc = make(chan T, size)
	}
}

func ErrorCooldown[T any](cooldown time.Duration) func(g *Generator[T]) {
	return func(g *Generator[T]) {
		g.errorCooldown = cooldown
	}
}

func StartDelay[T any](delay time.Duration) func(g *Generator[T]) {
	return func(g *Generator[T]) {
		g.startDelay = delay
	}
}
