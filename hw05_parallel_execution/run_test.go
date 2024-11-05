package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount atomic.Int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				runTasksCount.Add(1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount.Load(), int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		runTasksCount, elapsedTime, sumTime, err := createAndEvaluate(50, 0, 5, 1)
		require.NoError(t, err)
		require.Equal(t, runTasksCount, int32(50), "not all tasks were completed")
		require.LessOrEqual(t, elapsedTime, sumTime/2, "tasks were run sequentially?")
	})

	t.Run("tasks with all errors, but maxErrors set to ignore", func(t *testing.T) {
		runTasksCount, elapsedTime, sumTime, err := createAndEvaluate(0, 50, 10, 0)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(50), "not all tasks were completed")
		require.LessOrEqual(t, elapsedTime, sumTime/2, "tasks were run sequentially?")
	})

	t.Run("All tasks are finished, but error threshold exceeded", func(t *testing.T) {
		runTasksCount, _, _, err := createAndEvaluate(44, 1, 60, 1)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)

		require.Equal(t, runTasksCount, int32(45), "not all tasks were completed")
	})
}

func createAndEvaluate(fineJobCount, errorJobCount, wksCount, maxErrorsCount int) (int32, int64, time.Duration, error) {
	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount atomic.Int32
	var sumTime time.Duration

	for i := 0; i < errorJobCount; i++ {
		taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
		sumTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			runTasksCount.Add(1)
			return fmt.Errorf("error from task %d", i)
		})
	}
	for i := 0; i < fineJobCount; i++ {
		taskSleep := time.Millisecond * time.Duration(1)
		sumTime += taskSleep

		tasks = append(tasks, func() error {
			time.Sleep(taskSleep)
			runTasksCount.Add(1)
			return nil
		})
	}
	start := time.Now()
	err := Run(tasks, wksCount, maxErrorsCount)
	elapsedTime := time.Since(start)
	return runTasksCount.Load(), int64(elapsedTime), sumTime, err
}
