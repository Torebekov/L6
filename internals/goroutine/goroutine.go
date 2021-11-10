package goroutine

import (
	"context"
	"errors"
	"runtime"
	"sync"
)

func ExecuteChannel(tasks []func() error, E int) error {
	chErr := make(chan struct{}, runtime.GOMAXPROCS(-1))
	chOk := make(chan struct{}, runtime.GOMAXPROCS(-1))
	done := make(chan struct{})
	cntOk, cntErr := 0, 0
	go func(chErr <-chan struct{}, chOk <-chan struct{}, done chan<- struct{}) {
		for {
			select {
			case <-chOk:
				cntOk++
			case <-chErr:
				cntErr++
			default:
				if cntErr > E || cntErr+cntOk == len(tasks) {
					done <- struct{}{}
					return
				}
			}
		}
	}(chErr, chOk, done)
	for i := range tasks {
		go func(i int, chErr chan<- struct{}, chOk chan<- struct{}) {
			err := tasks[i]()
			if err != nil {
				chErr <- struct{}{}
			} else {
				chOk <- struct{}{}
			}
		}(i, chErr, chOk)
	}
	<-done
	if cntErr > E {
		return errors.New("too much errors")
	}
	return nil
}

func ExecuteMutex(tasks []func() error, E int) error {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	cntErr := 0
	for i := range tasks {
		wg.Add(1)
		go func(i int) {
			err := tasks[i]()
			if err != nil {
				mu.Lock()
				cntErr++
				mu.Unlock()
			}
			wg.Done()
		}(i)
		mu.Lock()
		if cntErr > E {
			return errors.New("too much errors")
		}
		mu.Unlock()
	}
	wg.Wait()
	if cntErr > E {
		return errors.New("too much errors")
	}
	return nil
}

func ExecuteCtx(tasks []func(ctx context.Context) error, E int) error {
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	cntErr := 0
	for i := range tasks {
		wg.Add(1)
		go func(i int) {
			err := tasks[i](ctx)
			if err != nil {
				mu.Lock()
				cntErr++
				mu.Unlock()
			}
			mu.Lock()
			if cntErr > E {
				cancelFunc()
			}
			mu.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	if cntErr > E {
		return errors.New("too much errors")
	}
	return nil
}
