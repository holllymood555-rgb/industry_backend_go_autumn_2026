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

	arrResult := make([]Result[R], len(in))

	var wg sync.WaitGroup
	ch := make(chan int)

	for g := 0; g < min(workers, len(in)); g++ {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case idx, ok := <-ch:
					if !ok {
						return
					}
					val, err := fn(ctx, in[idx])
					if err != nil {
						arrResult[idx].Err = err
						continue
					}
					arrResult[idx].Value = val
				}

			}
		})
	}

	wg.Go(func() {
		defer close(ch)
		for i := range in {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	})

	wg.Wait()

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return arrResult, nil
}