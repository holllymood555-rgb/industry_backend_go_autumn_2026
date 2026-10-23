package main

import (
	"context"
	"errors"
	"sync"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T, fn func(context.Context, T) (R, error)) ([]Result[R], error) {
	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if len(in) == 0 {
		return []Result[R]{}, nil
	}

	result := make([]Result[R], len(in))
	tasks := make(chan int, len(in))
	for i := range in {
		tasks <- i
	}
	close(tasks)

	var count sync.WaitGroup
	actualWorkers := workers
	if actualWorkers > len(in) {
		actualWorkers = len(in)
	}

	count.Add(actualWorkers)
	for i := 0; i < actualWorkers; i++ {
		go func() {
			defer count.Done()
			for {
				if ctx.Err() != nil {
					return
				}
				
				select {
				case <-ctx.Done():
					return
				case idx, ok := <-tasks:
					if !ok {
						return
					}

					val, err := fn(ctx, in[idx])
					if err != nil {
						var zero R
						result[idx] = Result[R]{Value: zero, Err: err}
					} else {
						result[idx] = Result[R]{Value: val, Err: nil}
					}
				}
			}
		}()
	}

	count.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	return result, nil
}