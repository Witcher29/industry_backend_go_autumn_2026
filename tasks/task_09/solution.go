package main

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type Result[R any] struct {
	Value R
	Err   error
}

var ErrInvalidWorkers = errors.New("workers must be positive")

func ParallelMap[T any, R any](ctx context.Context, workers int, in []T,
	fn func(context.Context, T) (R, error)) ([]Result[R], error) {

	if workers <= 0 {
		return nil, ErrInvalidWorkers
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if len(in) == 0 {
		return []Result[R]{}, nil
	}

	countOfTasks := len(in)
	if workers > countOfTasks {
		workers = countOfTasks
	}

	results := make([]Result[R], countOfTasks)

	var indexOfNextTask int64

	var wg sync.WaitGroup
	wg.Add(workers)

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}

				indexOfCurrentTask := int(atomic.AddInt64(&indexOfNextTask, 1)) - 1
				if indexOfCurrentTask >= countOfTasks {
					return
				}

				v, err := fn(ctx, in[indexOfCurrentTask])
				if err != nil {
					results[indexOfCurrentTask] = Result[R]{Err: err}
				} else {
					results[indexOfCurrentTask] = Result[R]{Value: v}
				}
			}
		}()
	}

	wg.Wait()

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
