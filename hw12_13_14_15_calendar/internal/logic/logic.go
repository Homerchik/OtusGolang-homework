package logic

import (
	"sort"
	"time"

	"github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
)

func EventAfterNow(event storage.Event) error {
	if event.StartDate < time.Now().Unix() {
		return storage.ErrStartTimeBeforeNow
	}
	return nil
}

func EventFitsSchedule(event storage.Event, schedule storage.Schedule) error {
	if len(schedule) == 0 {
		return nil
	}
	sort.Sort(schedule)
	for _, scheduleEvent := range schedule {
		if scheduleEvent.StartDate < event.StartDate && scheduleEvent.EndDate < event.StartDate {
			continue
		}
		if scheduleEvent.StartDate > event.StartDate && scheduleEvent.StartDate > event.EndDate {
			return nil
		}
		return storage.ErrEventIntersection
	}
	return nil
}

func CheckEvent(event storage.Event, schedule storage.Schedule) error {
	if err := EventAfterNow(event); err != nil {
		return err
	}
	if err := EventFitsSchedule(event, schedule); err != nil {
		return err
	}
	return nil
}
