package test

import (
	"context"
	"errors"
	"fmt"
	"github.com/Torebekov/L6/internals/goroutine"
	"reflect"
	"testing"
	"time"
)

func TestExecuteChannelMutex(t *testing.T) {
	tasks := make([]func() error, 0, 10100)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			tasks = append(tasks, func() error {
				time.Sleep(10 * time.Microsecond)
				return nil
			})
		}
		tasks = append(tasks, func() error {
			time.Sleep(10 * time.Microsecond)
			return errors.New("some error")
		})
	}
	type Input struct {
		Funcs  []func() error
		MaxErr int
	}
	testTable := []struct {
		input    Input
		expected error
	}{
		{
			input: Input{
				Funcs:  tasks,
				MaxErr: 60,
			},
			expected: errors.New("too much errors"),
		},
		{
			input: Input{
				Funcs:  tasks,
				MaxErr: 100,
			},
			expected: nil,
		},
	}
	//ExecuteChannel
	for _, testCase := range testTable {
		start := time.Now()
		result := goroutine.ExecuteChannel(testCase.input.Funcs, testCase.input.MaxErr)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Incorrect result. Expect %v, got %v",
				testCase.expected, result)
		}
		t.Log(fmt.Sprintf("Channel %v: ", result), time.Since(start))
	}
	//ExecuteMutex
	for _, testCase := range testTable {
		start := time.Now()
		result := goroutine.ExecuteMutex(testCase.input.Funcs, testCase.input.MaxErr)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Incorrect result. Expect %v, got %v",
				testCase.expected, result)
		}
		t.Log(fmt.Sprintf("Mutex %v: ", result), time.Since(start))
	}
}

func TestExecuteCtx(t *testing.T) {
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
						time.Sleep(10 * time.Microsecond)
						return nil
					}
				}
			})
		}
		tasks = append(tasks, func(ctx context.Context) error {
			for {
				select {
				case <-ctx.Done():
					//fmt.Println("error task was canceled")
					return nil
				default:
					time.Sleep(10 * time.Microsecond)
					return errors.New("some error")
				}
			}
		})
	}
	type Input struct {
		Funcs  []func(ctx context.Context) error
		MaxErr int
	}
	testTable := []struct {
		input    Input
		expected error
	}{
		{
			input: Input{
				Funcs:  tasks,
				MaxErr: 60,
			},
			expected: errors.New("too much errors"),
		},
		{
			input: Input{
				Funcs:  tasks,
				MaxErr: 100,
			},
			expected: nil,
		},
	}
	for _, testCase := range testTable {
		start := time.Now()
		result := goroutine.ExecuteCtx(testCase.input.Funcs, testCase.input.MaxErr)
		if !reflect.DeepEqual(result, testCase.expected) {
			t.Errorf("Incorrect result. Expect %v, got %v",
				testCase.expected, result)
		}
		t.Log(fmt.Sprintf("Context %v: ", result), time.Since(start))
	}
}
