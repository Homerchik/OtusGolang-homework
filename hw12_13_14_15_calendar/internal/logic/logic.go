package logic

import "github.com/homerchik/hw12_13_14_15_calendar/internal/storage"
import "time"
import "sort"


func EventAfterNow(event storage.Event) error {
	if event.StartDate.Before(time.Now()) {
		return storage.ErrStartTimeBeforeNow
	}
	return nil
}

func EventFitsSchedule(event storage.Event, schedule storage.Schedule) error {
	if len(schedule) == 0 {
		return nil
	}
	sort.Sort(schedule)
	for _, schedEvent := range schedule {
		if schedEvent.StartDate.Before(event.StartDate) && schedEvent.EndDate.Before(event.StartDate) {
			continue
		} else if schedEvent.StartDate.After(event.StartDate) && schedEvent.StartDate.After(event.EndDate) {
			return nil
		} else {
			return storage.ErrEventIntersection
		}
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