package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	wg := sync.WaitGroup{}
	ch := make(chan Task)
	var errorCount atomic.Int32
	maxErrors := int32(m)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range ch {
				if err := task(); err != nil {
					errorCount.Add(1)
				}
			}
		}()
	}
	for _, task := range tasks {
		if m != 0 && errorCount.Load() >= maxErrors {
			close(ch)
			return ErrErrorsLimitExceeded
		}
		ch <- task
	}
	close(ch)
	wg.Wait()
	if m != 0 && errorCount.Load() >= maxErrors {
		return ErrErrorsLimitExceeded
	}
	return nil
}
