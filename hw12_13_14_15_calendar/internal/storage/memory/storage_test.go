package memorystorage

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"                                       //nolint:depguard
	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/stretchr/testify/require"
)

func TestEventAddedToEmptyStorage(t *testing.T) {
	testStorage := New()
	startTime := time.Now().UTC()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1",
		startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	require.NoError(t, testStorage.AddEvent(event))
	require.Equal(t, testStorage.Events[startTime.Add(time.Hour)][0], event)
}

func TestEventAddedBeforeNow(t *testing.T) {
	testStorage := New()
	startTime := time.Now()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1",
		startTime.Add(-time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	require.Error(t, testStorage.AddEvent(event), storage.ErrStartTimeBeforeNow)
}

func TestEventIntersection(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}
	t.Run("Overlap with one of the events", func(t *testing.T) {
		minutes := 10 * time.Minute
		event := storage.NewEvent(
			userID, "Event 4", "Description 4",
			startTime.Add(minutes+time.Hour), startTime.Add(2*time.Hour-minutes), time.Minute,
		)
		require.Error(t, testStorage.AddEvent(event), storage.ErrEventIntersection)
	})

	t.Run("Overlap with two events", func(t *testing.T) {
		minutes := 10 * time.Minute
		event := storage.NewEvent(
			userID, "Event 4", "Description 4",
			startTime.Add(minutes+time.Hour), startTime.Add(3*time.Hour+minutes), time.Minute,
		)
		require.Error(t, testStorage.AddEvent(event), storage.ErrEventIntersection)
	})
}

func TestEventIntersectionWithOtherUsers(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	events := storage.Schedule{
		storage.NewEvent(
			uuid.New(), "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}
	minutes := 10 * time.Minute
	newEventStart := startTime.Add(minutes + 2*time.Hour)
	event := storage.NewEvent(
		userID, "Event 4", "Description 4",
		newEventStart, newEventStart.Add(2*time.Hour+minutes), time.Minute,
	)
	require.NoError(t, testStorage.AddEvent(event), storage.ErrEventIntersection)
	require.Equal(t, testStorage.Events[newEventStart][0], event)
}

func TestDeleteExistingEventOneForADate(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}
	require.NoError(t, testStorage.DeleteEvent(events[1].ID))
	require.NotContains(t, testStorage.Events, events[1].StartDate)
}

func TestDeleteExistingEventMultipleForADate(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 2-1", "Description 2-1",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 3-1", "Description 3-1",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}
	require.NoError(t, testStorage.DeleteEvent(events[1].ID))
	require.Contains(t, testStorage.Events, events[1].StartDate)
	require.NotContains(t, testStorage.Events[events[1].StartDate], events[1])
}

func TestDeleteUnexistingEvent(t *testing.T) {
	testStorage := New()
	startTime := time.Now().UTC()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1",
		startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	require.NoError(t, testStorage.AddEvent(event))
	require.Error(t, testStorage.DeleteEvent(uuid.New()), storage.ErrNoEventFound)
}

func TestEventUpdateSimpleFields(t *testing.T) {
	testStorage := New()
	startTime := time.Now().UTC()
	event := storage.NewEvent(
		uuid.New(), "Event 1", "Description 1",
		startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
	)
	require.NoError(t, testStorage.AddEvent(event))
	updatedEvent := storage.NewEvent(
		event.UserID, "Better than event 1", "Simple Des",
		startTime.Add(time.Hour), startTime.Add(2*time.Hour), 10*time.Minute,
	)
	updatedEvent.ID = event.ID
	require.NoError(t, testStorage.UpdateEvent(updatedEvent))
	require.Equal(t, updatedEvent, testStorage.Events[updatedEvent.StartDate][0])
}

func TestEventUpdateDateFields(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}
	t.Run("Check start date is changed, it's possible", func(t *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(30*time.Minute), events[0].EndDate, time.Minute,
		)
		updatedEvent.ID = events[0].ID
		require.NoError(t, testStorage.UpdateEvent(updatedEvent))
		require.Equal(t, updatedEvent, testStorage.Events[updatedEvent.StartDate][0])
	})

	t.Run("Check end date is changed, and it's possible", func(t *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1",
			events[0].StartDate, events[0].EndDate.Add(-30*time.Minute), time.Minute,
		)
		updatedEvent.ID = events[0].ID
		require.NoError(t, testStorage.UpdateEvent(updatedEvent))
		require.Equal(t, updatedEvent, testStorage.Events[updatedEvent.StartDate][0])
	})

	t.Run("Check start date changed, moved in the schedule", func(t *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(7*time.Hour), startTime.Add(8*time.Hour), time.Minute,
		)
		updatedEvent.ID = events[0].ID
		require.NoError(t, testStorage.UpdateEvent(updatedEvent))
		require.Equal(t, updatedEvent, testStorage.Events[updatedEvent.StartDate][0])
	})
}

func TestEventUpdateDateFieldsErrors(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		require.NoError(t, testStorage.AddEvent(event))
	}

	t.Run("Check start date changed, won't fit anymore, error thrown", func(t *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		)
		updatedEvent.ID = events[0].ID
		require.Error(t, testStorage.UpdateEvent(updatedEvent), storage.ErrEventCantBeUpdated)
		require.Equal(t, events[0], testStorage.Events[events[0].StartDate][0])
		require.Equal(t, events[1], testStorage.Events[events[1].StartDate][0])
		require.Equal(t, events[2], testStorage.Events[events[2].StartDate][0])
	})

	t.Run("Check end date changed, won't fit anymore, error thrown", func(t *testing.T) {
		updatedEvent := storage.NewEvent(
			userID, "Event 1", "Description 1",
			events[0].StartDate, startTime.Add(4*time.Hour), time.Minute,
		)
		updatedEvent.ID = events[0].ID
		require.Error(t, testStorage.UpdateEvent(updatedEvent), storage.ErrEventCantBeUpdated)
		require.Equal(t, events[0], testStorage.Events[events[0].StartDate][0])
		require.Equal(t, events[1], testStorage.Events[events[1].StartDate][0])
		require.Equal(t, events[2], testStorage.Events[events[2].StartDate][0])
	})
}

func TestAddingDeletingEventsInParallel(t *testing.T) {
	testStorage := New()
	userID := uuid.New()
	startTime := time.Now().UTC()
	wg := &sync.WaitGroup{}
	events := storage.Schedule{
		storage.NewEvent(
			userID, "Event 1", "Description 1",
			startTime.Add(time.Hour), startTime.Add(2*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 2", "Description 2",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			userID, "Event 3", "Description 3",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 2-1", "Description 2-1",
			startTime.Add(3*time.Hour), startTime.Add(4*time.Hour), time.Minute,
		),
		storage.NewEvent(
			uuid.New(), "Event 3-1", "Description 3-1",
			startTime.Add(5*time.Hour), startTime.Add(6*time.Hour), time.Minute,
		),
	}
	for _, event := range events {
		wg.Add(1)
		go func() {
			defer wg.Done()
			require.NoError(t, testStorage.AddEvent(event))
		}()
	}
	wg.Wait()
	require.Equal(t, len(testStorage.Events), 3)
	require.Equal(t, len(testStorage.Events[startTime.Add(3*time.Hour)]), 2)
	require.Equal(t, len(testStorage.Events[startTime.Add(5*time.Hour)]), 2)
	for _, event := range events {
		wg.Add(1)
		go func() {
			defer wg.Done()
			require.NoError(t, testStorage.DeleteEvent(event.ID))
		}()
	}
	wg.Wait()
	require.Equal(t, len(testStorage.Events), 0)
}
