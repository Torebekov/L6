package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/Torebekov/L6/internals/goroutine"
	"time"
)

func main() {
	tasks := make([]func(ctx context.Context) error, 0, 10100)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			tasks = append(tasks, func(ctx context.Context) error {
				for {
					select {
					case <-ctx.Done():
						//fmt.Println("task was canceled")
						return nil
					default:
						time.Sleep(time.Duration(i) * time.Microsecond)
						return nil
					}
				}
			})
		}
		tasks = append(tasks, func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("error task was canceled")
					return nil
				default:
					time.Sleep(10 * time.Microsecond)
					return errors.New("some error")
				}
			}
		})
	}
	fmt.Println(goroutine.ExecuteCtx(tasks, 90))
}
